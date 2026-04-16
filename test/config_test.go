//go:build testcases
// +build testcases

package test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/Ozoniuss/genconfig/test/t1"
	"github.com/Ozoniuss/genconfig/test/t2"
	"github.com/Ozoniuss/genconfig/test/t3"
	"github.com/Ozoniuss/genconfig/test/t4"
	"github.com/Ozoniuss/genconfig/test/t5"
	"github.com/Ozoniuss/genconfig/test/t6"
	"github.com/Ozoniuss/genconfig/test/t7"
)

type TestConfig1 = t1.TestConfig1
type TestConfigCopy = t2.TestConfigCopy
type TestConfigInts = t3.TestConfigInts
type TestConfigUints = t4.TestConfigUints
type TestConfigFloats = t5.TestConfigFloats
type TestConfigNested = t6.TestConfigNested
type TestConfigDefaults = t7.TestConfigDefaults

func callLoadFuncByName(t *testing.T, tc TestCase) (interface{}, error) {
	t.Helper()
	loadFunc := reflect.ValueOf(loadFuncRegistry[tc.LoadFuncName])
	if !loadFunc.IsValid() {
		t.Fatalf("load function not valid: %s", loadFunc)
	}

	results := loadFunc.Call(nil)
	actual := results[0].Interface()
	fmt.Println("actual", actual)
	var err error
	if results[1].Interface() != nil {
		err = results[1].Interface().(error)
	}
	return actual, err
}

var loadFuncRegistry map[string]any

