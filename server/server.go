package server

import (
	"fmt"
	"log"
	"net"
	"time"
	"bufio"
)

type ConnHandler func (bufin *bufio.Reader, bufout *bufio.Writer, address string, ts int64)

type ChatServer struct {
	port        string
	host        string
	connHandler ConnHandler
}

type Opts struct {
	Port        string
	Host        string
	ConnHandler ConnHandler
}

func NewChatServer(opts *Opts) *ChatServer {
	return &ChatServer{
		port:        opts.Port,
		host:        opts.Host,
		connHandler: opts.ConnHandler,
	}
}

func (s *ChatServer) Init() {	
	address := fmt.Sprintf("%s:%s", s.host, s.port)
	log.Printf("initing server on address %s", address)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("error listening - %+v",err)
	}
	log.Printf("server initiated on address %s", address)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("conection error - %+v", err)
			continue
		}
		go s.handleConn(conn, bufio.NewReader(conn), bufio.NewWriter(conn), conn.RemoteAddr().String(), time.Now().Unix())
	}	
}

func (s *ChatServer) handleConn(c net.Conn, bufin *bufio.Reader, bufout *bufio.Writer, address string, ts int64) {
	defer c.Close()

	if s.connHandler != nil {
		s.connHandler(bufin, bufout, address, ts)
	}
}