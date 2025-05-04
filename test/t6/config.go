//go:build testcases
// +build testcases

package t6

type Nested struct {
	InnerStr  string
	InnerBool bool
}

type TestConfigNested struct {
	AppName string
	Nested  Nested
}
