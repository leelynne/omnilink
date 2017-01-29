package omni

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/leelynne/omnilink/omni/proto"
)

// Client is an Omni-link II client
type Client struct {
	Addr string // IP:Port
	conn proto.Conn
}

func NewClient(addr string, key string) (*Client, error) {
	skey, err := parseKey(key)
	if err != nil {
		return nil, err
	}
	conn, err := proto.NewConnection(addr, skey)
	if err != nil {
		return nil, err
	}

	return &Client{
		Addr: addr,
		conn: conn,
	}, nil
}

type SystemInfo struct {
	ModelNumber      uint8
	MajorVersion     uint8
	MinorVerison     uint8
	Revesion         uint8
	LocalPhoneNumber [25]byte
}

type SystemStatus struct {
	DateValid   uint8
	Year        uint8
	Month       uint8
	Day         uint8
	DayOfWeek   uint8
	Hour        uint8
	Minute      uint8
	Second      uint8
	Daylight    uint8
	SunriseHour uint8
	SunriseMin  uint8
	SunsetHour  uint8
	SunsetMin   uint8
	Battery     uint8
}

func (c *Client) GetSystemInformation() (SystemInfo, error) {
	si := SystemInfo{}
	m := &proto.Msg{Type: proto.MsgReqSystemInfo}

	resp, err := c.get(m)
	if err != nil {
		return si, nil
	}

	err = binary.Read(resp, binary.LittleEndian, &si)
	return si, err
}

func (c *Client) get(m *proto.Msg) (*proto.Msg, error) {
	err := c.conn.Write(m, time.Second*10)
	if err != nil {
		return nil, err
	}
	return c.conn.Read(time.Second * 20)
}

func (c *Client) GetSystemStatus() (SystemStatus, error) {
	st := SystemStatus{}
	m := &proto.Msg{Type: proto.MsgReqSystemStatus}

	resp, err := c.get(m)
	if err != nil {
		return st, nil
	}
	err = binary.Read(resp, binary.LittleEndian, &st)
	return st, err
}

func parseKey(key string) (proto.StaticKey, error) {
	hexOnly := strings.Replace(key, "-", "", -1)
	keyBytes, err := hex.DecodeString(hexOnly)
	if err != nil {
		return nil, err
	}
	if len(keyBytes) != 16 {
		return nil, fmt.Errorf("Key %s must be 16 bytes long", key)
	}
	return proto.StaticKey(keyBytes), nil
}
