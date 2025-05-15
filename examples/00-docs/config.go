package main

import "time"

type Config struct {
	Apikey   string
	Loglevel string
	Server   ServerConfig
}

type ServerConfig struct {
	Host             string
	Port             int
	ShutdownInterval time.Duration
}
