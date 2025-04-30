# `genconfig`

`genconfig` is a command-line tool that allows you to generate code for reading your project's configuration through env variables. Based on your project's  config struct definition, it creates:

- A set of environment variables corresponding to each field of your config struct (see [rules](#how-are-environment-variable-names-generated));
- A set of exported constants corresponding to each environment variable;
- A `LoadConfig()` function that is able to read those values from environment variables, and return errors if those values are not set (empty) or not parsable to the config's type.

## Installation

### go tool (Go 1.24+)

```
go get -tool github.com/Ozoniuss/genconfig@latest
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

### binary installation


```
go install github.com/Ozoniuss/genconfig@latest
```

## Usage

## Supprted parsing functions

## How are environment variable names generated?