# Add `default:"..."` struct tag support to genconfig

## Context

Today genconfig treats every field as mandatory: a missing env var produces a `MissingEnvVarsError`. In practice this gets tedious when only a handful of fields ever change per environment. The goal is to let any supported leaf field opt into a default by adding a `default:"..."` struct tag. When the env var is unset, the generated loader falls back to the default; when set, the env var wins.

**Design decision:** the default value is treated as if it were the env var's value. It is *not* validated at generation time, and no typed Go literal is emitted. Instead the generated loader substitutes the raw tag string for the missing env var and runs it through the same parse function. An unparseable default therefore produces an `InvalidEnvVarsError` at runtime, identical to a malformed real env var. This keeps the generator simple and behavior consistent with env var handling.

## Files to modify

- `internal/configgen.go` — tag parsing, TemplateData, template
- `test/testgen/generate.go` — add a `t7` invocation
- `test/config_test.go` — add `t7` cases
- `test/t7/config.go` — new fixture (all supported types with defaults)
- `README.md` — remove the "no defaults" line; add a short example

## Implementation

### 1. `internal/configgen.go`

**Add to `TemplateData` (lines 17–28):**
```go
HasDefault bool
DefaultRaw string // literal string from the `default:"..."` tag; used as fallback value
```

Replace `ErrorVars []string` with two explicit fields to avoid fragile index math once the Missing entry can be absent:
```go
MissingErrVar string // empty iff HasDefault
InvalidErrVar string // empty iff !FormatErr
```

**Tag parsing in `insertTemplateDataEntryForStruct` (around line 186):**
```go
var hasDefault bool
var defaultRaw string
if f.Tag != nil {
    tag := reflect.StructTag(strings.Trim(f.Tag.Value, "`"))
    if raw, ok := tag.Lookup("default"); ok {
        hasDefault = true
        defaultRaw = raw
    }
}
```

Import `reflect`. Set `MissingErrVar` only when `!hasDefault`; set `InvalidErrVar` only when `canHaveFormatErr`.

Reject `default:"..."` on a field whose type resolves to a nested struct — the walker currently recurses into nested structs before inspecting tags and would silently ignore the tag. Return an error from the walker (change the signature to return `error`) and propagate out of `GenerateConfigLoader`.

**Template changes (`goTemplate`, lines 366–489):**

`var (...)` block: iterate both err var fields conditionally.
```
{{- if .MissingErrVar }}
    {{ .MissingErrVar }} = errors.New({{ .EnvVar }}_ENV)
{{- end }}
{{- if .InvalidErrVar }}
    {{ .InvalidErrVar }} = errors.New({{ .EnvVar }}_ENV)
{{- end }}
```

Per-field block: substitute the default when the env var is missing, then fall through into the existing parse block unchanged. Sketch:
```
{{ .AssignmentName }}, ok := os.LookupEnv({{ .EnvVar }}_ENV)
if !ok {
{{- if .HasDefault }}
    {{ .AssignmentName }} = {{ printf "%q" .DefaultRaw }}
    ok = true
{{- else }}
    missingVars = append(missingVars, {{ .MissingErrVar }})
{{- end }}
}
if ok {
    <existing parse-switch block, using .InvalidErrVar in place of (index .ErrorVars 1)>
}
```

This restructuring turns the old `if !ok { ... } else { ... parse ... }` into `if !ok { substitute-or-error }; if ok { parse }`. Generated output for fields WITHOUT defaults should be logically identical to today. A cosmetic diff on regenerated `t1`–`t6` is acceptable; a behavioral diff is not.

**Dotenv output (line 148):** write the default as the value when present, else empty:
```go
fmt.Fprintf(outEnv, "%s=%s\n", field.EnvVar, field.DefaultRaw)
```

### 2. `test/t7/config.go` (new)

A struct with at least one field per supported type carrying a default, plus one required field to confirm mixing works:
```go
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
```

### 3. `test/testgen/generate.go`

Append one more `GenerateConfigLoader(...)` call for `t7/config.go` → `t7/config_gen.go`, prefix `TESTCONFIGDEFAULTS`.

### 4. `test/config_test.go`

Register `t7.LoadTestConfigDefaults` in `loadFuncRegistry` and add test cases:

| Case | Setup | Expected |
|------|-------|----------|
| `t7_all_defaults` | set only `TESTCONFIGDEFAULTS_REQUIRED` | every defaulted field equals its tag default; no error |
| `t7_overrides` | set every env to a different value | every field matches override; no error |
| `t7_required_missing` | no envs set | `IsError = true` (Required is missing) |
| `t7_malformed_env_on_default_field` | `TESTCONFIGDEFAULTS_REQUIRED=x`, `TESTCONFIGDEFAULTS_I=abc` | `IsError = true` via `InvalidEnvVarsError` |

The "malformed default itself" case (an unparseable default string) is best tested manually with a throwaway fixture rather than as an automated test, since a bad default in `t7/config.go` would also poison other cases. Noted but not required for the initial change.

### 5. `README.md`

Remove the line in "Considerations" claiming there are no defaults. Add a short subsection under "Supported parsing functions" with a snippet:
```go
type Config struct {
    Port    int           `default:"8080"`
    Timeout time.Duration `default:"30s"`
    Name    string // no default — still required
}
```
One sentence noting that defaults are parsed at runtime through the same function as env var values and an unparseable default will surface as `InvalidEnvVarsError`.

## Verification

1. `go build ./...` — generator compiles after refactor.
2. `go run test/testgen/generate.go` (or `cd test && go run testgen/generate.go`) — regenerates `t1`–`t7` without error. Diff of `t1`–`t6` generated files should be empty or cosmetic-only; any behavioral diff must be investigated before moving on.
3. `go test --tags=testcases ./test/...` — all existing cases plus the new `t7` cases pass.
4. Eyeball the generated `t7/config_gen.go`: each defaulted field's snippet should read `if !ok { val_X = "<raw>"; ok = true }`, then the usual parse block.
5. Manually run a loader in a scratch program with no env vars set (except `REQUIRED`) to confirm defaults populate.

## Notes and edge cases

- **Tag on a struct-typed field:** explicitly reject with an error from the walker. Silent ignore would confuse users.
- **Empty string default on a non-string type:** will fail `strconv.ParseBool("")` etc. at runtime and surface as `InvalidEnvVarsError` — acceptable per the chosen design.
- **Dotenv quoting:** if a default contains `=` or whitespace the naive `KEY=value` emission may not round-trip through every dotenv parser. Acceptable for now; add a TODO in the code.
- **Template whitespace:** keep aggressive `{{- -}}` trimming; run `gofmt` on the regenerated files.
- **`reflect` import:** added to `internal/configgen.go` imports only, not to generated code.

## Recommended ordering

1. Refactor `TemplateData` (replace `ErrorVars` with `MissingErrVar` / `InvalidErrVar`) and the template accordingly. Regenerate `t1`–`t6` and confirm `go test --tags=testcases ./...` still passes. This isolates the refactor from the feature.
2. Add tag parsing + `HasDefault` / `DefaultRaw` wiring + template branch. Regenerate `t1`–`t6` again — diff should still be empty.
3. Add the `t7` fixture, `testgen` entry, and test cases. Regenerate and run tests.
4. Update `README.md`.
5. Dotenv change (smallest, last; verify with a one-off invocation that outputs a dotenv file).
