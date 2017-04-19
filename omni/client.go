package omni

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/leelynne/omnilink/omni/proto"
	"github.com/pkg/errors"
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

func (c *Client) GetSystemTroubles() (SystemTroubles, error) {
	m := &proto.Msg{Type: proto.MsgReqSystemTroubles}

	resp, err := c.get(m)
	if err != nil {
		return SystemTroubles{}, errors.Wrap(err, "Failed to get system troubles")
	}

	numTroubles := len(resp.Data) - 1
	troubles := make([]SystemTrouble, numTroubles)
	fmt.Printf("Type %x\n", resp.Type)
	for i := range troubles {
		troubles[i] = SystemTrouble(resp.Data[i])
	}
	return SystemTroubles{
		Troubles: troubles,
	}, nil
}

func (c *Client) GetSystemFeatures() (SystemFeatures, error) {
	m := &proto.Msg{Type: proto.MsgReqSystemFeatures}

	resp, err := c.get(m)
	if err != nil {
		return SystemFeatures{}, errors.Wrap(err, "Failed to get system features")
	}

	numFeatures := len(resp.Data) - 1
	features := make([]SystemFeature, numFeatures)
	fmt.Printf("Type %x\n", resp.Type)
	for i := range features {
		features[i] = SystemFeature(resp.Data[i])
	}
	return SystemFeatures{
		Features: features,
	}, nil
}

func (c *Client) GetSystemFormats() (SystemFormats, error) {
	m := &proto.Msg{Type: proto.MsgReqSystemFormats}

	resp, err := c.get(m)
	if err != nil {
		return SystemFormats{}, errors.Wrap(err, "Failed to get system formats")
	}
	fmt.Printf("Type %x\n", resp.Type)
	sf := SystemFormats{}
	err = binary.Read(resp, binary.LittleEndian, &sf)
	return sf, err
}

func (c *Client) get(m *proto.Msg) (*proto.Msg, error) {
	err := c.conn.Write(m, time.Second*10)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to write")
	}
	return c.conn.Read(time.Second * 20)
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
