package client

import (
	"log"
	"net"
	"bufio"
	"fmt"
	"time"
)

type ClientEventType int32

const (	
	ClientEventType_StopProcessing  ClientEventType = 1
	ClientEventType_Message         ClientEventType = 2
)

type ClientEvent struct {
	ClientEventType ClientEventType
	Payload        interface{}
	ResCh          chan interface{}
}

type ChatClient struct {
	port string
	host string
	conn net.Conn
	
	Bufin *bufio.Reader
	Bufout *bufio.Writer
	Address string
	Ts int64

	InChan chan *ClientEvent
	OutChan chan *ClientEvent
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
		log.Printf("client - error on opening - %+v", err)
		return err
	}

	// goroutine - read from socket and write/pipe to OutChan
	go func(outChan chan *ClientEvent) {
		defer func() {
			// recover from panic and releasing resources
			if r := recover(); r != nil {
				log.Printf("client - writeToOutChan recovered - error: %+v\n", r)
			}
			c.closeChan <- true
		}()

		for {
			text, err := c.Bufin.ReadString('\n')
			if err != nil {
				log.Printf("client - error reading %+v", err)
				c.closeChan <- true
				return
			}
			log.Printf("client - received: %s", text)
			text = text[:len(text)-1]
			outChan <- &ClientEvent {
				ClientEventType: ClientEventType_Message,
				Payload: text,
			}
		}
	}(c.OutChan)

	// goroutine - read from InChan and write/pipe to the socket
	go func(inChan chan *ClientEvent) {
		defer func() {
			// recover from panic and releasing resources
			if r := recover(); r != nil {
				log.Printf("client - readFromInChan recovered - error: %+v\n", r)
			}
			c.closeChan <- true
		}()

		for {
			select {
			case evt := <- inChan:
				if evt.ClientEventType != ClientEventType_Message {
					continue
				}
				text, ok := evt.Payload.(string)
				if !ok {
					continue
				}

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
			log.Printf("client - closing")
			c.close()
			c.OutChan <- &ClientEvent {
				ClientEventType: ClientEventType_StopProcessing,
			}			
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
	c.InChan = make(chan *ClientEvent)
	c.OutChan = make(chan *ClientEvent)
	c.closeChan = make(chan bool)

	return nil
}

func (c *ChatClient) close() {
	c.conn.Close()
}