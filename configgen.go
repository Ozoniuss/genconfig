package genconfig

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"
)

const (
	defaultInputFile          = "config.go"
	defaultOutputDotenv       = ".env"
	defaultOutputConfigLoader = "config_gen.go"
)

type TemplateData struct {
	Name      string
	EnvVar    string
	Type      string
	ParseFunc string
	ErrorVar  string
	FormatErr bool
	BitSize   int    // used to determine how to call parseFunc
	CastFunc  string // parseInt and parseUint return 64bit numbers, need to cast
}

func GenerateConfigLoader(projectPrefix, configStructName, inputFile, outputLoader, outputDotenv string, generateEnv bool, testBuildTag string) error {

	prefix, err := getProjectNamePrefix(projectPrefix)
	if err != nil {
		panic(err)
	}
	fmt.Printf("using project name prefix %s\n", prefix)

	outputImports := setupImportsAlwaysNeeded()

	// Parse config.go
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, inputFile, nil, parser.AllErrors)
	if err != nil {
		panic(err)
	}

	var fields []TemplateData

	for _, decl := range node.Decls {

		// ignore function and bad declarations
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.TYPE {
			// only look for type (struct) declarations
			continue
		}

		for _, spec := range gen.Specs {
			ts, ok := spec.(*ast.TypeSpec)

			if !ok || ts.Name.Name != configStructName {
				continue
			}

			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				continue
			}

			if st.Fields != nil {
				for _, f := range st.Fields.List {
					if len(f.Names) == 0 {
						continue
					}
					fmt.Println("names", f.Names)

					// in the same struct, you can have multiple fields of the
					// same type declared on the same line
					for _, n := range f.Names {
						name := n.Name
						typ := convertTypeIdentifierToString(f.Type)
						envKey := prefix + "_" + strings.ToUpper(snakeCase(name))

						parseFunc, formatErr, bitSize, castFunc, ok := lookupParseFunc(typ)
						if !ok {
							panic("unsupported type in config: " + typ)
						}

						if p := pkgForParseFunc(parseFunc); p != "" {
							outputImports[p] = struct{}{}
						}

						fields = append(fields, TemplateData{
							Name:      name,
							EnvVar:    envKey,
							Type:      typ,
							ParseFunc: parseFunc,
							ErrorVar:  "Err" + name + "NotSet",
							FormatErr: formatErr,
							BitSize:   bitSize,
							CastFunc:  castFunc,
						})

					}

				}
			}

		}
	}

	importList := generateImportsListAsTemplateString(outputImports)

	// Generate config_gen.go
	outGo, _ := os.Create(outputLoader)
	defer outGo.Close()

	goTemplate.Execute(outGo, struct {
		Prefix       string
		StructName   string
		Fields       []TemplateData
		TestBuildTag string
		ImportList   string
	}{
		Prefix:       prefix,
		StructName:   configStructName,
		Fields:       fields,
		TestBuildTag: testBuildTag,
		ImportList:   importList,
	})

	// Generate .env
	if generateEnv {
		outEnv, _ := os.Create(outputDotenv)
		defer outEnv.Close()
		for _, field := range fields {
			fmt.Fprintf(outEnv, "%s=\n", field.EnvVar)
		}
	}

	return nil
}

func main() {
	err := GenerateConfigLoader("", "Config", defaultInputFile, defaultOutputConfigLoader, defaultOutputDotenv, true, "")
	if err != nil {
		fmt.Printf("failed to generate config: %v", err.Error())
		os.Exit(1)
	}
}

func setupImportsAlwaysNeeded() map[string]struct{} {
	return map[string]struct{}{
		`"fmt"`:     {},
		`"os"`:      {},
		`"strings"`: {},
	}

}

func generateImportsListAsTemplateString(outputImports map[string]struct{}) string {
	pkgs := make([]string, 0, len(outputImports))
	for p := range outputImports {
		pkgs = append(pkgs, p)
	}
	slices.Sort(pkgs) // Go 1.21+
	importList := strings.Join(pkgs, "\n    ")
	return importList
}

func getProjectNamePrefix(suppliedPrefix string) (string, error) {
	if suppliedPrefix != "" {
		return suppliedPrefix, nil
	}
	var prefix string
	if len(os.Args) > 1 {
		prefix = os.Args[1]
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("cannot get working directory: %w", err)
		}
		base := filepath.Base(cwd)
		prefix = base
	}
	prefix = strings.ToUpper(prefix)
	return prefix, nil
}

func convertTypeIdentifierToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return convertTypeIdentifierToString(t.X) + "." + t.Sel.Name
	default:
		panic("expected identifier or selector expression")
	}
}

func snakeCase(in string) string {
	var out []rune
	for i, r := range in {
		if i > 0 && r >= 'A' && r <= 'Z' {
			out = append(out, '_')
		}
		out = append(out, r)
	}
	return strings.ToLower(string(out))
}

