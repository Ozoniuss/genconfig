//go:build testcases
// +build testcases

package t7

import "time"

type TestConfigDefaults struct {
	Required string        // no default — must still be set
	Str      string        `default:"hello"`
	B        bool          `default:"true"`
	I        int           `default:"-5"`
	I8       int8          `default:"-1"`
	I16      int16         `default:"16"`
	I32      int32         `default:"32"`
	I64      int64         `default:"64"`
	U        uint          `default:"5"`
	U8       uint8         `default:"1"`
	U16      uint16        `default:"16"`
	U32      uint32        `default:"32"`
	U64      uint64        `default:"64"`
	F32      float32       `default:"1.5"`
	F64      float64       `default:"2.5"`
	D        time.Duration `default:"1500ms"`
}
