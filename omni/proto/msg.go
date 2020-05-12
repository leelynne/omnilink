package proto

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
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
	MsgReqObjectStatus         AppMsgType = 0x22
)

// Msg is the raw application data message
type Msg struct {
	Type AppMsgType
	Data []byte
}

// NewMsg creates a Msg from a packet received by a connection.
func NewMsg(p *packet) (*Msg, error) {
	m := &Msg{}
	buf := bytes.NewBuffer(p.data)
	var start [1]byte
	binary.Read(buf, binary.LittleEndian, &start)

	var crc, expectedCRC uint16
	var dataLen uint8
	binary.Read(buf, binary.LittleEndian, &dataLen)
	crc = updateCRC(crc, dataLen)
	binary.Read(buf, binary.LittleEndian, &m.Type)
	crc = updateCRC(crc, byte(m.Type))
	m.Data = make([]byte, dataLen-1, dataLen-1) // Data Length includes one byte for the Type field.
	err := binary.Read(buf, binary.LittleEndian, m.Data)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to read Data from packet")
	}
	crc = updateCRC(crc, m.Data...)

	err = binary.Read(buf, binary.LittleEndian, &expectedCRC)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to read CRC from packet")
	}

	if crc != expectedCRC {
		return nil, errors.Errorf("CRC mismatch on received packet")
	}
	return m, nil
}

// Reader returns an io.Reader from the underlying message data.
func (m *Msg) Reader() io.Reader {
	return bytes.NewBuffer(m.Data)
}

// Packet creates a packet suitable for sending over a connection.
func (m *Msg) packet(seqNum uint16) *packet {
	plaintext := m.serialize()
	p := packet{
		seqNum:  seqNum,
		msgType: msgAppData,
		data:    plaintext,
	}
	return &p
}

func (m *Msg) serialize() []byte {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, appMsgStart)
	binary.Write(buf, binary.LittleEndian, uint8(len(m.Data)+1))
	binary.Write(buf, binary.LittleEndian, m.Type)
	binary.Write(buf, binary.LittleEndian, m.Data)
	binary.Write(buf, binary.LittleEndian, m.crc())
	out := buf.Bytes()
	return out
}

func (m *Msg) crc() uint16 {
	crc := updateCRC(0, uint8(len(m.Data)+1))
	crc = updateCRC(crc, byte(m.Type))
	return updateCRC(crc, m.Data...)
}

func updateCRC(crc uint16, in ...byte) uint16 {
	for _, b := range in {
		crc ^= uint16(b)
		for i := 0; i < 8; i++ {
			flag := (crc & 1) == 1
			crc = crc >> 1
			if flag {
				crc ^= haiPoly
			}
		}
	}
	return crc
}
