package conf

import (
	"flag"
	"log"
)

type Config struct {
	Port    string
	Host    string
	Server  bool
}

func GetConfig() *Config {
	config := &Config{}

	flag.StringVar(&config.Host, "host", "localhost", "host")
	flag.StringVar(&config.Port, "port", "9090", "port")
	flag.BoolVar(&config.Server, "server", false, "server mode")
	flag.Parse()

	return config
}

func (c *Config) Dump() {
	log.Printf("config: Host: %s", c.Host)
	log.Printf("config: Port: %s", c.Port)
	log.Printf("config: Server: %t", c.Server)
}