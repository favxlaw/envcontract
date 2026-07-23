package cli

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type sourceField struct {
	Name       string
	EnvKey     string
	Required   bool
	HasDefault bool
	Default    string
	Kind       string
}

type sourceStruct struct {
	Name   string
	Fields []sourceField
}

func parseGoFile(path string) ([]sourceStruct, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}

	var structs []sourceStruct

	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			s := sourceStruct{Name: typeSpec.Name.Name}
			for _, field := range structType.Fields.List {
				if len(field.Names) == 0 {
					continue
				}

				tag := ""
				if field.Tag != nil {
					tag = field.Tag.Value
					tag = strings.Trim(tag, "`")
				}

				envVal, ok := parseStructTag(tag, "env")
				if !ok || envVal == "" || envVal == "-" {
					continue
				}

				ft := resolveFieldType(field.Type)
				if ft == "" {
					continue
				}

				parsed := parseEnvTagValue(envVal)
				s.Fields = append(s.Fields, sourceField{
					Name:       field.Names[0].Name,
					EnvKey:     parsed.key,
					Required:   parsed.required,
					HasDefault: parsed.hasDefault,
					Default:    parsed.defaultValue,
					Kind:       ft,
				})
			}

			if len(s.Fields) > 0 {
				structs = append(structs, s)
			}
		}
	}

	if len(structs) == 0 {
		return nil, fmt.Errorf("no structs with env tags found in %s", path)
	}

	return structs, nil
}

func resolveFieldType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return mapGoKind(t.Name)
	case *ast.StarExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			return mapGoKind(ident.Name)
		}
		if sel, ok := t.X.(*ast.SelectorExpr); ok {
			return resolveSelector(sel)
		}
		return ""
	case *ast.SelectorExpr:
		return resolveSelector(t)
	case *ast.ArrayType:
		return ""
	default:
		return ""
	}
}

func resolveSelector(sel *ast.SelectorExpr) string {
	if pkg, ok := sel.X.(*ast.Ident); ok {
		if pkg.Name == "time" && sel.Sel.Name == "Duration" {
			return "duration"
		}
	}
	return ""
}

func mapGoKind(name string) string {
	switch name {
	case "string":
		return "string"
	case "int":
		return "int"
	case "int64":
		return "int64"
	case "float64":
		return "float64"
	case "bool":
		return "bool"
	default:
		return ""
	}
}

type envTagValue struct {
	key          string
	required     bool
	hasDefault   bool
	defaultValue string
}

func parseEnvTagValue(tag string) envTagValue {
	parts := strings.Split(tag, ",")
	if len(parts) == 0 || parts[0] == "" {
		return envTagValue{}
	}

	parsed := envTagValue{key: strings.TrimSpace(parts[0])}

	for _, rawOpt := range parts[1:] {
		opt := strings.TrimSpace(rawOpt)
		switch {
		case opt == "":
		case opt == "required":
			parsed.required = true
		case strings.HasPrefix(opt, "default="):
			parsed.hasDefault = true
			parsed.defaultValue = strings.TrimPrefix(opt, "default=")
		}
	}

	return parsed
}

func parseStructTag(tag, key string) (string, bool) {
	for tag != "" {
		i := strings.Index(tag, key+":")
		if i < 0 {
			break
		}

		remaining := tag[i+len(key)+1:]
		if remaining == "" {
			break
		}

		if remaining[0] != '"' {
			tag = remaining
			continue
		}

		remaining = remaining[1:]
		j := strings.IndexByte(remaining, '"')
		if j < 0 {
			break
		}

		return remaining[:j], true
	}

	return "", false
}
