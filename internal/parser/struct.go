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

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("envcontract: input must be a pointer to a struct, got %s", rv.Kind())
	}
	if rv.IsNil() {
		return nil, fmt.Errorf("envcontract: input must be a non-nil pointer to a struct")
	}

	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("envcontract: input must be a pointer to a struct, got pointer to %s", rv.Kind())
	}

	return parseStruct(rv)
}

func parseStruct(rv reflect.Value) ([]envcontract.FieldContract, error) {
	rt := rv.Type()
	var contracts []envcontract.FieldContract

	for i := range rt.NumField() {
		field := rt.Field(i)
		fieldVal := rv.Field(i)

		tag, ok := field.Tag.Lookup("env")
		if !ok {
			kind := fieldVal.Kind()
			if kind == reflect.Struct {
				nested, err := parseStruct(fieldVal)
				if err != nil {
					return nil, err
				}
				contracts = append(contracts, nested...)
				continue
			}
			if kind == reflect.Ptr && !fieldVal.IsNil() && fieldVal.Elem().Kind() == reflect.Struct {
				nested, err := parseStruct(fieldVal.Elem())
				if err != nil {
					return nil, err
				}
				contracts = append(contracts, nested...)
			}
			continue
		}

		parts := strings.Split(tag, ",")
		envKey := parts[0]
		var required bool
		var hasDefault bool
		var defaultVal string

		for _, opt := range parts[1:] {
			if opt == "required" {
				required = true
			} else if strings.HasPrefix(opt, "default=") {
				hasDefault = true
				defaultVal = strings.TrimPrefix(opt, "default=")
			}
		}
		kind := fieldVal.Kind()
		if kind == reflect.Struct {
			nested, err := parseStruct(fieldVal)
			if err != nil {
				return nil, err
			}
			contracts = append(contracts, nested...)
			continue
		}
		if kind == reflect.Ptr {
			if fieldVal.IsNil() {
				continue
			}
			fieldVal = fieldVal.Elem()
			kind = fieldVal.Kind()
		}

		kindStr, ok := supportedKind(kind)
		if !ok {
			continue
		}

		contracts = append(contracts, envcontract.FieldContract{
			Name:       field.Name,
			EnvKey:     envKey,
			Required:   required,
			HasDefault: hasDefault,
			Default:    defaultVal,
			Kind:       kindStr,
		})

	}

	return contracts, nil
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
	}
	return "", false
}
