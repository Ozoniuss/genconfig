package genconfig

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"
	"unicode"
)

type TemplateData struct {
	PackageName    string
	Name           string
	AssignmentName string
	EnvVar         string
	ParseFunc      string
	// first is always missing error, second is invalid, if exists
	ErrorVars []string
	FormatErr bool
	BitSize   int    // used to determine how to call parseFunc
	CastFunc  string // parseInt and parseUint return 64bit numbers, need to cast
}

func printformat(debug bool, format string, a ...any) {
	if debug {
		fmt.Printf(format, a...)
	}
}
func printline(debug bool, a ...any) {
	if debug {
		fmt.Println(a...)
	}
}

func createOutputFile(outputGeneratedConfigFile string) (*os.File, error) {
	outFullPath, err := filepath.Abs(outputGeneratedConfigFile)
	if err != nil {
		return nil, fmt.Errorf("cannot get expected absolute path: %w", err)
	}
	baseDir := filepath.Dir(outFullPath)
	fileName := filepath.Base(outFullPath)

	if !strings.HasSuffix(fileName, ".go") {
		outFullPath = outFullPath + ".go"
		fmt.Fprintf(os.Stderr, "output file %s does not have .go suffix; adding it", fileName)
	}

	err = os.MkdirAll(baseDir, 0750)
	if err != nil {
		return nil, fmt.Errorf("could not create directory for output file: %w", err)
	}
	generatedFile, err := os.Create(outFullPath)
	if err != nil {
		return nil, fmt.Errorf("cannot generate output file %s: %w", outFullPath, err)
	}
	return generatedFile, nil
}

func GenerateConfigLoader(projectPrefix, configStructName, inputFile, outputGeneratedConfigFile, outputDotenv string, testBuildTag string, debug bool) error {

	printformat(debug, "using project name prefix %s\n", projectPrefix)

	generatedFile, err := createOutputFile(outputGeneratedConfigFile)
	if err != nil {
		return fmt.Errorf("could not create output config loader file: %w", err)
	}
	defer generatedFile.Close()

	outputImports := setupImportsAlwaysNeeded()

	// Parse config.go
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, inputFile, nil, parser.AllErrors)
	if err != nil {
		return fmt.Errorf("could not parse file with config struct: %w", err)
	}
	if node == nil || node.Name == nil {
		printline(debug, "node name is nil")
		return errors.New("could not parse package name")
	}
	packageName := node.Name.Name
	printformat(debug, "using package name %s", packageName)

	var fields []TemplateData

	printformat(debug, "node %+v", *node)

	allTopLevelStructDefinitions := getAllTopLevelStructDefinitions(node)
	printline(debug, "all struct defintions", allTopLevelStructDefinitions)

	configTypeDefinition := allTopLevelStructDefinitions[configStructName]
	parentNames := []string{}

	insertTemplateDataEntryForStruct(configTypeDefinition, configStructName, &parentNames, projectPrefix, outputImports, &fields, allTopLevelStructDefinitions, debug)

	importList := generateImportsListAsTemplateString(outputImports)

	goTemplate.Execute(generatedFile, struct {
		Prefix       string
		StructName   string
		Fields       []TemplateData
		TestBuildTag string
		ImportList   string
		PackageName  string
	}{
		Prefix:       projectPrefix,
		StructName:   configStructName,
		Fields:       fields,
		TestBuildTag: testBuildTag,
		ImportList:   importList,
		PackageName:  packageName,
	})

	printline(debug, fields)

	// Generate .env
	if outputDotenv != "" {
		outEnv, _ := os.Create(outputDotenv)
		defer outEnv.Close()
		for _, field := range fields {
			fmt.Fprintf(outEnv, "%s=\n", field.EnvVar)
		}
	}

	return nil
}

func insertTemplateDataEntryForStruct(structDefinition *ast.StructType, structName string, parentNames *[]string, projectPrefix string, outputImports map[string]struct{}, templateData *[]TemplateData, allTopLevelStructDefinitions map[string]*ast.StructType, debug bool) {

	if structDefinition.Fields == nil {
		return
	}

	for _, f := range structDefinition.Fields.List {
		if len(f.Names) == 0 {
			continue
		}
		// debug multiple names in a row
		// printline(debug, "names", f.Names)

		// in the same struct, you can have multiple fields of the
		// same type declared on the same line
		for _, n := range f.Names {
			// fullname := n.Name
			fullname := strings.Join(append(*parentNames, n.Name), ".")
			assignmentName := "val_" + strings.Join(append(*parentNames, n.Name), "_")
			typ := convertTypeIdentifierToString(f.Type)
			printline(debug, "identifier type", typ, "field name", n.Name)

			// we have encoutnered a struct defined in the same file
			if childDefinition, ok := allTopLevelStructDefinitions[typ]; ok {
				*parentNames = append(*parentNames, n.Name)
				insertTemplateDataEntryForStruct(childDefinition, typ, parentNames, projectPrefix, outputImports, templateData, allTopLevelStructDefinitions, debug)
				*parentNames = (*parentNames)[:len(*parentNames)-1]
			} else {
				canonicalNameList := append([]string{projectPrefix}, *parentNames...)
				canonicalNameList = append(canonicalNameList, n.Name)
				envKey := getEnvKey(canonicalNameList)

				parseFunc, canHaveFormatErr, bitSize, castFunc, ok := lookupParseFunc(typ)
				if !ok {
					panic("unsupported type in config: " + typ)
				}
				errorsVars := []string{getErrKey(canonicalNameList) + "Missing"}
				if canHaveFormatErr {
					errorsVars = append(errorsVars, getErrKey(canonicalNameList)+"Invalid")
				}

				if p := pkgForParseFunc(parseFunc); p != "" {
					outputImports[p] = struct{}{}
				}

				*templateData = append(*templateData, TemplateData{
					Name:           fullname,
					AssignmentName: assignmentName,
					EnvVar:         envKey,
					ParseFunc:      parseFunc,
					ErrorVars:      errorsVars,
					FormatErr:      canHaveFormatErr,
					BitSize:        bitSize,
					CastFunc:       castFunc,
				})
			}

		}

	}
}

