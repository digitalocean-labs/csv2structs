package csv2structs

import (
	"fmt"
	"reflect"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// snakeToTitle converts a snake_case string to a TitleCase string
func snakeToTitle(header string) string {
	caser := cases.Title(language.English)

	var cleaned []string
	for _, part := range strings.Split(header, "_") {
		cleaned = append(cleaned, caser.String(part))
	}

	return strings.Join(cleaned, "")
}

// getHeader will apply the HeaderType option to the header or HeaderTransform if provided
func getHeader(opts *options, header []string) []string {
	if opts.HeaderTransform != nil {
		return remap(header, opts.HeaderTransform)
	}

	switch opts.HeaderType {
	case HeaderTypeSnake:
		return remap(header, snakeToTitle)
	default:
		return header
	}
}

func mapHeader(fields []reflect.StructField, header []string) (map[string]int, error) {
	var foundHeaders = map[string]int{}
	var missingHeaders []string
	for _, f := range fields {
		for i, h := range header {
			if f.Name == h {
				foundHeaders[h] = i
				break
			}
		}
		if _, ok := foundHeaders[f.Name]; !ok {
			missingHeaders = append(missingHeaders, f.Name)
		}
	}

	for _, h := range header {
		if _, ok := foundHeaders[h]; !ok {
			missingHeaders = append(missingHeaders, h)
		}
	}

	if len(missingHeaders) > 0 || len(foundHeaders) != len(fields) {
		return nil, fmt.Errorf(
			"missing header%s: %v",
			map[bool]string{true: "s", false: ""}[len(missingHeaders) > 1],
			strings.Join(missingHeaders, ", "),
		)
	}

	return foundHeaders, nil
}
