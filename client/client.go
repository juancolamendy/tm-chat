package client

import (
	"net"
	"bufio"
	"fmt"
	"time"
)

type ChatClient struct {
	port string
	host string
	conn net.Conn
	
	Bufin *bufio.Reader
	Bufout *bufio.Writer
	Address string
	Ts int64	
}

type Opts struct {
	Port string
	Host string
}

func NewChatClient(opts *Opts) *ChatClient {
	return &ChatClient{
		port: opts.Port,
		host: opts.Host,
	}	
}

func (c *ChatClient) Open() error {
	address := fmt.Sprintf("%s:%s", c.host, c.port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err	
	}

	c.conn = conn
	c.Bufin = bufio.NewReader(conn)
	c.Bufout = bufio.NewWriter(conn)
	c.Address = conn.RemoteAddr().String()
	c.Ts = time.Now().Unix()

	return nil
}

func (c *ChatClient) Close() {
	c.conn.Close()
}