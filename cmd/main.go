package main

import (
	"github.com/juancolamendy/tm-chat/conf"
	"github.com/juancolamendy/tm-chat/server"
)

func main() {
	config := conf.GetConfig()
	config.Dump()

	opts := &server.Opts {
		Port: config.Port,
		Host: config.Host,
	}
	svr := server.NewChatServer(opts)
	svr.Init()
}