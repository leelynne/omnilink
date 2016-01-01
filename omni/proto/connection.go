package proto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

type StaticKey []byte
type sessionKey []byte

type Conn interface {
	Read(timeout time.Duration) (*Msg, error)
	Write(m *Msg, timeout time.Duration) error
	Close() error
}

type ConnError struct {
	// Op is the operation which caused the error, such as "read" or "write".
	Op string
	// Addr is the network address of the connection
	Addr string

	// Err is the error that occurred during the operation.
	Err error
}

func (ce ConnError) Error() string {
	return fmt.Sprintf("%s %s %s", ce.Op, ce.Addr, ce.Err.Error())
}

type conn struct {
	mu           sync.Mutex
	nconn        net.Conn
	protoVersion uint16 // Protocol version used by the controller
	sessionKey   sessionKey
	cipher       cipher.Block
	seqNum       uint16
	addr         string
	err          error
	closed       bool
}

// NewConnection will create a new connection and session with the controller
func NewConnection(addr string, key StaticKey) (Conn, error) {
	logger := log.New(os.Stdout, "connection: ", log.LstdFlags)
	nconn, err := net.DialTimeout("tcp", addr, time.Duration(10*time.Second))
	if err != nil {
		return nil, ConnError{Op: "dial", Addr: addr, Err: err}
	}

	oconn := &conn{
		addr:   addr,
		nconn:  nconn,
		seqNum: 1,
	}
	newp := packet{
		seqNum:  oconn.nextSeqNum(),
		msgType: msgClientReqNewSession,
	}
	timeout := time.Now().Add(time.Second * 15)
	err = oconn.sendPacket(newp, timeout)
	if err != nil {
		return nil, err
	}

	ackSessionp, err := oconn.recvPacket(timeout)
	if err != nil {
		return nil, err
	}
	if ackSessionp.msgType != msgControllerAckNewSession {
		return nil, fmt.Errorf("Could not establish new session with controller")
	}
	as := ackNewSession{}
	err = ackSessionp.deserialize(&as)
	if err != nil {
		return nil, err
	}

	oconn.sessionKey, err = createSessionKey(key, as.SessionID[:])
	if err != nil {
		return nil, fmt.Errorf("Failed to create session key - %s", err)
	}

	oconn.cipher, err = aes.NewCipher(oconn.sessionKey[:])
	if err != nil {
		return nil, fmt.Errorf("Failed to create client cipher - %s", err.Error())
	}
	secp := packet{
		seqNum:  oconn.nextSeqNum(),
		msgType: msgClientReqSecureConnection,
		data:    as.SessionID[:],
	}
	err = oconn.sendPacket(secp, timeout)
	if err != nil {
		return nil, err
	}

	ackSecurep, err := oconn.recvPacket(timeout)
	if err != nil {
		return nil, err
	}
	if ackSecurep.msgType != msgControllerAckSecureConnection {
		return nil, fmt.Errorf("Client generated wrong session key")
	}
	sec := ackSecureSession{}
	err = ackSecurep.deserialize(&sec)
	if err != nil {
		return nil, err
	}
	if sec.SessionID != as.SessionID {
		return nil, fmt.Errorf("Failed to match session id on secure connection.")
	}

	return oconn, err
}

func (c *conn) Read(timeout time.Duration) (*Msg, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.ok() {
		return nil, c.err
	}

	p, err := c.recvPacket(time.Now().Add(timeout))
	if err != nil {
		return nil, err
	}

	return fromPacket(p), nil
}

func (c *conn) Write(m *Msg, timeout time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.ok() {
		return c.err
	}

	return c.sendPacket(c.makePacket(m), time.Now().Add(timeout))
}

// Close will close the connection. Close can be called multiple times.
func (c *conn) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return nil
	}
	c.closed = true
	if c.err == nil {
		c.err = fmt.Errorf("Connection is closed.")
	}
	neterr := c.nconn.Close()
	if neterr != nil {
		c.err = neterr
	}
	return neterr
}

