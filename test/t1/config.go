//go:build testcases
// +build testcases

package t1

import "time"

type TestConfig1 struct {
	AppName string
	Debug   bool
	Timeout time.Duration
	Retries int
}
