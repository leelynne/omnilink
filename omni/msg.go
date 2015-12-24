package omni

import (
	"bytes"
	"crypto/cipher"
	"encoding/binary"
	"fmt"
	"io"
)

type msgType uint8
type appMsgType uint8

const (
	NoMsg = iota
	ClientReqNewSession
	ControllerAckNewSession
	ClientReqSecureConnection
	ControllerAckSecureConnection
	ClientSessionTerminated
	ControllerSessionTerminated
	ControllerCannotStartNewSession

	AppDataMsg = 32

	maxSeqNum = 65535

	appMsgStart        = 0x21
	ackMsgType         = 0x01
	negativeAckMsgType = 0x02
	eomMsgType         = 0x03
)

var AckMsg = appmsg{
	Start:    appMsgStart,
	Length:   0x01,
	Type:     ackMsgType,
	CRCLeast: 0xC0,
	CRCMost:  0x50,
}

var NegativeAckMsg = appmsg{
	Start:    appMsgStart,
	Length:   0x01,
	Type:     negativeAckMsgType,
	CRCLeast: 0x80,
	CRCMost:  0x51,
}

var ReqSystemInfoMsg = appmsg{
	Start:    appMsgStart,
	Length:   0x01,
	Type:     0x16,
	CRCLeast: 0x80,
	CRCMost:  0x5E,
}
var msgReqSystemStatus = appmsg{
	Start:    appMsgStart,
	Length:   0x01,
	Type:     0x18,
	CRCLeast: 0x01,
	CRCMost:  0x9A,
}

// genmsg is the generic message format (non-application)
type genmsg struct {
	SeqNum   uint16  // Sequence Number From 1 - 65535. Zero means no sequence tracking
	Type     msgType // Type of this message. Specifies the format of the MsgData
	reserved byte    // Not used
	Data     []byte  // Application message data
}

/* appmsg is the application level message that rides on top of the genmsg
The CRC value is calculated using all bytes of the message, except the “start character” and the CRC fields.
*/
type appmsg struct {
	Start    byte
	Length   byte
	Type     appMsgType
	Data     []byte
	CRCLeast byte // LSB of 16-bit CRC
	CRCMost  byte // MSB of 16-bit CRC
}

func (m *genmsg) serialize(c cipher.Block) []byte {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, m.SeqNum)
	binary.Write(buf, binary.LittleEndian, m.Type)
	binary.Write(buf, binary.LittleEndian, m.reserved)
	data := m.Data
	if c != nil {
		data = m.encrypt(c)
	}
	binary.Write(buf, binary.LittleEndian, data)
	return buf.Bytes()
}

func (m *genmsg) encrypt(b cipher.Block) []byte {
	seqBytes := [2]byte{}
	binary.LittleEndian.PutUint16(seqBytes[:], m.SeqNum)
	extra := len(m.Data) % 16
	padLen := 0
	if extra > 0 {
		padLen = 16 - extra
	}
	padding := make([]byte, padLen)
	plainWithPad := append(m.Data, padding...)
	fmt.Printf("Plain %v\n", plainWithPad)
	ciphertext := make([]byte, len(plainWithPad))
	for i := 0; i < len(ciphertext); i += 16 {
		end := i + 16
		copy(ciphertext[i:end], plainWithPad[i:end])
		ciphertext[i] ^= seqBytes[0]
		ciphertext[i+1] ^= seqBytes[1]
		b.Encrypt(ciphertext[i:end], ciphertext[i:end])
	}

	fmt.Printf("Encrypt result %v\n", ciphertext)
	return ciphertext
}

func (m *genmsg) decrypt(b cipher.Block) {
	fmt.Printf("DECRYPTING!")
	seqBytes := [2]byte{}
	binary.LittleEndian.PutUint16(seqBytes[:], m.SeqNum)
	for i := 0; i < len(m.Data); i += 16 {
		end := i + 16
		b.Decrypt(m.Data[i:end], m.Data[i:end])
		m.Data[i] ^= seqBytes[0]
		m.Data[i+1] ^= seqBytes[1]
	}
}

func deserialize(buf io.Reader) (*genmsg, error) {
	m := genmsg{}
	err := binary.Read(buf, binary.LittleEndian, &m.SeqNum)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buf, binary.LittleEndian, &m.Type)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buf, binary.LittleEndian, &m.reserved)
	if err != nil {
		return nil, err
	}
	w := &bytes.Buffer{}
	_, err = io.Copy(w, buf)

	if err != nil {
		return nil, err
	}
	m.Data = w.Bytes()
	return &m, nil
}

func (m *appmsg) serialize(c *Client, seqNum uint16) []byte {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, m.Start)
	binary.Write(buf, binary.LittleEndian, m.Length)
	binary.Write(buf, binary.LittleEndian, m.Type)
	binary.Write(buf, binary.LittleEndian, m.Data)
	binary.Write(buf, binary.LittleEndian, m.CRCLeast)
	binary.Write(buf, binary.LittleEndian, m.CRCMost)
	return buf.Bytes()
}

func deserializeAppMsg(c *Client, buf io.Reader) (*appmsg, error) {
	m := appmsg{}
	err := binary.Read(buf, binary.LittleEndian, &m.Start)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buf, binary.LittleEndian, &m.Length)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buf, binary.LittleEndian, &m.Type)
	if err != nil {
		return nil, err
	}
	w := &bytes.Buffer{}
	_, err = io.CopyN(w, buf, int64(m.Length))

	if err != nil {
		return nil, err
	}
	m.Data = w.Bytes()
	return &m, nil
}

type SystemInfo struct {
	ModelNumber      uint8
	MajorVersion     uint8
	MinorVerison     uint8
	Revesion         uint8
	LocalPhoneNumber [25]byte
}

type SystemStatus struct {
	DateValid uint8
	Year      uint8
	Month     uint8
	Day       uint8
	DayOfWeek uint8
	Hour      uint8
	Minute    uint8
	Second    uint8
}
