# `genconfig`

`genconfig` is a command-line tool that generates code for reading your project's configuration from environment variables. It enables you to easily sync all the code required to read your configuration by just updating your struct definition. All you have to do is define a struct representing your project's settings and `genconfig` will do the rest:)

## Example

See [Installation and Usage](#installation-and-usage) for more detailed instructions.

The struct `Config` below defines your project's settings:

<!-- sync:begin file=examples/00-docs/config.go lang=go skip=4 -->
```go
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
<!-- sync:end -->

Based on this, running `genconfig` will generate [config_gen.go](examples/00-docs/config_gen.go) containing code to populate the struct from env vars. In order to do so, call `LoadConfig()` at the start of the program:

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

<!-- sync:begin file=examples/00-docs/.env lang= -->
```
APP_APIKEY=''
APP_LOGLEVEL=''
APP_SERVER_HOST=''
APP_SERVER_PORT=''
APP_SERVER_SHUTDOWNINTERVAL=''
```
<!-- sync:end -->

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
- Every environment variable must be parseable into its corresponding type in the config struct. Fields without a `default` tag must be explicitly set; see [Defaults](#defaults) below. In order to facilitate explicitly setting them, `genconfig` can be configured to output a .env file.

> Note: I'm open to changing those assumptions in the future.

## Supprted parsing functions

Based on the field's type, a different parsing function will be used to convert its value from string. Currently, the following ones are supported:

- `strconv.ParseInt` for all integer types
- `strconv.ParseFloat` for all float types
- `strconv.ParseBool` for bool
- `time.ParseDuration` for `time.Duration`

Support for parsing more types (including custom parsing functions) may be added in the future.

## Defaults

Any supported leaf field can opt into a fallback value by adding a `default:"..."` struct tag. When the environment variable is unset, the generated loader substitutes the raw tag string and runs it through the same parse function as a real env var value. When the variable is set, it takes priority.

```go
type Config struct {
    Port    int           `default:"8080"`
    Timeout time.Duration `default:"30s"`
    Name    string        // no default — still required
}
```

An unparseable default (e.g. `default:"abc"` on an `int` field) will surface as an `InvalidEnvVarsError` at runtime, identical to a malformed env var value.