func lookupParseFunc(typ string) (parseFunc string, canHaveFormatErr bool, bitSize int, castFuncForIntAndFloat string, ok bool) {
	switch typ {
	case "string":
		return "raw", false, 0, "", true
	case "bool":
		return "strconv.ParseBool", true, 0, "", true

	case "int":
		return "strconv.Atoi", true, 0, "", true
	case "int64":
		return "strconv.ParseInt", true, 64, "int64", true
	case "int32":
		return "strconv.ParseInt", true, 32, "int32", true
	case "int16":
		return "strconv.ParseInt", true, 16, "int16", true
	case "int8":
		return "strconv.ParseInt", true, 8, "int8", true

	case "uint":
		return "strconv.ParseUint", true, 0, "uint", true
	case "uint64":
		return "strconv.ParseUint", true, 64, "uint64", true
	case "uint32":
		return "strconv.ParseUint", true, 32, "uint32", true
	case "uint16":
		return "strconv.ParseUint", true, 16, "uint16", true
	case "uint8":
		return "strconv.ParseUint", true, 8, "uint8", true

	case "float32":
		return "strconv.ParseFloat", true, 32, "float32", true
	case "float64":
		return "strconv.ParseFloat", true, 64, "", true

	case "time.Duration":
		return "time.ParseDuration", true, 0, "", true

	default:
		return "", false, 0, "", false
	}
}

func pkgForParseFunc(fn string) string {
	switch {
	case strings.HasPrefix(fn, "strconv."):
		return `"strconv"`
	case fn == "time.ParseDuration":
		return `"time"`
	default:
		return ""
	}
}

var goTemplate = template.Must(template.New("config").Parse(`// Code generated by configgen.go; DO NOT EDIT.

{{- if .TestBuildTag }}
//go:build {{ .TestBuildTag }}
// +build {{ .TestBuildTag }}

{{- end }}

package main

import (
    {{ .ImportList }}
)

const (
{{- range .Fields }}
    {{ .EnvVar }}_ENV = "{{ .EnvVar }}"
{{- end }}
)

func Load{{ .StructName }}() ({{ .StructName }}, error) {
    var config {{ .StructName }}
    var missingVars []string
    var formatVars []string

{{- range .Fields }}
    val_{{ .Name }}, ok := os.LookupEnv({{ .EnvVar }}_ENV)
    if !ok {
        missingVars = append(missingVars, {{ .EnvVar }}_ENV)
    } else {
        {{- if eq .ParseFunc "raw" }}
        config.{{ .Name }} = val_{{ .Name }}
        {{- else if eq .ParseFunc "strconv.Atoi" }}
        parsed, err := strconv.Atoi(val_{{ .Name }})
        if err != nil {
            formatVars = append(formatVars, {{ .EnvVar }}_ENV)
        } else {
            config.{{ .Name }} = {{ if .CastFunc }}{{ .CastFunc }}(parsed){{ else }}parsed{{ end }}
        }
        {{- else if or (eq .ParseFunc "strconv.ParseInt") (eq .ParseFunc "strconv.ParseUint") }}
        parsed, err := {{ .ParseFunc }}(val_{{ .Name }}, 10, {{ .BitSize }})
        if err != nil {
            formatVars = append(formatVars, {{ .EnvVar }}_ENV)
        } else {
            config.{{ .Name }} = {{ if .CastFunc }}{{ .CastFunc }}(parsed){{ else }}parsed{{ end }}
        }
		{{- else if eq .ParseFunc "strconv.ParseFloat" }}
        parsed, err := strconv.ParseFloat(val_{{ .Name }}, {{ .BitSize }})
        if err != nil {
            formatVars = append(formatVars, {{ .EnvVar }}_ENV)
        } else {
            config.{{ .Name }} = {{ if .CastFunc }}{{ .CastFunc }}(parsed){{ else }}parsed{{ end }}
        }
        {{- else }}
        parsed, err := {{ .ParseFunc }}(val_{{ .Name }})
        if err != nil {
            formatVars = append(formatVars, {{ .EnvVar }}_ENV)
        } else {
            config.{{ .Name }} = parsed
        }
        {{- end }}
    }
{{- end }}

    if len(missingVars) > 0 || len(formatVars) > 0 {
        var parts []string
        if len(missingVars) > 0 {
            parts = append(parts, fmt.Sprintf("missing env vars: %s", strings.Join(missingVars, ", ")))
        }
        if len(formatVars) > 0 {
            parts = append(parts, fmt.Sprintf("invalid format in env vars: %s", strings.Join(formatVars, ", ")))
        }
        return {{ .StructName }}{}, fmt.Errorf(strings.Join(parts, "; "))
    }

    return config, nil
}
`))
