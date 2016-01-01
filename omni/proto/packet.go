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
type ackSecureSession struct {
	SessionID [5]byte
}

// encrypt performs the AES encryption specified by the protocol
func (p packet) encrypt(b cipher.Block) []byte {
	seqBytes := [2]byte{}
	binary.LittleEndian.PutUint16(seqBytes[:], p.seqNum)
	padLen := padLength(len(p.data))
	padding := make([]byte, padLen)
	plainWithPad := append(p.data, padding...)
	ciphertext := make([]byte, len(plainWithPad))
	for i := 0; i < len(ciphertext); i += blockSize {
		end := i + blockSize
		copy(ciphertext[i:end], plainWithPad[i:end])
		ciphertext[i] ^= seqBytes[0]
		ciphertext[i+1] ^= seqBytes[1]
		b.Encrypt(ciphertext[i:end], ciphertext[i:end])
	}

	return ciphertext
}

func padLength(unencryptedLength int) int {
	extra := unencryptedLength % blockSize
	padLen := 0
	if extra > 0 {
		padLen = blockSize - extra
	}
	return padLen
}

// decrypt performs the AES decryption specified by the protocol
func (p packet) decrypt(b cipher.Block) []byte {
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

func (p packet) deserialize(v interface{}) error {
	r := bytes.NewBuffer(p.data)
	return binary.Read(r, binary.LittleEndian, v)
}
