package parser

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/favxlaw/envcontract"
)

func ParseStruct(v any) ([]envcontract.FieldContract, error) {
	if v == nil {
		return nil, fmt.Errorf("envcontract: input must be a non-nil pointer to a struct")
	}

	rt := reflect.TypeOf(v)
	if rt.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("envcontract: input must be a pointer to a struct, got %s", rt.Kind())
	}

	if rt.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("envcontract: input must be a pointer to a struct, got pointer to %s", rt.Elem().Kind())
	}

	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return nil, fmt.Errorf("envcontract: input must be a non-nil pointer to a struct")
	}

	return parseStructType(rt.Elem())
}

func parseStructType(rt reflect.Type) ([]envcontract.FieldContract, error) {
	var contracts []envcontract.FieldContract

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)

		if field.PkgPath != "" {
			continue
		}

		tag, hasTag := field.Tag.Lookup("env")
		if hasTag && tag == "-" {
			continue
		}

		fieldType := field.Type
		fieldKind := fieldType.Kind()

		if fieldKind == reflect.Ptr {
			fieldType = fieldType.Elem()
			fieldKind = fieldType.Kind()
		}

		if fieldKind == reflect.Struct && !hasTag {
			nested, err := parseStructType(fieldType)
			if err != nil {
				return nil, err
			}
			contracts = append(contracts, nested...)
			continue
		}

		if !hasTag {
			continue
		}

		envTag, err := parseEnvTag(tag)
		if err != nil {
			return nil, fmt.Errorf("field %s: %w", field.Name, err)
		}

		if envTag.key == "" {
			continue
		}

		if fieldKind == reflect.Struct {
			nested, err := parseStructType(fieldType)
			if err != nil {
				return nil, err
			}
			contracts = append(contracts, nested...)
			continue
		}

		kind, ok := supportedKind(fieldKind)
		if !ok {
			continue
		}

		contracts = append(contracts, envcontract.FieldContract{
			Name:       field.Name,
			EnvKey:     envTag.key,
			Required:   envTag.required,
			HasDefault: envTag.hasDefault,
			Default:    envTag.defaultValue,
			Kind:       kind,
		})
	}

	return contracts, nil
}

type envTag struct {
	key          string
	required     bool
	hasDefault   bool
	defaultValue string
}

func parseEnvTag(tag string) (envTag, error) {
	parts := strings.Split(tag, ",")
	if len(parts) == 0 {
		return envTag{}, nil
	}

	parsed := envTag{
		key: strings.TrimSpace(parts[0]),
	}

	for _, rawOpt := range parts[1:] {
		opt := strings.TrimSpace(rawOpt)

		switch {
		case opt == "":
			continue
		case opt == "required":
			parsed.required = true
		case strings.HasPrefix(opt, "default="):
			parsed.hasDefault = true
			parsed.defaultValue = strings.TrimPrefix(opt, "default=")
		default:
			return envTag{}, fmt.Errorf("unknown env tag option %q", opt)
		}
	}

	return parsed, nil
}

func supportedKind(k reflect.Kind) (string, bool) {
	switch k {
	case reflect.String:
		return "string", true
	case reflect.Int:
		return "int", true
	case reflect.Int64:
		return "int64", true
	case reflect.Float64:
		return "float64", true
	case reflect.Bool:
		return "bool", true
	default:
		return "", false
	}
}
