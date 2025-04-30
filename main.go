package main

import (
	"flag"
	"fmt"
	"os"

	gncfg "github.com/Ozoniuss/genconfig/internal"
)

const (
	defaultInputFile          = "config.go"
	defaultOutputDotenv       = ""
	defaultOutputConfigLoader = "config_gen.go"
	defaultConfigStructName   = "Config"
	defaultProjectName        = ""
)

var (
	flagOutputDotenvFile string
	flagConfigStructName string
	flagPrintUsage       bool
	flagProjectName      string
)

func main() {
	flag.StringVar(&flagProjectName, "project", defaultProjectName, "Name of the project. If not empty, will be used as prefix for environment variables.")
	flag.StringVar(&flagConfigStructName, "struct", defaultConfigStructName, "Name of the config struct.")
	flag.StringVar(&flagOutputDotenvFile, "env", defaultOutputDotenv, "Name of the output .env file, if you want to generate one with all possible config values. An empty value will not generate a .env file.")
	flag.BoolVar(&flagPrintUsage, "help", false, "Show usage.")

	flag.Parse()
	if flagPrintUsage {
		flag.Usage()
		os.Exit(0)
	}

	if flagConfigStructName == "" {
		fmt.Fprint(os.Stderr, "Please specify a non-empty config struct.")
		os.Exit(1)
	}

	err := gncfg.GenerateConfigLoader(flagProjectName, flagConfigStructName, defaultInputFile, defaultOutputConfigLoader, defaultOutputDotenv, flagOutputDotenvFile, false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not generate config: %s", err.Error())
		os.Exit(1)
	}
}