func getTestCases() []TestCase {
	tcs := []TestCase{
		{
			TestName:     "TestConfig1",
			LoadFuncName: "LoadTestConfig1",
			SetEnvs: func(t *testing.T) {
				t.Setenv("TESTCONFIG1_APPNAME", "SuperApp")
				t.Setenv("TESTCONFIG1_DEBUG", "true")
				t.Setenv("TESTCONFIG1_TIMEOUT", "5s")
				t.Setenv("TESTCONFIG1_RETRIES", "3")
			},
			Expected: TestConfig1{
				AppName: "SuperApp",
				Debug:   true,
				Timeout: 5 * time.Second,
				Retries: 3,
			},
		},
		{
			TestName:     "TestConfig1_malformed_retries",
			LoadFuncName: "LoadTestConfig1",
			SetEnvs: func(t *testing.T) {
				t.Setenv("TESTCONFIG1_APPNAME", "SuperApp")
				t.Setenv("TESTCONFIG1_DEBUG", "true")
				t.Setenv("TESTCONFIG1_TIMEOUT", "5s")
				t.Setenv("TESTCONFIG1_RETRIES", "notanumber")
			},
			IsError: true,
		},
		{
			TestName:     "TestConfigCopy",
			LoadFuncName: "LoadTestConfigCopy",
			SetEnvs: func(t *testing.T) {
				t.Setenv("TESTCONFIGCOPY_APPNAME", "SuperApp")
				t.Setenv("TESTCONFIGCOPY_DEBUG", "true")
				t.Setenv("TESTCONFIGCOPY_TIMEOUT", "5s")
				t.Setenv("TESTCONFIGCOPY_RETRIES", "3")
			},
			Expected: TestConfigCopy{
				AppName: "SuperApp",
				Debug:   true,
				Timeout: 5 * time.Second,
				Retries: 3,
			},
		},
		{
			TestName:     "TestConfigInts",
			LoadFuncName: "LoadTestConfigInts",
			SetEnvs: func(t *testing.T) {
				t.Setenv("TESTCONFIGINTS_INT8VAL", "8")
				t.Setenv("TESTCONFIGINTS_INT16VAL", "1600")
				t.Setenv("TESTCONFIGINTS_INT32VAL", "320000")
				t.Setenv("TESTCONFIGINTS_INT64VAL", "6400000000")
			},
			Expected: TestConfigInts{
				Int8Val:  8,
				Int16Val: 1600,
				Int32Val: 320000,
				Int64Val: 6400000000,
			},
		},
		{
			TestName:     "TestConfigUints",
			LoadFuncName: "LoadTestConfigUints",
			SetEnvs: func(t *testing.T) {
				t.Setenv("TESTCONFIGUINTS_UINT8VAL", "8")
				t.Setenv("TESTCONFIGUINTS_UINT16VAL", "1600")
				t.Setenv("TESTCONFIGUINTS_UINT32VAL", "320000")
				t.Setenv("TESTCONFIGUINTS_UINT64VAL", "6400000000")
			},
			Expected: TestConfigUints{
				Uint8Val:  8,
				Uint16Val: 1600,
				Uint32Val: 320000,
				Uint64Val: 6400000000,
			},
		},
		{
			TestName:     "TestConfigFloats",
			LoadFuncName: "LoadTestConfigFloats",
			SetEnvs: func(t *testing.T) {
				t.Setenv("TESTCONFIGFLOATS_FLOAT32VAL", "1.23")
				t.Setenv("TESTCONFIGFLOATS_FLOAT64VAL", "3.1415")
			},
			Expected: TestConfigFloats{
				Float32Val: 1.23,
				Float64Val: 3.1415,
			},
		},
		{
			TestName:     "TestConfigNested",
			LoadFuncName: "LoadTestConfigNested",
			SetEnvs: func(t *testing.T) {
				t.Setenv("TESTCONFIGNESTED_APPNAME", "NestedApp")
				t.Setenv("TESTCONFIGNESTED_NESTED_INNERSTR", "hello")
				t.Setenv("TESTCONFIGNESTED_NESTED_INNERBOOL", "true")
			},
			Expected: TestConfigNested{
				AppName: "NestedApp",
				Nested: t6.Nested{
					InnerStr:  "hello",
					InnerBool: true,
				},
			},
		},
		{
			TestName:     "t7_all_defaults",
			LoadFuncName: "LoadTestConfigDefaults",
			SetEnvs: func(t *testing.T) {
				t.Setenv("TESTCONFIGDEFAULTS_REQUIRED", "mustbeset")
			},
			Expected: TestConfigDefaults{
				Required: "mustbeset",
				Str:      "hello",
				B:        true,
				I:        -5,
				I8:       -1,
				I16:      16,
				I32:      32,
				I64:      64,
				U:        5,
				U8:       1,
				U16:      16,
				U32:      32,
				U64:      64,
				F32:      1.5,
				F64:      2.5,
				D:        1500 * time.Millisecond,
			},
		},
		{
			TestName:     "t7_overrides",
			LoadFuncName: "LoadTestConfigDefaults",
			SetEnvs: func(t *testing.T) {
				t.Setenv("TESTCONFIGDEFAULTS_REQUIRED", "overridden")
				t.Setenv("TESTCONFIGDEFAULTS_STR", "world")
				t.Setenv("TESTCONFIGDEFAULTS_B", "false")
				t.Setenv("TESTCONFIGDEFAULTS_I", "99")
				t.Setenv("TESTCONFIGDEFAULTS_I8", "7")
				t.Setenv("TESTCONFIGDEFAULTS_I16", "100")
				t.Setenv("TESTCONFIGDEFAULTS_I32", "200")
				t.Setenv("TESTCONFIGDEFAULTS_I64", "300")
				t.Setenv("TESTCONFIGDEFAULTS_U", "10")
				t.Setenv("TESTCONFIGDEFAULTS_U8", "2")
				t.Setenv("TESTCONFIGDEFAULTS_U16", "20")
				t.Setenv("TESTCONFIGDEFAULTS_U32", "200")
				t.Setenv("TESTCONFIGDEFAULTS_U64", "2000")
				t.Setenv("TESTCONFIGDEFAULTS_F32", "9.9")
				t.Setenv("TESTCONFIGDEFAULTS_F64", "8.8")
				t.Setenv("TESTCONFIGDEFAULTS_D", "3s")
			},
			Expected: TestConfigDefaults{
				Required: "overridden",
				Str:      "world",
				B:        false,
				I:        99,
				I8:       7,
				I16:      100,
				I32:      200,
				I64:      300,
				U:        10,
				U8:       2,
				U16:      20,
				U32:      200,
				U64:      2000,
				F32:      9.9,
				F64:      8.8,
				D:        3 * time.Second,
			},
		},
		{
			TestName:     "t7_required_missing",
			LoadFuncName: "LoadTestConfigDefaults",
			SetEnvs:      func(t *testing.T) {},
			IsError:      true,
		},
		{
			TestName:     "t7_malformed_env_on_default_field",
			LoadFuncName: "LoadTestConfigDefaults",
			SetEnvs: func(t *testing.T) {
				t.Setenv("TESTCONFIGDEFAULTS_REQUIRED", "x")
				t.Setenv("TESTCONFIGDEFAULTS_I", "abc")
			},
			IsError: true,
		},
	}

	// Build the registry manually (unfortunately Go can't discover this automatically)
	loadFuncRegistry = map[string]any{
		"LoadTestConfig1":        t1.LoadTestConfig1,
		"LoadTestConfigCopy":     t2.LoadTestConfigCopy,
		"LoadTestConfigInts":     t3.LoadTestConfigInts,
		"LoadTestConfigUints":    t4.LoadTestConfigUints,
		"LoadTestConfigFloats":   t5.LoadTestConfigFloats,
		"LoadTestConfigNested":   t6.LoadTestConfigNested,
		"LoadTestConfigDefaults": t7.LoadTestConfigDefaults,
	}

	return tcs
}

type TestCase struct {
	TestName     string
	LoadFuncName string
	SetEnvs      func(t *testing.T)
	Expected     any
	IsError      bool
}

func TestGenerator(t *testing.T) {
	testcases := getTestCases()
	for _, tc := range testcases {
		t.Run(tc.TestName, func(t *testing.T) {
			tc.SetEnvs(t)
			config, err := callLoadFuncByName(t, tc)
			if !tc.IsError && err != nil {
				t.Errorf("unexpected error when parsing config: %s", err)
			}
			if tc.IsError && err == nil {
				t.Errorf("expected error to occur during parse, got nil")
			}
			if !tc.IsError && !reflect.DeepEqual(tc.Expected, config) {
				t.Errorf("expected %+v (%v), got %+v (%v)\n", tc.Expected, reflect.TypeOf(tc.Expected), config, reflect.TypeOf(config))
			}
		})
	}
}
