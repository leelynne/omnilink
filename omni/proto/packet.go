package proto

import (
	"bytes"
	"crypto/cipher"
	"encoding/binary"
)

type msgType uint8

const (
	msgNone = iota
	msgClientReqNewSession
	msgControllerAckNewSession
	msgClientReqSecureConnection
	msgControllerAckSecureConnection
	msgClientSessionTerminated
	msgControllerSessionTerminated
	msgControllerCannotStartNewSession

	msgAppData = 32

	maxSeqNum = 65535

	blockSize = 16
)

// packet corresponds to omnilink application packet
type packet struct {
	seqNum   uint16  // Sequence Number From 1 - 65535. Zero means no sequence tracking
	msgType  msgType // Type of this message. Specifies the format of the data field
	reserved byte    // Not used. Always zero
	data     []byte  // Application message data
}

/* Packet based messages */
// ackNewSession is a response to the initial new session request
type ackNewSession struct {
	ProtoVersion uint16
	SessionID    [5]byte
}

// ackSecureSession is a response to Client Request Secure Connection
type ackSecureSession struct {
	SessionID [5]byte
}

// encrypt performs the AES encryption specified by the protocol
func (p *packet) encrypt(b cipher.Block) []byte {
	seqBytes := [2]byte{}
	binary.LittleEndian.PutUint16(seqBytes[:], p.seqNum)
	padLen := padLength(len(p.data))
	padding := make([]byte, padLen)
	plainWithPad := append(p.data, padding...)
	ciphertext := make([]byte, len(plainWithPad))
	for i := 0; i < len(ciphertext); i += blockSize {
		end := i + blockSize
		copy(ciphertext[i:end], plainWithPad[i:end])
		// The protocol requires XOR the first two bytes of the block with the sequence number.
		// Although not stated, this seems to avoid ECB mode issues by ensuring the same
		// message doesn't generate the same cipher.
		ciphertext[i] ^= seqBytes[0]
		ciphertext[i+1] ^= seqBytes[1]
		b.Encrypt(ciphertext[i:end], ciphertext[i:end])
	}

	return ciphertext
}

// decrypt performs the AES decryption specified by the protocol
func (p *packet) decrypt(b cipher.Block) []byte {
	seqBytes := [2]byte{}
	binary.LittleEndian.PutUint16(seqBytes[:], p.seqNum)
	plaintext := make([]byte, len(p.data))
	for i := 0; i < len(p.data); i += blockSize {
		end := i + blockSize
		b.Decrypt(plaintext[i:end], p.data[i:end])
		plaintext[i] ^= seqBytes[0]
		plaintext[i+1] ^= seqBytes[1]
	}
	return plaintext
}

// serialize turns a packet into its byte representation suitable for sending over the wire. If the block cipher is not nil the serialized packet will be encrypted.
func (p *packet) serialize(cb cipher.Block) []byte {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, p.seqNum)
	binary.Write(buf, binary.LittleEndian, p.msgType)
	binary.Write(buf, binary.LittleEndian, p.reserved)
	data := p.data
	if cb != nil {
		data = p.encrypt(cb)
	}
	binary.Write(buf, binary.LittleEndian, data)
	return buf.Bytes()
}

// unmarshal reads message data into the given interface via encoding/binary.
func (p *packet) unmarshal(v interface{}) error {
	r := bytes.NewBuffer(p.data)
	return binary.Read(r, binary.LittleEndian, v)
}

func padLength(unencryptedLength int) int {
	extra := unencryptedLength % blockSize
	padLen := 0
	if extra > 0 {
		padLen = blockSize - extra
	}
	return padLen
}
