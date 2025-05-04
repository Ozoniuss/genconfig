package main

import (
	"fmt"

	genconfig "github.com/Ozoniuss/genconfig/internal"
)

func main() {
	var err error
	err = genconfig.GenerateConfigLoader("TESTCONFIG1", "TestConfig1", "config_test.go", "c1/c111_test.go", "", "testcases", false)
	if err != nil {
		fmt.Println("TESTCONFIG1", err)
	}
	err = genconfig.GenerateConfigLoader("TESTCONFIGCOPY", "TestConfigCopy", "config_test.go", "c1_1_test.go", "", "testcases", false)
	if err != nil {
		fmt.Println("TESTCONFIGCOPY", err)
	}
	err = genconfig.GenerateConfigLoader("TESTCONFIGINTS", "TestConfigInts", "config_test.go", "c2_test.go", "", "testcases", false)
	if err != nil {
		fmt.Println("TESTCONFIGINTS", err)
	}
	err = genconfig.GenerateConfigLoader("TESTCONFIGUINTS", "TestConfigUints", "config_test.go", "c3_test.go", "", "testcases", false)
	if err != nil {
		fmt.Println("TESTCONFIGUINTS", err)
	}
	err = genconfig.GenerateConfigLoader("TESTCONFIGFLOATS", "TestConfigFloats", "config_test.go", "c4_test.go", "", "testcases", false)
	if err != nil {
		fmt.Println("TESTCONFIGFLOATS", err)
	}
	err = genconfig.GenerateConfigLoader("TESTCONFIGNESTED", "TestConfigNested", "config_test.go", "c5_test.go", "", "testcases", false)
	if err != nil {
		fmt.Println("TESTCONFIGNESTED", err)
	}
}
