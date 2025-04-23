//go:build testcases
// +build testcases

package main

import (
	"testing"
	"time"
)

type TestConfig1 struct {
	AppName string
	Debug   bool
	Timeout time.Duration
	Retries int
}

func TestGeneratedConfig1(t *testing.T) {

	t.Setenv("TEST_APP_NAME", "SuperApp")
	t.Setenv("TEST_DEBUG", "true")
	t.Setenv("TEST_TIMEOUT", "5s")
	t.Setenv("TEST_RETRIES", "3")

	expected := TestConfig1{
		AppName: "SuperApp",
		Debug:   true,
		Timeout: 5 * time.Second,
		Retries: 3,
	}

	actual, err := LoadTestConfig1()
	if err != nil {
		t.Fatalf("expected no error, got: %s", err.Error())
	}

	if expected != actual {
		t.Fatalf("expected configs to be equal: want %+v, got %+v", expected, actual)
	}
}

type TestConfigInts struct {
	Int8Val  int8
	Int16Val int16
	Int32Val int32
	Int64Val int64
}

func TestGeneratedConfigInts(t *testing.T) {
	t.Setenv("TEST_INT8_VAL", "8")
	t.Setenv("TEST_INT16_VAL", "1600")
	t.Setenv("TEST_INT32_VAL", "320000")
	t.Setenv("TEST_INT64_VAL", "6400000000")

	expected := TestConfigInts{
		Int8Val:  8,
		Int16Val: 1600,
		Int32Val: 320000,
		Int64Val: 6400000000,
	}

	actual, err := LoadTestConfigInts()
	if err != nil {
		t.Fatalf("expected no error, got: %s", err)
	}

	if expected != actual {
		t.Fatalf("expected: %+v, got: %+v", expected, actual)
	}
}

type TestConfigMixed struct {
	UInt8Val   uint8
	UInt32Val  uint32
	Float32Val float32
	Float64Val float64
}

// TestConfigUints: unsigned ints
type TestConfigUints struct {
	Uint8Val  uint8
	Uint16Val uint16
	Uint32Val uint32
	Uint64Val uint64
}

func TestGeneratedConfigUints(t *testing.T) {
	t.Setenv("TEST_UINT8_VAL", "8")
	t.Setenv("TEST_UINT16_VAL", "1600")
	t.Setenv("TEST_UINT32_VAL", "320000")
	t.Setenv("TEST_UINT64_VAL", "6400000000")

	expected := TestConfigUints{
		Uint8Val:  8,
		Uint16Val: 1600,
		Uint32Val: 320000,
		Uint64Val: 6400000000,
	}

	actual, err := LoadTestConfigUints()
	if err != nil {
		t.Fatalf("expected no error, got: %s", err)
	}
	if expected != actual {
		t.Fatalf("expected: %+v, got: %+v", expected, actual)
	}
}

type TestConfigFloats struct {
	Float32Val float32
	Float64Val float64
}

func TestGeneratedConfigFloats(t *testing.T) {
	t.Setenv("TEST_FLOAT32_VAL", "1.23")
	t.Setenv("TEST_FLOAT64_VAL", "3.1415")

	expected := TestConfigFloats{
		Float32Val: 1.23,
		Float64Val: 3.1415,
	}

	actual, err := LoadTestConfigFloats()
	if err != nil {
		t.Fatalf("expected no error, got: %s", err)
	}

	// Loose float comparisons
	if delta := actual.Float32Val - expected.Float32Val; delta < -0.001 || delta > 0.001 {
		t.Errorf("float32 mismatch: got %f, want %f", actual.Float32Val, expected.Float32Val)
	}
	if delta := actual.Float64Val - expected.Float64Val; delta < -0.000001 || delta > 0.000001 {
		t.Errorf("float64 mismatch: got %f, want %f", actual.Float64Val, expected.Float64Val)
	}
}
