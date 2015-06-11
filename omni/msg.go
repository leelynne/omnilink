package omni

import (
	"bytes"
	"encoding/binary"
	"io"
)

type msgType uint8

const (
	NoMsg = iota
	ClientReqNewSession
	ControllerAckNewSession
	ClientReqSecureConnection
	ControllerAckSecureConnection
	ClientSessionTerminated
	ControllerSessionTerminated
	ControllerCannotStartNewSession
	AppDataMsg
)

// genmsg is the generic message format (non-application)
type genmsg struct {
	SeqNum   uint16  // Sequence Number From 1 - 65535. Zero means no sequence tracking
	Type     msgType // Type of this message. Specifies the format of the MsgData
	reserved byte    // Not used
	Data     []byte  // Application message data
}

type contAckNewSession struct {
	ProtoVersion uint16 // HAI Network protocol version used by the controller
	SessionID    []byte // 40 bit session ID
}

func (m *genmsg) serialize() []byte {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, m.SeqNum)
	binary.Write(buf, binary.LittleEndian, m.Type)
	binary.Write(buf, binary.LittleEndian, m.reserved)
	binary.Write(buf, binary.LittleEndian, m.Data)
	return buf.Bytes()
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
