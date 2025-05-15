package main

import (
	genconfig "github.com/Ozoniuss/genconfig/internal"
)

func main() {
	genconfig.GenerateConfigLoader("App", "Config", "config.go", "config_gen.go", ".env", "", true)
}
