package main

import (
	"fmt"

	genconfig "github.com/Ozoniuss/genconfig/internal"
)

func main() {
	var err error
	err = genconfig.GenerateConfigLoader("TESTCONFIG1", "TestConfig1", "t1/config.go", "t1/config_gen.go", "", "testcases", false)
	if err != nil {
		fmt.Println("TESTCONFIG1", err)
	}
	err = genconfig.GenerateConfigLoader("TESTCONFIGCOPY", "TestConfigCopy", "t2/config.go", "t2/config_gen.go", "", "testcases", false)
	if err != nil {
		fmt.Println("TESTCONFIGCOPY", err)
	}
	err = genconfig.GenerateConfigLoader("TESTCONFIGINTS", "TestConfigInts", "t3/config.go", "t3/config_gen.go", "", "testcases", false)
	if err != nil {
		fmt.Println("TESTCONFIGINTS", err)
	}
	err = genconfig.GenerateConfigLoader("TESTCONFIGUINTS", "TestConfigUints", "t4/config.go", "t4/config_gen.go", "", "testcases", false)
	if err != nil {
		fmt.Println("TESTCONFIGUINTS", err)
	}
	err = genconfig.GenerateConfigLoader("TESTCONFIGFLOATS", "TestConfigFloats", "t5/config.go", "t5/config_gen.go", "", "testcases", false)
	if err != nil {
		fmt.Println("TESTCONFIGFLOATS", err)
	}
	err = genconfig.GenerateConfigLoader("TESTCONFIGNESTED", "TestConfigNested", "t6/config.go", "t6/config_gen.go", "", "testcases", false)
	if err != nil {
		fmt.Println("TESTCONFIGNESTED", err)
	}
}
