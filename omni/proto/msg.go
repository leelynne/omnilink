package proto

import (
	"bytes"
	"encoding/binary"
	"sync"
)

type AppMsgType uint8

const (
	appMsgStart                byte       = 0x21
	haiPoly                    uint16     = 0xA001
	MsgReqSystemInfo           AppMsgType = 0x16
	MsgReqSystemStatus         AppMsgType = 0x18
	MsgReqSystemTroubles       AppMsgType = 0x1A
	MsgReqSystemFeatures       AppMsgType = 0x1C
	MsgReqSystemFormats        AppMsgType = 0x28
	MsgReqObjectTypeCapacities AppMsgType = 0x1E
	MsgReqObjectProperties     AppMsgType = 0x20
)

// Msg is the application data message
type Msg struct {
	Type     AppMsgType
	Data     []byte
	crcRecvd []byte
	mu       sync.Mutex
	buf      *bytes.Buffer
}

func (m *Msg) Read(p []byte) (n int, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.buf == nil {
		m.buf = bytes.NewBuffer(m.Data)
	}
	return m.buf.Read(p)
}

func (m *Msg) serialize() []byte {
	m.mu.Lock()
	defer m.mu.Unlock()
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, appMsgStart)
	binary.Write(buf, binary.LittleEndian, uint8(len(m.Data)+1))
	binary.Write(buf, binary.LittleEndian, m.Type)
	binary.Write(buf, binary.LittleEndian, m.Data)
	binary.Write(buf, binary.LittleEndian, m.crc())
	out := buf.Bytes()
	return out
}

func (m *Msg) crc() []byte {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, uint8(len(m.Data)+1))
	binary.Write(buf, binary.LittleEndian, m.Type)
	binary.Write(buf, binary.LittleEndian, m.Data)
	var crc uint16
	for _, b := range buf.Bytes() {
		crc ^= uint16(b)
		for i := 0; i < 8; i++ {
			flag := (crc & 1) == 1
			crc = crc >> 1
			if flag {
				crc ^= haiPoly
			}
		}
	}
	out := make([]byte, 2)
	binary.LittleEndian.PutUint16(out, crc)
	return out
}

func fromPacket(p packet) *Msg {
	m := &Msg{}
	buf := bytes.NewBuffer(p.data)
	var start [1]byte
	binary.Read(buf, binary.LittleEndian, &start)
	var len uint8
	binary.Read(buf, binary.LittleEndian, &len)
	binary.Read(buf, binary.LittleEndian, &m.Type)
	rest := buf.Bytes()
	m.Data = rest[0:len]
	m.crcRecvd = rest[len : len+2]
	return m
}
