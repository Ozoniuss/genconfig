package main

import "testmain/github/genconfig"

func main() {
	genconfig.GenerateConfigLoader("TESTCONFIG1", "TestConfig1", "config_test.go", "c1_test.go", "", false, "testcases")
	genconfig.GenerateConfigLoader("TESTCONFIGINTS", "TestConfigInts", "config_test.go", "c2_test.go", "", false, "testcases")
	genconfig.GenerateConfigLoader("TESTCONFIGUINTS", "TestConfigUints", "config_test.go", "c3_test.go", "", false, "testcases")
	genconfig.GenerateConfigLoader("TESTCONFIGFLOATS", "TestConfigFloats", "config_test.go", "c4_test.go", "", false, "testcases")
	genconfig.GenerateConfigLoader("TESTCONFIGNESTED", "TestConfigNested", "config_test.go", "c5_test.go", "", false, "testcases")
}
