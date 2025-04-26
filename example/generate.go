package main

import (
	genconfig "testmain/github/genconfig/internal"
)

func main() {
	genconfig.GenerateConfigLoader("", "Config", "config.go", "config_gen.go", "", true, "", true)
}