func (c *conn) ok() bool {
	if c.closed {
		return false
	}
	return true
}

func (c *conn) makePacket(m *Msg) packet {
	plaintext := m.serialize()
	p := packet{
		seqNum:  c.nextSeqNum(),
		msgType: msgAppData,
		data:    plaintext,
	}
	return p
}

func (c *conn) nextSeqNum() uint16 {
	next := c.seqNum
	c.seqNum++
	if c.seqNum == maxSeqNum {
		c.seqNum = 1
	}
	return next
}

func (c *conn) sendPacket(p packet, timeout time.Time) error {
	b := c.serializePacket(p)
	c.nconn.SetWriteDeadline(timeout)
	for written := 0; written < len(b); {
		n, err := c.nconn.Write(b[written:])
		if err != nil {
			return ConnError{Op: "write", Addr: c.addr, Err: err}
		}
		written += n
	}
	return nil
}

func (c *conn) recvPacket(timeout time.Time) (packet, error) {
	header, err := c.getBytes(4, timeout)
	if err != nil {
		return packet{}, err
	}
	p := packet{}
	buf := bytes.NewBuffer(header)
	err = binary.Read(buf, binary.LittleEndian, &p.seqNum)
	if err != nil {
		return p, ConnError{Op: "read", Addr: c.addr, Err: err}
	}
	err = binary.Read(buf, binary.LittleEndian, &p.msgType)
	if err != nil {
		return p, ConnError{Op: "read", Addr: c.addr, Err: err}
	}
	err = binary.Read(buf, binary.LittleEndian, &p.reserved)
	if err != nil {
		return p, ConnError{Op: "read", Addr: c.addr, Err: err}
	}

	dataLen := 0
	encrypted := false
	switch p.msgType {
	case msgControllerCannotStartNewSession, msgControllerSessionTerminated:
	case msgControllerAckNewSession:
		dataLen = 7
	case msgControllerAckSecureConnection, msgAppData:
		encrypted = true
	default:
		return p, fmt.Errorf("Unknown message type %d", p.msgType)
	}
	if encrypted {
		// Read the first block of the encrypted data in order to get size of the message
		p.data, err = c.getBytes(blockSize, timeout)
		if err != nil {
			return p, err
		}
		msgHeader := p.decrypt(c.cipher)
		unencryptedLength := int(msgHeader[1:2][0])  // Length of unencrypted data + message type
		totalLength := unencryptedLength + 4         // Length + start char, length, and crc fields
		dataLen = padLength(totalLength) - blockSize // Subtract what was already read
	}
	data, err := c.getBytes(dataLen, timeout)
	if err != nil {
		return p, err
	}
	p.data = append(p.data, data...)
	if encrypted {
		p.data = p.decrypt(c.cipher)
	}
	return p, nil
}

func (c *conn) getBytes(numBytes int, timeout time.Time) ([]byte, error) {
	if numBytes <= 0 {
		return []byte{}, nil
	}
	buf := make([]byte, numBytes)
	c.nconn.SetReadDeadline(timeout)
	// Grab at least the header header
	read := 0
	for read < numBytes {
		n, err := c.nconn.Read(buf)
		if err != nil {
			return nil, ConnError{Op: "read", Addr: c.addr, Err: err}
		}
		read += n
	}
	return buf, nil
}

func (c *conn) serializePacket(p packet) []byte {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, p.seqNum)
	binary.Write(buf, binary.LittleEndian, p.msgType)
	binary.Write(buf, binary.LittleEndian, p.reserved)
	data := p.data
	if c.cipher != nil {
		data = p.encrypt(c.cipher)
	}
	binary.Write(buf, binary.LittleEndian, data)
	return buf.Bytes()
}

func createSessionKey(key StaticKey, sessionID []byte) (sessionKey, error) {
	skey := make([]byte, 16)
	copy(skey[:], key[:])
	for i := 11; i < 16; i++ {
		skey[i] ^= sessionID[i-11]
	}
	return skey, nil
}
