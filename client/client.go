package client

import (
	"log"
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

	InChan chan string
	OutChan chan string
	closeChan chan bool
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

func (c *ChatClient) Init() error {
	err := c.open()
	if err != nil {
		log.Printf("error on opening - %+v", err)
		return err
	}

	// goroutine - read from socket and write/pipe to OutChan
	go func(outChan chan string) {
		for {
			text, err := c.Bufin.ReadString('\n')
			if err != nil {
				log.Printf("client - error reading %+v", err)
				c.closeChan <- true
				return
			}
			log.Printf("client - received: %s", text)
			text = text[:len(text)-1]
			outChan <- text
		}
	}(c.OutChan)

	// goroutine - read from InChan and write/pipe to the socket
	go func(inChan chan string) {
		for {
			select {
			case text := <- inChan:
				log.Printf("client - sending: %s", text)
				_, err := c.Bufout.WriteString(fmt.Sprintf("%s\n",text))
				if err != nil {
					log.Printf("client - error writing %+v", err)
					c.closeChan <- true
					return
				}
				err = c.Bufout.Flush()
				if err != nil {
					log.Printf("client - error flushing %+v", err)
					c.closeChan <- true
					return
				}
			}
		}
	}(c.InChan)

	// goroutine - listen for closing socket
	go func() {
		select {
		case <- c.closeChan:
			c.close()
			return			
		}
	}()

	return nil
}

func (c *ChatClient) open() error {
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
	c.InChan = make(chan string)
	c.OutChan = make(chan string)
	c.closeChan = make(chan bool)

	return nil
}

func (c *ChatClient) close() {
	c.conn.Close()
}