package client

import (
	"net"
	"sync"
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
