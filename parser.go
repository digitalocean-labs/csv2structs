package csv2structs

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
	"strconv"
)

// Parser can be used to read a single or all structs from an io.Reader containing CSV data
type Parser[T any] interface {
	Read() (*T, error)
	ReadAll() ([]*T, error)
}

type parser[T any] struct {
	opts   *options
	reader *csv.Reader
	fields []reflect.StructField
	header map[string]int
}

// NewParser returns a new Parser implementation for the given io.Reader containing CSV data
func NewParser[T any](reader io.Reader, opts ...Option) (Parser[T], error) {
	fields, err := getFields[T]()
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(reader)
	header, err := r.Read()
	if err != nil {
		return nil, err
	}

	o := getOptions(opts)

	h := getHeader(o, header)

	headerMap, err := mapHeader(fields, h)
	if err != nil {
		return nil, err
	}

	p := &parser[T]{
		opts:   o,
		reader: r,
		fields: fields,
		header: headerMap,
	}

	return p, nil
}

// Read a single struct from the CSV data
func (p *parser[T]) Read() (*T, error) {
	row, err := p.reader.Read()
	if err != nil {
		return nil, err
	}

	var obj T
	for f, i := range p.header {
		field := reflect.ValueOf(&obj).Elem().FieldByName(f)

		switch field.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val, err := strconv.ParseInt(row[i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf(`failed to convert "%s" to int`, row[i])
			}
			field.SetInt(val)
		case reflect.String:
			field.SetString(row[i])
		case reflect.Bool:
			val, err := strconv.ParseBool(row[i])
			if err != nil {
				return nil, fmt.Errorf(`failed to convert "%s" to bool`, row[i])
			}
			field.SetBool(val)
		case reflect.Float64, reflect.Float32:
			val, err := strconv.ParseFloat(row[i], 64)
			if err != nil {
				return nil, fmt.Errorf(`failed to convert "%s" to float`, row[i])
			}
			field.SetFloat(val)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val, err := strconv.ParseUint(row[i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf(`failed to convert "%s" to uint`, row[i])
			}
			field.SetUint(val)
		default:
			return nil, fmt.Errorf("unsupported type: %s", field.Kind())
		}
	}

	return &obj, nil
}

// ReadAll reads all structs from the CSV data
func (p *parser[T]) ReadAll() ([]*T, error) {
	var objects []*T
	for {
		obj, err := p.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		objects = append(objects, obj)
	}

	return objects, nil
}
