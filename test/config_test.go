//go:build testcases

package test

import (
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

type testCase interface {
	name() string
	run(t *testing.T)
}

type Case[T comparable] struct {
	TestName string
	Load     func() (T, error)
	SetEnvs  func(t *testing.T)
	Expected T
	IsError  bool
}

func (c Case[T]) name() string { return c.TestName }

func (c Case[T]) run(t *testing.T) {
	c.SetEnvs(t)
	config, err := c.Load()
	if !c.IsError && err != nil {
		t.Errorf("unexpected error when parsing config: %s", err)
	}
	if c.IsError && err == nil {
		t.Errorf("expected error to occur during parse, got nil")
	}
	if !c.IsError && config != c.Expected {
		t.Errorf("expected %+v, got %+v (%T)", c.Expected, config, config)
	}
}

func TestGenerator(t *testing.T) {
	for _, tc := range getTestCases() {
		t.Run(tc.name(), tc.run)
	}
}

func getTestCases() []testCase {
	return []testCase{
		Case[t1.TestConfig1]{
			TestName: "TestConfig1",
			Load:     t1.LoadTestConfig1,
			SetEnvs: func(t *testing.T) {
				t.Setenv("TESTCONFIG1_APPNAME", "SuperApp")
				t.Setenv("TESTCONFIG1_DEBUG", "true")
				t.Setenv("TESTCONFIG1_TIMEOUT", "5s")
				t.Setenv("TESTCONFIG1_RETRIES", "3")
			},
			Expected: t1.TestConfig1{
				AppName: "SuperApp",
				Debug:   true,
				Timeout: 5 * time.Second,
				Retries: 3,
			},
		},
		Case[t1.TestConfig1]{
			TestName: "TestConfig1_malformed_retries",
			Load:     t1.LoadTestConfig1,
			SetEnvs: func(t *testing.T) {
				t.Setenv("TESTCONFIG1_APPNAME", "SuperApp")
				t.Setenv("TESTCONFIG1_DEBUG", "true")
				t.Setenv("TESTCONFIG1_TIMEOUT", "5s")
				t.Setenv("TESTCONFIG1_RETRIES", "notanumber")
			},
			IsError: true,
		},
		Case[t2.TestConfigCopy]{
			TestName: "TestConfigCopy",
			Load:     t2.LoadTestConfigCopy,
			SetEnvs: func(t *testing.T) {
				t.Setenv("TESTCONFIGCOPY_APPNAME", "SuperApp")
				t.Setenv("TESTCONFIGCOPY_DEBUG", "true")
				t.Setenv("TESTCONFIGCOPY_TIMEOUT", "5s")
				t.Setenv("TESTCONFIGCOPY_RETRIES", "3")
			},
			Expected: t2.TestConfigCopy{
				AppName: "SuperApp",
				Debug:   true,
				Timeout: 5 * time.Second,
				Retries: 3,
			},
		},
		Case[t3.TestConfigInts]{
			TestName: "TestConfigInts",
			Load:     t3.LoadTestConfigInts,
			SetEnvs: func(t *testing.T) {
				t.Setenv("TESTCONFIGINTS_INT8VAL", "8")
				t.Setenv("TESTCONFIGINTS_INT16VAL", "1600")
				t.Setenv("TESTCONFIGINTS_INT32VAL", "320000")
				t.Setenv("TESTCONFIGINTS_INT64VAL", "6400000000")
			},
			Expected: t3.TestConfigInts{
				Int8Val:  8,
				Int16Val: 1600,
				Int32Val: 320000,
				Int64Val: 6400000000,
			},
		},
		Case[t4.TestConfigUints]{
			TestName: "TestConfigUints",
			Load:     t4.LoadTestConfigUints,
			SetEnvs: func(t *testing.T) {
				t.Setenv("TESTCONFIGUINTS_UINT8VAL", "8")
				t.Setenv("TESTCONFIGUINTS_UINT16VAL", "1600")
				t.Setenv("TESTCONFIGUINTS_UINT32VAL", "320000")
				t.Setenv("TESTCONFIGUINTS_UINT64VAL", "6400000000")
			},
			Expected: t4.TestConfigUints{
				Uint8Val:  8,
				Uint16Val: 1600,
				Uint32Val: 320000,
				Uint64Val: 6400000000,
			},
		},
		Case[t5.TestConfigFloats]{
			TestName: "TestConfigFloats",
			Load:     t5.LoadTestConfigFloats,
			SetEnvs: func(t *testing.T) {
				t.Setenv("TESTCONFIGFLOATS_FLOAT32VAL", "1.23")
				t.Setenv("TESTCONFIGFLOATS_FLOAT64VAL", "3.1415")
			},
			Expected: t5.TestConfigFloats{
				Float32Val: 1.23,
				Float64Val: 3.1415,
			},
		},
		Case[t6.TestConfigNested]{
			TestName: "TestConfigNested",
			Load:     t6.LoadTestConfigNested,
			SetEnvs: func(t *testing.T) {
				t.Setenv("TESTCONFIGNESTED_APPNAME", "NestedApp")
				t.Setenv("TESTCONFIGNESTED_NESTED_INNERSTR", "hello")
				t.Setenv("TESTCONFIGNESTED_NESTED_INNERBOOL", "true")
			},
			Expected: t6.TestConfigNested{
				AppName: "NestedApp",
				Nested: t6.Nested{
					InnerStr:  "hello",
					InnerBool: true,
				},
			},
		},
		Case[t7.TestConfigDefaults]{
			TestName: "t7_all_defaults",
			Load:     t7.LoadTestConfigDefaults,
			SetEnvs: func(t *testing.T) {
				t.Setenv("TESTCONFIGDEFAULTS_REQUIRED", "mustbeset")
			},
			Expected: t7.TestConfigDefaults{
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
		Case[t7.TestConfigDefaults]{
			TestName: "t7_overrides",
			Load:     t7.LoadTestConfigDefaults,
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
			Expected: t7.TestConfigDefaults{
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
		Case[t7.TestConfigDefaults]{
			TestName: "t7_required_missing",
			Load:     t7.LoadTestConfigDefaults,
			SetEnvs:  func(t *testing.T) {},
			IsError:  true,
		},
		Case[t7.TestConfigDefaults]{
			TestName: "t7_malformed_env_on_default_field",
			Load:     t7.LoadTestConfigDefaults,
			SetEnvs: func(t *testing.T) {
				t.Setenv("TESTCONFIGDEFAULTS_REQUIRED", "x")
				t.Setenv("TESTCONFIGDEFAULTS_I", "abc")
			},
			IsError: true,
		},
	}
}
