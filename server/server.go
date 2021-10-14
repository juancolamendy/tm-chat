package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"	
)

type ChatServer struct {
	port string
	host string
}

type Opts struct {
	Port string
	Host string
}

func NewChatServer(opts *Opts) *ChatServer {
	return &ChatServer{
		port: opts.Port,
		host: opts.Host,
	}
}

func (s *ChatServer) Init() {
	address := fmt.Sprintf("%s:%s", s.host, s.port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		handleConn(conn)
	}	
}

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
		if err != nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
}