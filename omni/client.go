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

// NewClient returns a Client
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

	resp, err := c.sendMessage(m)
	if err != nil {
		return si, nil
	}

	err = unmarshalMessage(resp, &si)
	return si, err
}

func (c *Client) GetSystemStatus() (SystemStatus, error) {
	st := SystemStatus{}
	m := &proto.Msg{Type: proto.MsgReqSystemStatus}

	resp, err := c.sendMessage(m)
	if err != nil {
		return st, nil
	}

	err = unmarshalMessage(resp, &st)
	return st, err
}

func (c *Client) GetSystemTroubles() (SystemTroubles, error) {
	m := &proto.Msg{Type: proto.MsgReqSystemTroubles}

	resp, err := c.sendMessage(m)
	if err != nil {
		return SystemTroubles{}, errors.Wrap(err, "Failed to get system troubles")
	}

	numTroubles := len(resp.Data) - 1
	troubles := make([]SystemTrouble, numTroubles)
	for i := range troubles {
		troubles[i] = SystemTrouble(resp.Data[i])
	}
	return SystemTroubles{
		Troubles: troubles,
	}, nil
}

func (c *Client) GetSystemFeatures() (SystemFeatures, error) {
	m := &proto.Msg{Type: proto.MsgReqSystemFeatures}

	resp, err := c.sendMessage(m)
	if err != nil {
		return SystemFeatures{}, errors.Wrap(err, "Failed to get system features")
	}

	numFeatures := len(resp.Data) - 1

	features := make([]SystemFeature, numFeatures)
	for i := range features {
		features[i] = SystemFeature(resp.Data[i])
	}
	return SystemFeatures{
		Features: features,
	}, nil
}

func (c *Client) GetSystemFormats() (SystemFormats, error) {
	m := &proto.Msg{Type: proto.MsgReqSystemFormats}

	resp, err := c.sendMessage(m)
	if err != nil {
		return SystemFormats{}, errors.Wrap(err, "Failed to get system formats")
	}

	sf := SystemFormats{}
	err = unmarshalMessage(resp, &sf)
	return sf, err
}

func (c *Client) GetObjectTypeCapacity(t ObjectType) (ObjectTypeCapacities, error) {
	m := &proto.Msg{
		Type: proto.MsgReqObjectTypeCapacities,
		Data: []byte{byte(t)},
	}

	resp, err := c.sendMessage(m)
	if err != nil {
		return ObjectTypeCapacities{}, errors.Wrapf(err, "Failed to get object type capacity for type %s", t)
	}

	otc := ObjectTypeCapacities{}
	err = unmarshalMessage(resp, &otc)
	return otc, err
}

func (c *Client) GetObjectProperties(objectType ObjectType) (properties interface{}, numObject int, e error) {
	msgs := []*proto.Msg{}
	for lsb := 0; true; lsb++ {
		m := &proto.Msg{
			Type: proto.MsgReqObjectProperties,
			Data: []byte{
				byte(objectType),
				byte(0),
				byte(lsb),
				byte(1),
				byte(0),
				byte(255), // Area
				byte(0),
			},
		}

		resp, err := c.sendMessage(m)

		if err != nil {
			return ObjectProperties{}, 0, errors.Wrap(err, "Failed to get object property")
		}
		if len(resp.Data) <= 1 {
			break
		}
		msgs = append(msgs, resp)
	}
	var out interface{}
	switch objectType {
	case Thermostat:
		tprops := []ThermostatProperties{}
		for _, msg := range msgs {
			tprop := ThermostatProperties{}
			err := unmarshalMessage(msg, &tprop)
			if err != nil {
				return nil, 0, errors.Wrap(err, "Failed to marshal data into Property")
			}
			tprops = append(tprops, tprop)
		}
		out = tprops
		numObject = len(tprops)
	}

	return out, numObject, nil
}

func (c *Client) GetObjectStatus(objectType ObjectType, numObjects int) (interface{}, error) {
	m := &proto.Msg{
		Type: proto.MsgReqObjectStatus,
		Data: []byte{
			byte(objectType),
			byte(0),
			byte(1),
			byte(0),
			byte(numObjects - 1),
		},
	}

	resp, err := c.sendMessage(m)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get object status")
	}

	reader := resp.Reader()
	var respType uint8
	binary.Read(reader, binary.LittleEndian, &respType)
	if respType != uint8(objectType) {
		return nil, errors.Errorf("Wrong return typed '%d' for input type '%d'", respType, objectType)
	}

	statusSize := StatusSizes[objectType]
	total := len(resp.Data) / statusSize

	var out interface{}
	switch objectType {
	case Thermostat:
		otc := make([]ThermostatStatus, total)
		err = binary.Read(reader, binary.LittleEndian, &otc)
		out = otc
	}

	return out, err
}

// sendMessage sends an application data message to the controller and returns a response
func (c *Client) sendMessage(m *proto.Msg) (*proto.Msg, error) {
	err := c.conn.Write(m, time.Second*10)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to write")
	}
	return c.conn.Read(time.Second * 20)
}

// unmarshalMessage unpacks an application data message into a struct
func unmarshalMessage(msg *proto.Msg, data interface{}) error {
	reader := msg.Reader()
	return binary.Read(reader, binary.LittleEndian, data)
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
