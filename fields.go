package csv2structs

import (
	"errors"
	"reflect"
)

var (
	errNoVisibleFields  = errors.New("no visible fields")
	errNoExportedFields = errors.New("no exported fields")
	errInvalidType      = errors.New("invalid type; must be a struct")
)

func getFields[T any]() ([]reflect.StructField, error) {
	var t T

	tt := reflect.TypeOf(t)
	if tt == nil || tt.Kind() != reflect.Struct {
		return nil, errInvalidType
	}

	allFields := reflect.VisibleFields(tt)
	if len(allFields) == 0 {
		return nil, errNoVisibleFields
	}

	var fields []reflect.StructField
	for _, field := range allFields {
		if field.IsExported() {
			fields = append(fields, field)
		}
	}
	if len(fields) == 0 {
		return nil, errNoExportedFields
	}

	return fields, nil
}
