package main

import "time"

type Config struct {
	Apikey   string
	Loglevel string `default:"info"`
	Server   ServerConfig
}

type ServerConfig struct {
	Host             string
	Port             int
	ShutdownInterval time.Duration
}
