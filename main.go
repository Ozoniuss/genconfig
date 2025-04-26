package main

import (
	"fmt"
	"os"

	gncfg "github.com/Ozoniuss/genconfig/internal"
)

const (
	defaultInputFile          = "config.go"
	defaultOutputDotenv       = ".env"
	defaultOutputConfigLoader = "config_gen.go"
)

func main() {
	err := gncfg.GenerateConfigLoader("", "Config", defaultInputFile, defaultOutputConfigLoader, defaultOutputDotenv, "", false)
	if err != nil {
		fmt.Printf("failed to generate config: %v", err.Error())
		os.Exit(1)
	}
}
