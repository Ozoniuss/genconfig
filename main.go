package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	gncfg "github.com/Ozoniuss/genconfig/internal"
)

const (
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
	flagConfigFilePath   string
	flagDebugLogs        bool
)

func main() {

	flag.StringVar(&flagProjectName, "project", defaultProjectName, "Name of the project. If not empty, will be used as prefix for environment variables.")
	flag.StringVar(&flagConfigStructName, "struct", defaultConfigStructName, "Name of the config struct.")
	flag.StringVar(&flagConfigFilePath, "path", "", "File path of the config struct. Defaults to the location of the go:generate directive. Note that because of this, running genconfig as an executable without providing this flag can behave unpredictably.")
	flag.StringVar(&flagOutputDotenvFile, "env", defaultOutputDotenv, "Name of the output .env file, if you want to generate one with all possible config values. An empty value will not generate a .env file.")
	flag.BoolVar(&flagDebugLogs, "debug", false, "Debug generation issues.")
	flag.BoolVar(&flagPrintUsage, "help", false, "Show usage.")

	flag.Parse()
	if flagPrintUsage {
		flag.Usage()
		os.Exit(0)
	}

	// At the moment, all paths supplied are relative to where the executable
	// is called. This means if you have a certain parent for the input config,
	// you need the same parent for the output config. However, there is a
	// reasonable use case in which you may provide the path of the parent, but
	// don't want to provide a path for the output (the reverse doesn't really
	// make a reasonable use case). In that case, detect only the input path
	// was specified, and appent that path's parent to the output too.
	//
	// For now, the generated config name is not configurable, but will be
	// in the future.
	parent, _ := filepath.Split(flagConfigFilePath)
	outputConfigLoaderPath := parent + defaultOutputConfigLoader

	var configFilePath string
	if flagConfigFilePath != "" {
		configFilePath = flagConfigFilePath
	} else {
		configFilePath = os.Getenv("GOFILE")
	}

	if configFilePath == "" {
		fmt.Fprintf(os.Stderr, "Cannot find the go file where the struct is located.")
		os.Exit(1)
	}

	if flagConfigStructName == "" {
		fmt.Fprint(os.Stderr, "Please specify a non-empty config struct.")
		os.Exit(1)
	}

	fmt.Printf("Using the struct %s from file %s\n", flagConfigStructName, configFilePath)

	err := gncfg.GenerateConfigLoader(flagProjectName, flagConfigStructName, configFilePath, outputConfigLoaderPath, defaultOutputDotenv, flagOutputDotenvFile, flagDebugLogs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not generate config: %s", err.Error())
		os.Exit(1)
	}
}
