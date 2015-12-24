package omni

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"strings"
	"time"
)

type sessionkey []byte

// Client is an Omni-link II client
type Client struct {
	Addr         string // IP:Port
	conn         net.Conn
	protoVersion uint16 // Protocol version used by the controller
	sessionID    []byte
	sessionKey   sessionkey
	cipher       cipher.Block
	seqNum       uint16
}

func NewClient(addr string, key string) (*Client, error) {
	conn, err := net.DialTimeout("tcp", addr, time.Duration(10*time.Second))
	if err != nil {
		return nil, err
	}
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	fmt.Println("Connected!")
	client := Client{
		Addr:   addr,
		conn:   conn,
		seqNum: 1,
	}
	m := genmsg{
		SeqNum: client.nextSeqNum(),
		Type:   ClientReqNewSession,
	}
	fmt.Printf("Sending new session req %+v\n", m)
	err = client.Send(m)
	if err != nil {
		return nil, fmt.Errorf("Failed to send %s", err.Error())
	}
	ackMsg, err := client.Receive()
	if err != nil {
		return nil, fmt.Errorf("Failed to receive %s", err.Error())
	}
	buf := bytes.NewReader(ackMsg.Data)

	err = binary.Read(buf, binary.LittleEndian, &client.protoVersion)
	if err != nil {
		return nil, fmt.Errorf("Failed to read version %s", err.Error())
	}
	fmt.Printf("Protocol version %v\n", client.protoVersion)
	idbuf := bytes.Buffer{}
	_, err = idbuf.ReadFrom(buf)
	if err != nil {
		return nil, fmt.Errorf("Failed to read sessionID  %s", err.Error())
	}
	client.sessionID = idbuf.Bytes()
	client.sessionKey, err = createSessionKey(key, client.sessionID)
	if err != nil {
		return nil, fmt.Errorf("Failed to create session key - %s", err)
	}
	fmt.Printf("New session %+v\n", client)

	verify := genmsg{
		SeqNum: client.nextSeqNum(),
		Type:   ClientReqSecureConnection,
		Data:   client.sessionID,
	}
	err = client.Send(verify)
	if err != nil {
		return nil, fmt.Errorf("Failed to send req secure conn  %s", err.Error())
	}
	secMsg, err := client.Receive()
	if err != nil {
		return nil, fmt.Errorf("Failed to setup secure connection - %s", err.Error())
	}
	if secMsg.Type != ControllerAckSecureConnection {
		return nil, fmt.Errorf("Client generated wrong session key")
	}
	fmt.Printf("Secure connection established %+v\n", *secMsg)
	client.cipher, err = aes.NewCipher(client.sessionKey)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client cipher - %s", err.Error())
	}
	return &client, err
}

func (c *Client) GetSystemInformation() (SystemInfo, error) {
	seqNum := c.nextSeqNum()
	data := ReqSystemInfoMsg.serialize(c, seqNum)
	err := c.Send(genmsg{
		SeqNum: seqNum,
		Type:   AppDataMsg,
		Data:   data,
	})
	if err != nil {
		return SystemInfo{}, fmt.Errorf("Faile dto send - %s", err.Error())
	}
	msg, err := c.Receive()
	if err != nil {
		return SystemInfo{}, fmt.Errorf("Failed to receive system info %s", err.Error())
	}
	fmt.Printf("sysinfo %+v\n", msg)
	buf := bytes.NewBuffer(msg.Data)
	si := SystemInfo{}

	err = binary.Read(buf, binary.LittleEndian, &si.ModelNumber)
	err = binary.Read(buf, binary.LittleEndian, &si.MajorVersion)
	err = binary.Read(buf, binary.LittleEndian, &si.MinorVerison)
	err = binary.Read(buf, binary.LittleEndian, &si.Revesion)
	err = binary.Read(buf, binary.LittleEndian, &si.LocalPhoneNumber)

	return si, err
}

func (c *Client) nextSeqNum() uint16 {
	next := c.seqNum
	c.seqNum++
	if c.seqNum == maxSeqNum {
		c.seqNum = 1
	}
	return next
}

func createSessionKey(key string, sessionID []byte) (sessionkey, error) {
	keyb, err := parseKey(key)
	if err != nil {
		return sessionkey{}, err
	}
	for i := 11; i < 16; i++ {
		keyb[i] ^= sessionID[i-11]
	}
	return keyb, nil
}

func parseKey(key string) (sessionkey, error) {
	hexOnly := strings.Replace(key, "-", "", -1)
	keyBytes, err := hex.DecodeString(hexOnly)
	if err != nil {
		return sessionkey{}, err
	}
	if len(keyBytes) != 16 {
		return sessionkey{}, fmt.Errorf("Key %s must be 16 bytes long", key)
	}
	return sessionkey(keyBytes), nil
}

func (c *Client) Receive() (*genmsg, error) {
	fmt.Printf("Receving msg...")
	w := bytes.NewBuffer([]byte{})
	buf := make([]byte, 1024)
	for {
		n, err := c.conn.Read(buf)
		if err != nil {
			return nil, err
		}
		w.Write(buf[0:n])
		if n < 1024 {
			break
		}
	}

	m, err := deserialize(w)
	if err == nil {
		fmt.Printf(" %+v\n", m)
	}
	if m.Type == AppDataMsg {
		m.decrypt(c.cipher)
	}

	return m, err
}

func (c *Client) Send(m genmsg) error {
	fmt.Printf("Sending msg type %+v\n", m)
	out := m.serialize(c.cipher)
	fmt.Printf("Sending msg bytes %v\n", out)
	_, err := c.conn.Write(out)
	return err
}
