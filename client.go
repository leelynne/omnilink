package main

import (
	"net"
	"sync"
	"time"
)

// Client is an Omni-link II client
type Client struct {
	Addr string // IP:Port
	mu   sync.Mutex
	conn net.Conn
}

func (c *Client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.Addr, time.Duration(10*time.Second))
	if err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.conn = conn
	return err
}

func (c *Client) Receive() (*msg, error) {
	buf := make([]byte, 20)
	_, err := c.conn.Read(buf)
	if err != nil {
		return nil, err
	}

	return deserialize(buf)
}

func (c *Client) Send(m msg) error {
	_, err := c.conn.Write(m.serialize())
	return err
}
