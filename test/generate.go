package main

import "testmain/github/genconfig"

func main() {
	genconfig.GenerateConfigLoader("TestConfig1", "config_test.go", "c1_test.go", "", false, "testcases")
	genconfig.GenerateConfigLoader("TestConfigInts", "config_test.go", "c2_test.go", "", false, "testcases")
	genconfig.GenerateConfigLoader("TestConfigUints", "config_test.go", "c3_test.go", "", false, "testcases")
	genconfig.GenerateConfigLoader("TestConfigFloats", "config_test.go", "c4_test.go", "", false, "testcases")
}
