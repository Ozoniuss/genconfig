# `genconfig`

`genconfig` is a command-line tool that generates code for reading your project's configuration from environment variables. It enables you to easily sync all the code required to read your configuration by just updating your struct definition. All you have to do is define a struct representing your project's settings and `genconfig` will do the rest:)

## Example

See [Installation and Usage](#installation-and-usage) for more detailed instructions.

The struct `Config` below defines your project's settings:

```go
//go:generate go tool genconfig -struct=Config -project=App
type Config struct {
	Apikey   string
	Loglevel string
	Server   ServerConfig
}

type ServerConfig struct {
	Host             string
	Port             int
	ShutdownInterval time.Duration
}
```

Based on this, running `genconfig` will generate the following file:

```go
package main

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	APP_APIKEY_ENV                  = "APP_APIKEY"
	APP_LOGLEVEL_ENV                = "APP_LOGLEVEL"
	APP_SERVER_HOST_ENV             = "APP_SERVER_HOST"
	APP_SERVER_PORT_ENV             = "APP_SERVER_PORT"
	APP_SERVER_SHUTDOWNINTERVAL_ENV = "APP_SERVER_SHUTDOWNINTERVAL"
)

var (
	ErrAppApikeyEnvMissing                 = errors.New(APP_APIKEY_ENV)
	ErrAppLoglevelEnvMissing               = errors.New(APP_LOGLEVEL_ENV)
	ErrAppServerHostEnvMissing             = errors.New(APP_SERVER_HOST_ENV)
	ErrAppServerPortEnvMissing             = errors.New(APP_SERVER_PORT_ENV)
	ErrAppServerPortEnvInvalid             = errors.New(APP_SERVER_PORT_ENV)
	ErrAppServerShutdownintervalEnvMissing = errors.New(APP_SERVER_SHUTDOWNINTERVAL_ENV)
	ErrAppServerShutdownintervalEnvInvalid = errors.New(APP_SERVER_SHUTDOWNINTERVAL_ENV)
)

func LoadConfig() (Config, error) {
	var config Config
	var missingVars []error
	var formatVars []error
	val_Apikey, ok := os.LookupEnv(APP_APIKEY_ENV)
	if !ok {
		missingVars = append(missingVars, ErrAppApikeyEnvMissing)
	} else {
		config.Apikey = val_Apikey
	}
	val_Loglevel, ok := os.LookupEnv(APP_LOGLEVEL_ENV)
	if !ok {
		missingVars = append(missingVars, ErrAppLoglevelEnvMissing)
	} else {
		config.Loglevel = val_Loglevel
	}
	val_Server_Host, ok := os.LookupEnv(APP_SERVER_HOST_ENV)
	if !ok {
		missingVars = append(missingVars, ErrAppServerHostEnvMissing)
	} else {
		config.Server.Host = val_Server_Host
	}
	val_Server_Port, ok := os.LookupEnv(APP_SERVER_PORT_ENV)
	if !ok {
		missingVars = append(missingVars, ErrAppServerPortEnvMissing)
	} else {
		parsed, err := strconv.Atoi(val_Server_Port)
		if err != nil {
			formatVars = append(formatVars, ErrAppServerPortEnvInvalid)
		} else {
			config.Server.Port = parsed
		}
	}
	val_Server_ShutdownInterval, ok := os.LookupEnv(APP_SERVER_SHUTDOWNINTERVAL_ENV)
	if !ok {
		missingVars = append(missingVars, ErrAppServerShutdownintervalEnvMissing)
	} else {
		parsed, err := time.ParseDuration(val_Server_ShutdownInterval)
		if err != nil {
			formatVars = append(formatVars, ErrAppServerShutdownintervalEnvInvalid)
		} else {
			config.Server.ShutdownInterval = parsed
		}
	}

	if len(missingVars) > 0 || len(formatVars) > 0 {
		var verr error
		if len(missingVars) > 0 {
			verr = errors.Join(verr, MissingEnvVarsError{vars: missingVars})
		}
		if len(formatVars) > 0 {
			verr = errors.Join(verr, InvalidEnvVarsError{vars: missingVars})
		}
		return Config{}, verr
	}

	return config, nil
}

type MissingEnvVarsError struct {
	vars []error
}

func (m MissingEnvVarsError) Unwrap() []error {
	return m.vars
}

func (m MissingEnvVarsError) Error() string {
	if len(m.vars) == 0 {
		return ""
	}
	varsstr := make([]string, 0, len(m.vars))
	for _, v := range m.vars {
		varsstr = append(varsstr, v.Error())
	}
	return "envs " + strings.Join(varsstr, ",") + " are not set"
}

type InvalidEnvVarsError struct {
	vars []error
}

func (m InvalidEnvVarsError) Unwrap() []error {
	return m.vars
}

func (m InvalidEnvVarsError) Error() string {
	if len(m.vars) == 0 {
		return ""
	}
	varsstr := make([]string, 0, len(m.vars))
	for _, v := range m.vars {
		varsstr = append(varsstr, v.Error())
	}
	return "envs " + strings.Join(varsstr, ",") + " have an invalid value"
}
```

You can then use the exported `LoadConfig()` function to populate your struct:

```go
package main

func main() {
	c, err := LoadConfig()
	if err != nil {
		fmt.Printf("failed to read config: %s", err.Error())
		os.Exit(1)
	}
	
	fmt.Printf("%+v", c)
}
```

Optionally, it can also generate a `.env` file containing all the environment variables that it reads (without setting a value for them):

```
APP_APIKEY=
APP_LOGLEVEL=
APP_SERVER_HOST=
APP_SERVER_PORT=
APP_SERVER_SHUTDOWNINTERVAL=
```

## Installation and Usage

There are several ways you can use `genconfig`. The flag `-struct` is used to denote your config struct, and optionally `-project` denotes your project name, if you want a prefix for the environment variables. Note that the struct does not necessarily need to be named `Config`.

When using the `//go:generate` directive, you can generate your loader by calling `go generate` from the root of your module.

All methods illustrated below have a working example in the `examples` directory. Based on your preferred approach, running the generator may vary. Please refer to the instructions for each approach for more details.

### go tool (Go 1.24+)

```
go get -tool github.com/Ozoniuss/genconfig
```

Then, next to your struct definition, add

```
//go:generate go tool genconfig -struct=Config -project=Myapp
```

### tools.go 

See how to use the tools.go pattern [here](https://www.jvt.me/posts/2022/06/15/go-tools-dependency-management/). This pattern was popular before the `go tool` directive.

Add the following to your tools.go file:

```go
//go:build tools
// +build tools

package main

import (
	_ "github.com/Ozoniuss/genconfig"
)
```

Next to your struct definition, add

```
//go:generate go run github.com/Ozoniuss/genconfig -struct=Config -project=Myapp
```

### binary installation

You can also install the `genconfig` binary and either call the downloaded executable in the `//go:generate` directive, or call the executable directly. For example, if you do 

```
go install github.com/Ozoniuss/genconfig
```

then you can use `genconfig` directly in your tags

```
//go:generate go run genconfig -struct=Config -project=Myapp
```

or call it directly from the command line:

```
genconfig -struct=Config -path=config.go -project=Myapp
```

> ⚠️ You need to provide the `-path` flag if you use the executable directly, otherwise `genconfig` will not be able to locate your config.

## Considerations

`genconfig` is very opinionated in its approach. It specifically assumes that:

- There must be only one exported config struct definition for the whole project.
- The generated config loader will be created in the same package as the config struct definition. This is to allow one import to reference both the `LoadConfig()` function and the config struct definition, as well as to make it easier for the `LoadConfig()` function to return that struct.
- For each struct field, an associated environment variable name will automatically be created and follows [this rule](https://github.com/Ozoniuss/genconfig/blob/283a5252de20a4fa9693499412b861b348ea1a75/internal/configgen.go#L196). The name is not configurable.
- Every environment variable must be set explicity (no empty values are allowed), thus there are no defaults. Also, they must be parseable into their corresponding type in the config struct. In order to facilitate explicitly setting them, `genconfig` can be configured to output a .env file where they're all set to the empty value.

> Note: I'm open to changing those assumptions in the future.

## Supprted parsing functions

Based on the field's type, a different parsing function will be used to convert its value from string. Currently, the following ones are supported:

- `strconv.ParseInt` for all integer types
- `strconv.ParseFloat` for all float types
- `strconv.ParseBool` for bool
- `time.ParseDuration` for `time.Duration`

Support for parsing more types (including custom parsing functions) may be added in the future.