func getEnvKey(canonicalNameList []string) string {
	sb := &strings.Builder{}
	for _, part := range canonicalNameList {
		for _, r := range part {
			// keep only letters and digits in the env var name. This is prone
			// to errors e.g. if someone names a field "my_field" and has a field
			// called "my" of type struct that has a field called "field" but
			// come on
			if unicode.IsDigit(r) {
				sb.WriteRune(r)
			}
			if unicode.IsLetter(r) {
				sb.WriteRune(unicode.ToUpper(r))
			}
		}
		sb.WriteRune('_')
	}
	return strings.TrimSuffix(sb.String(), "_")
}

func getErrKey(canonicalNameList []string) string {
	sb := &strings.Builder{}
	sb.WriteString("Err")
	for _, part := range canonicalNameList {
		for i, r := range part {
			// keep only letters and digits in the env var name. This is prone
			// to errors e.g. if someone names a field "my_field" and has a field
			// called "my" of type struct that has a field called "field" but
			// come on
			if unicode.IsDigit(r) {
				sb.WriteRune(r)
			}
			if unicode.IsLetter(r) && i == 0 {
				sb.WriteRune(unicode.ToUpper(r))
			} else if unicode.IsLetter(r) && i != 0 {
				sb.WriteRune(unicode.ToLower(r))
			}
		}
	}
	sb.WriteString("Env")
	return sb.String()
}

func setupImportsAlwaysNeeded() map[string]struct{} {
	return map[string]struct{}{
		`"os"`:      {},
		`"strings"`: {},
		`"errors"`:  {},
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

func getAllTopLevelStructDefinitions(node *ast.File) map[string]*ast.StructType {
	allTopLevelStructDefinitions := map[string]*ast.StructType{}
	for _, decl := range node.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.TYPE {
			continue
		}
		for _, spec := range gen.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			// only care about struct type definitions for now
			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				continue
			}
			// add it to the list of definitions
			allTopLevelStructDefinitions[ts.Name.Name] = st
		}
	}
	return allTopLevelStructDefinitions
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

var goTemplate = template.Must(template.New("config").Parse(`// Code generated by configgen.go; EDIT AT YOUR OWN RISK.

{{- if .TestBuildTag }}
//go:build {{ .TestBuildTag }}
// +build {{ .TestBuildTag }}

{{- end }}

package {{ .PackageName }}

import (
    {{ .ImportList }}
)

const (
{{- range .Fields }}
    {{ .EnvVar }}_ENV = "{{ .EnvVar }}"
{{- end }}
)

var (
{{- range $f := .Fields }}
{{- range $f.ErrorVars }}
    {{ . }} = errors.New({{ $f.EnvVar }}_ENV)
{{- end }}
{{- end }}
)

func Load{{ .StructName }}() ({{ .StructName }}, error) {
    var config {{ .StructName }}
    var missingVars []error
    var formatVars []error

{{- range .Fields }}
    {{ .AssignmentName }}, ok := os.LookupEnv({{ .EnvVar }}_ENV)
    if !ok {
        missingVars = append(missingVars, {{ index .ErrorVars 0 }})
    } else {
        {{- if eq .ParseFunc "raw" }}
        config.{{ .Name }} = {{ .AssignmentName }}
        {{- else if eq .ParseFunc "strconv.Atoi" }}
        parsed, err := strconv.Atoi({{ .AssignmentName }})
        if err != nil {
            formatVars = append(formatVars, {{ index .ErrorVars 1 }})
        } else {
            config.{{ .Name }} = {{ if .CastFunc }}{{ .CastFunc }}(parsed){{ else }}parsed{{ end }}
        }
        {{- else if or (eq .ParseFunc "strconv.ParseInt") (eq .ParseFunc "strconv.ParseUint") }}
        parsed, err := {{ .ParseFunc }}({{ .AssignmentName }}, 10, {{ .BitSize }})
        if err != nil {
            formatVars = append(formatVars, {{ index .ErrorVars 1 }})
        } else {
            config.{{ .Name }} = {{ if .CastFunc }}{{ .CastFunc }}(parsed){{ else }}parsed{{ end }}
        }
		{{- else if eq .ParseFunc "strconv.ParseFloat" }}
        parsed, err := strconv.ParseFloat({{ .AssignmentName }}, {{ .BitSize }})
        if err != nil {
            formatVars = append(formatVars, {{ index .ErrorVars 1 }})
        } else {
            config.{{ .Name }} = {{ if .CastFunc }}{{ .CastFunc }}(parsed){{ else }}parsed{{ end }}
        }
        {{- else }}
        parsed, err := {{ .ParseFunc }}({{ .AssignmentName }})
        if err != nil {
            formatVars = append(formatVars, {{ index .ErrorVars 1 }})
        } else {
            config.{{ .Name }} = parsed
        }
        {{- end }}
    }
{{- end }}

    if len(missingVars) > 0 || len(formatVars) > 0 {
        var verr error
        if len(missingVars) > 0 {
            verr = errors.Join(verr, MissingEnvVarsError{vars: missingVars})
        }
        if len(formatVars) > 0 {
            verr = errors.Join(verr, InvalidEnvVarsError{vars: missingVars})
        }
        return {{ .StructName }}{}, verr
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
`))
