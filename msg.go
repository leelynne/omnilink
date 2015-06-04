package main

import (
	"bytes"
	"encoding/binary"
)

type msgType byte

const (
	NO_MSG_ = iota
	CLIENT_REQ_NEW_SESSION
	CONTROLLER_ACK_NEW_SESSION
	CLIENT_REQ_SECURE_CONNECTION
	CONTROLLER_ACK_SECURE_CONNECTION
	CLIENT_SESSION_TERMINATED
	CONTROLLER_SESSION_TERMINATED
	CONTROLLER_CANNOT_START_NEW_SESSION
	APP_DATA_MSG
)

type msg struct {
	SeqNum   uint16
	MsgType  byte
	reserved byte
	MsgData  []byte
}

func (m *msg) serialize() []byte {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, m.SeqNum)
	binary.Write(buf, binary.LittleEndian, m.MsgType)
	binary.Write(buf, binary.LittleEndian, m.reserved)
	binary.Write(buf, binary.LittleEndian, m.MsgData)
	return buf.Bytes()
}

func deserialize(b []byte) (*msg, error) {
	m := msg{}
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.LittleEndian, &m.SeqNum)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buf, binary.LittleEndian, &m.MsgType)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buf, binary.LittleEndian, &m.reserved)
	if err != nil {
		return nil, err
	}
	m.MsgData = buf.ReadAt(b, 4)
	_, err = buf.Read(m.MsgData)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
