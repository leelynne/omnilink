package omni

import (
	"encoding/hex"
	"fmt"
	"strings"

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

/*
func (c *Client) GetSystemInformation() (SystemInfo, error) {
}


func (c *Client) GetSystemInformation() (SystemInfo, error) {
	seqNum := c.nextSeqNum()
	data := ReqSystemInfoMsg.serialize(c, seqNum)
	err := c.Send(genmsg{
		SeqNum: seqNum,
		Type:   AppDataMsg,
		Data:   data,
	})
	if err != nil {
		return SystemInfo{}, fmt.Errorf("Faile dto send - %s", err.Error())
	}
	msg, err := c.ReceiveMsg()
	if err != nil {
		return SystemInfo{}, fmt.Errorf("Failed to receive system info %s", err.Error())
	}
	fmt.Printf("sysinfo %+v\n", msg)
	buf := bytes.NewBuffer(msg.Data)
	si := SystemInfo{}

	err = binary.Read(buf, binary.LittleEndian, &si.ModelNumber)
	err = binary.Read(buf, binary.LittleEndian, &si.MajorVersion)
	err = binary.Read(buf, binary.LittleEndian, &si.MinorVerison)
	err = binary.Read(buf, binary.LittleEndian, &si.Revesion)
	err = binary.Read(buf, binary.LittleEndian, &si.LocalPhoneNumber)

	return si, err
}

func (c *Client) GetSystemStatus() (SystemStatus, error) {
	seqNum := c.nextSeqNum()
	data := msgReqSystemStatus.serialize(c, seqNum)
	err := c.Send(genmsg{
		SeqNum: seqNum,
		Type:   AppDataMsg,
		Data:   data,
	})
	if err != nil {
		return SystemStatus{}, fmt.Errorf("Faile dto send - %s", err.Error())
	}
	msg, err := c.ReceiveMsg()
	if err != nil {
		return SystemStatus{}, fmt.Errorf("Failed to receive system info %s", err.Error())
	}
	fmt.Printf("sysinfo %+v\n", msg)
	buf := bytes.NewBuffer(msg.Data)
	si := SystemStatus{}

	err = binary.Read(buf, binary.LittleEndian, &si.DateValid)
	err = binary.Read(buf, binary.LittleEndian, &si.Year)
	err = binary.Read(buf, binary.LittleEndian, &si.Month)
	err = binary.Read(buf, binary.LittleEndian, &si.Day)
	err = binary.Read(buf, binary.LittleEndian, &si.DayOfWeek)
	err = binary.Read(buf, binary.LittleEndian, &si.Hour)
	err = binary.Read(buf, binary.LittleEndian, &si.Minute)
	err = binary.Read(buf, binary.LittleEndian, &si.Second)

	return si, nil
}
*/
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
