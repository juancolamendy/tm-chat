package main

import (
	"log"
	"bufio"
	"time"

	"github.com/juancolamendy/tm-chat/conf"
	"github.com/juancolamendy/tm-chat/server"
	"github.com/juancolamendy/tm-chat/client"
)

func main() {
	config := conf.GetConfig()
	config.Dump()

	if config.Server {
		handler := func (bufin *bufio.Reader, bufout *bufio.Writer, address string, ts int64) {
			defer func() {
				log.Printf("handler - leaving the handler for %s", address)
			}()
			log.Printf("handler - conntected from %s", address)
			for {
				// write
				n, err := bufout.WriteString("server time:" + time.Now().Format("15:04:05\n"))
				if err != nil {
					log.Printf("error writing %+v", err)
					return
				}
				err = bufout.Flush()
				if err != nil {
					log.Printf("error writing %+v", err)
					return
				}
				log.Printf("written %d bytes", n)

				// read
				text, err := bufin.ReadString('\n')
				if err != nil {
					log.Printf("error reading %+v", err)
					return
				}
				log.Printf("read text: %s", text)

				// sleep
				time.Sleep(5 * time.Second)
			}
		}

		opts := &server.Opts {
			Port: config.Port,
			Host: config.Host,
			ConnHandler: handler,
		}
		svr := server.NewChatServer(opts)
		svr.Init()
	} else {
		opts := &client.Opts {
			Port: config.Port,
			Host: config.Host,
		}
		client := client.NewChatClient(opts)
		err := client.Open()
		if err != nil {
			log.Printf("error on opening - %+v", err)
			return
		}
		defer client.Close()

		for {
			// read
			text, err := client.Bufin.ReadString('\n')
			if err != nil {
				log.Printf("error reading %+v", err)
				return
			}
			log.Printf("read text: %s", text)

			// write
			n, err := client.Bufout.WriteString("client time:" + time.Now().Format("15:04:05\n"))
			if err != nil {
				log.Printf("error writing %+v", err)
				return
			}
			err = client.Bufout.Flush()
			if err != nil {
				log.Printf("error writing %+v", err)
				return
			}
			log.Printf("written %d bytes", n)			
		}
	}
}