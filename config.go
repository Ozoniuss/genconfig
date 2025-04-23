//go:generate bash -c "go run configgen.go --project=myapp"

package genconfig

import "time"

type Config struct {
	HddSyncPath string
	DryRun, Lol bool
	Timeout     time.Duration
	Port        int
	Port32      uint32
	Port16      int16
}
