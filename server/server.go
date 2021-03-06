package server

import (
	"fmt"
	"log"
	"net"
	"time"
	"bufio"

	"github.com/juancolamendy/tm-chat/utils/ioutils"
)

type ChatServer struct {
	port        string
	host        string
}

type Opts struct {
	Port        string
	Host        string
}

func NewChatServer(opts *Opts) *ChatServer {
	return &ChatServer{
		port:        opts.Port,
		host:        opts.Host,
	}
}

func (s *ChatServer) Init() {
	address := fmt.Sprintf("%s:%s", s.host, s.port)
	log.Printf("server - initing server on address %s", address)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("server - error listening - %+v",err)
	}

	log.Printf("server - server initiated on address %s - accepting connections", address)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("server - conection error - %+v", err)
			continue
		}
		go s.handleConn(conn, bufio.NewReader(conn), bufio.NewWriter(conn), conn.RemoteAddr().String(), time.Now().Unix())
	}	
}

func (s *ChatServer) handleConn(c net.Conn, bufin *bufio.Reader, bufout *bufio.Writer, address string, ts int64) {
	defer func() {
		// recover from panic and releasing resources
		if r := recover(); r != nil {
			log.Printf("server - recovered - error: %+v\n", r)
		}
		c.Close()
	}()

	log.Printf("server - accepted connection from %s at %d", address, ts)
	for {
		// read text ending with \n
		text, err := ioutils.ReceiveString(bufin)
		if err != nil {
			log.Printf("server - error receiving %+v", err)
			return
		}
		
		log.Printf("server - received text: %s", text)
		cmd, err := parseToCommand(text)
		if err != nil {
			log.Printf("server - error parsing commands %+v", err)
			continue
		}
		log.Printf("server - command: %+v", cmd)

		// write text with \n
		out := fmt.Sprintf("server echo:%s",text)
		err = ioutils.SendString(bufout, out)
		if err != nil {
			log.Printf("server - error sending %+v", err)
			return
		}
		log.Printf("server - sent %s", out)
	}
}