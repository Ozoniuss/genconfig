package main

import genconfig "testmain/github/genconfig/internal"

func main() {
	genconfig.GenerateConfigLoader("TESTCONFIG1", "TestConfig1", "config_test.go", "c1_test.go", "", false, "testcases", false)
	genconfig.GenerateConfigLoader("TESTCONFIGCOPY", "TestConfigCopy", "config_test.go", "c1_1_test.go", "", false, "testcases", false)
	genconfig.GenerateConfigLoader("TESTCONFIGINTS", "TestConfigInts", "config_test.go", "c2_test.go", "", false, "testcases", false)
	genconfig.GenerateConfigLoader("TESTCONFIGUINTS", "TestConfigUints", "config_test.go", "c3_test.go", "", false, "testcases", false)
	genconfig.GenerateConfigLoader("TESTCONFIGFLOATS", "TestConfigFloats", "config_test.go", "c4_test.go", "", false, "testcases", false)
	genconfig.GenerateConfigLoader("TESTCONFIGNESTED", "TestConfigNested", "config_test.go", "c5_test.go", "", false, "testcases", false)
}
