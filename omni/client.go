package omni

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"strings"
	"time"
)

type sessionkey []byte

// Client is an Omni-link II client
type Client struct {
	Addr         string // IP:Port
	conn         net.Conn
	protoVersion uint16 // Protocol version used by the controller
	sessionID    []byte
	sessionKey   sessionkey
}

func NewClient(addr string, key string) (*Client, error) {
	conn, err := net.DialTimeout("tcp", addr, time.Duration(10*time.Second))
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected!")
	client := Client{
		Addr: addr,
		conn: conn,
	}
	m := genmsg{
		SeqNum: 1,
		Type:   ClientReqNewSession,
	}
	err = client.Send(m)
	if err != nil {
		return nil, fmt.Errorf("Failed to send %s", err.Error())
	}

	ackMsg, err := client.Receive()
	if err != nil {
		return nil, fmt.Errorf("Failed to receive %s", err.Error())
	}
	fmt.Printf("Received msg %+v\n", *ackMsg)
	buf := bytes.NewReader(ackMsg.Data)

	err = binary.Read(buf, binary.LittleEndian, &client.protoVersion)
	if err != nil {
		return nil, fmt.Errorf("Failed to read version %s", err.Error())
	}
	idbuf := bytes.Buffer{}
	_, err = idbuf.ReadFrom(buf)
	if err != nil {
		return nil, fmt.Errorf("Failed to read sessionID  %s", err.Error())
	}
	client.sessionID = idbuf.Bytes()
	client.sessionKey, err = createSessionKey(key, client.sessionID)
	if err != nil {
		return nil, fmt.Errorf("Failed to create session key - %s", err)
	}
	fmt.Printf("New session %+v\n", client)

	return &client, err
}

func createSessionKey(key string, sessionID []byte) (sessionkey, error) {
	keyb, err := parseKey(key)
	if err != nil {
		return sessionkey{}, err
	}
	skey := make([]byte, 16)
	i := 0
	for ; i <= 11; i++ {
		skey[i] = keyb[i]
	}
	for j := i; j < 16; j++ {
		skey[j] = keyb[j] ^ sessionID[j-12]
	}
	return skey, nil
}

func parseKey(key string) (sessionkey, error) {
	parts := strings.Split(key, "-")
	hexOnly := strings.Join(parts, "")

	keyBytes, err := hex.DecodeString(hexOnly)
	if err != nil {
		return sessionkey{}, err
	}
	if len(keyBytes) != 16 {
		return sessionkey{}, fmt.Errorf("Key %s must be 16 bytes long", key)
	}
	return sessionkey(keyBytes), nil
}

func (c *Client) Receive() (*genmsg, error) {
	w := bytes.NewBuffer([]byte{})
	buf := make([]byte, 1024)
	for {
		n, err := c.conn.Read(buf)
		if err != nil {
			return nil, err
		}
		w.Write(buf[0:n])
		if n < 1024 {
			break
		}
	}

	return deserialize(w)
}

func (c *Client) Send(m genmsg) error {
	_, err := c.conn.Write(m.serialize())
	return err
}
