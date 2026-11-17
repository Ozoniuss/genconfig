# Changelog

## v1.6.0

- `genconfig` no longer panics when the struct passed with `-struct` is not found in the config file; it now reports `struct <name> not found in <file>` and exits with status 1.
- A failed run no longer empties an existing generated loader. The output file is now only written once generation succeeds.
- CLI error messages now end with a newline, so they no longer run into the output of `go generate`.

## v1.5.1

- The legacy `// +build` constraint line is no longer written to generated files; only the `//go:build` line is emitted.
