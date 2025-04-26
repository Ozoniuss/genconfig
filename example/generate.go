package main

import (
	"testmain/github/genconfig"
)

func main() {
	genconfig.GenerateConfigLoader("", "Config", "config.go", "config_gen.go", "", true, "", true)
}
