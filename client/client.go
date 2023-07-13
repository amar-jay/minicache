package client

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/amar-jay/minicache/proto"
)

type Client struct {
	lock sync.Mutex
	conn net.Conn
}

// create a new client from an already existing connection
func NewWithConn(conn net.Conn) *Client {

	return &Client{
		conn: conn,
		lock: sync.Mutex{},
	}
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Set(ctx context.Context, key []byte, val []byte, ttl int64) error {

	c.lock.Lock()
	defer c.lock.Unlock()
	cmd := &proto.Set{
		Key: key,
		Val: val,
		TTL: ttl,
	}

	b, _ := cmd.Bytes()
	_, err := c.conn.Write(b)
	if err != nil {
		return err
	}

	resp, err := proto.ParseSetResponse(c.conn)
	if err != nil {
		return err
	}
	if resp.Status != http.StatusOK {
		return fmt.Errorf("server responsed with non OK status [%s]", resp.Status)
	}

	return nil
}
