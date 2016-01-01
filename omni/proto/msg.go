package proto

import (
	"bytes"
	"encoding/binary"
)

type AppMsgType uint8

const (
	appMsgStart = 0x21
)

// Msg is the application data message
type Msg struct {
	Type AppMsgType
	Data []byte
}

func (m *Msg) serialize() []byte {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, appMsgStart)
	binary.Write(buf, binary.LittleEndian, len(m.Data)+2)
	binary.Write(buf, binary.LittleEndian, m.Type)
	binary.Write(buf, binary.LittleEndian, m.Data)
	binary.Write(buf, binary.LittleEndian, m.crc())
	return buf.Bytes()
}

func (m *Msg) crc() []byte {
	return nil
}

func fromPacket(p packet) *Msg {
	m := &Msg{}
	buf := bytes.NewBuffer(p.data)
	binary.Read(buf, binary.LittleEndian, &m.Type)
	copy(m.Data, buf.Bytes())
	return m
}
