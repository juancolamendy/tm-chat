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

	flag.StringVar(&config.Host, "host", "localhost", "Server host")
	flag.StringVar(&config.Port, "port", "9090", "Server port")
	flag.BoolVar(&config.Server, "server", false, "Is only server")
	flag.Parse()

	return config
}

func (c *Config) Dump() {
	log.Printf("Config: Host: %s", c.Host)
	log.Printf("Config: Port: %s", c.Port)
	log.Printf("Config: Server: %t", c.Server)
}