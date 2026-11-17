//go:build testcases

package t5

type Nested struct {
	InnerStr  string
	InnerBool bool
}

type TestConfigNested struct {
	AppName string
	Nested  Nested
}
