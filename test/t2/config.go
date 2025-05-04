//go:build testcases
// +build testcases

package t2

import "time"

type TestConfigCopy struct {
	AppName string
	Debug   bool
	Timeout time.Duration
	Retries int
}
