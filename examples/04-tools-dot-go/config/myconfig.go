package config

import "time"

//go:generate go run github.com/Ozoniuss/genconfig -struct=MyConfig -project=Myapp
type MyConfig struct {
	HddSyncPath string
	DryRun, Lol bool
	Timeout     time.Duration
	Port        int
	Port32      uint32
	Port16      int16

	Ne Nested
}

type Nested struct {
	Name string
	Age  int
}

// does not get covered yet
type Nested2 = Nested
