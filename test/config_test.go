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
)

type TestConfig1 = t1.TestConfig1
type TestConfigCopy = t2.TestConfigCopy
type TestConfigInts = t3.TestConfigInts
type TestConfigUints = t4.TestConfigUints
type TestConfigFloats = t5.TestConfigFloats
type TestConfigNested = t6.TestConfigNested

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
			StructName: "TestConfig1",
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
			StructName: "TestConfigCopy",
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
			StructName: "TestConfigInts",
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
			StructName: "TestConfigUints",
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
			StructName: "TestConfigFloats",
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
			StructName: "TestConfigNested",
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
	}

	for i := range tcs {
		tcs[i].LoadFuncName = "Load" + tcs[i].StructName
	}

	// Build the registry manually (unfortunately Go can't discover this automatically)
	loadFuncRegistry = map[string]any{
		"LoadTestConfig1":      t1.LoadTestConfig1,
		"LoadTestConfigCopy":   t2.LoadTestConfigCopy,
		"LoadTestConfigInts":   t3.LoadTestConfigInts,
		"LoadTestConfigUints":  t4.LoadTestConfigUints,
		"LoadTestConfigFloats": t5.LoadTestConfigFloats,
		"LoadTestConfigNested": t6.LoadTestConfigNested,
	}

	return tcs
}

type TestCase struct {
	StructName   string
	SetEnvs      func(t *testing.T)
	Expected     any
	LoadFuncName string
	IsError      bool
}

func TestGenerator(t *testing.T) {
	testcases := getTestCases()
	for _, tc := range testcases {
		t.Run(tc.StructName, func(t *testing.T) {
			tc.SetEnvs(t)
			config, err := callLoadFuncByName(t, tc)
			if !tc.IsError && err != nil {
				t.Errorf("unexpected error when parsing config: %s", err)
			}
			if tc.IsError && err != nil {
				t.Errorf("expected errror to occur during parse: %s", err)
			}
			if !reflect.DeepEqual(tc.Expected, config) {
				t.Errorf("expected %+v (%v), got %+v (%v)\n", tc.Expected, reflect.TypeOf(tc.Expected), config, reflect.TypeOf(config))
			}
		})
	}
}
