package main

import (
	"log"

	"github.com/juancolamendy/tm-chat/conf"
	"github.com/juancolamendy/tm-chat/server"
	"github.com/juancolamendy/tm-chat/client"
	"github.com/juancolamendy/tm-chat/shell"
)

func main() {
	config := conf.GetConfig()
	config.Dump()

	if config.Server {
		// If server mode, init a server
		opts := &server.Opts {
			Port: config.Port,
			Host: config.Host,
		}
		svr := server.NewChatServer(opts)
		svr.Init()
	} else {
		// If client mode, init a client and shell
		opts := &client.Opts {
			Port: config.Port,
			Host: config.Host,
		}
		client := client.NewChatClient(opts)
		err := client.Init()
		if err != nil {
			log.Fatalf("app - error init client - %v", err)
		}

		// launch shell
		shell := shell.NewShell(client)
		shell.Init()
	}
}