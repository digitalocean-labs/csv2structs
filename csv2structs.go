// Package csv2structs parses CSV data into a slice of structs.
//
// # Example Usage
//
//	package main
//
//	import (
//		"fmt"
//		"strings"
//
//		"github.com/digitalocean-labs/csv2structs"
//	)
//
//	func main() {
//		csvData := `first_name,age
//	Alice,30
//	Bob,25
//	Charlie,35`
//
//		type Person struct {
//			FirstName string
//			Age       int
//		}
//
//		r := strings.NewReader(csvData)
//		people, err := csv2structs.Parse[Person](r)
//		if err != nil {
//			fmt.Println("error:", err)
//			return
//		}
//
//		for _, p := range people {
//			fmt.Printf("%+v\n", p)
//		}
//	}
//
// # Headers
//
// All exported fields in the struct passed must match the headers in the CSV data.
//
// By default, the headers in the CSV data are transformed from snake_case to TitleCase.
//
// If you want to disable the header transformation, you can use the WithHeaderType option with the HeaderTypeNone value:
//
//	people, err := csv2structs.Parse[Person](r, csv2structs.WithHeaderType(csv2structs.HeaderTypeNone))
//
// If your CSV data has headers in a different format, you can implement your own
// header transformation function and pass it to the WithHeaderTransform option:
//
//	func customHeaderTransform(header string) string {
//		// your custom header transformation logic
//	}
//
//	people, err := csv2structs.Parse[Person](r, csv2structs.WithHeaderTransform(customHeaderTransform))
//
// Or, if your CSV data has headers in snake_case format and you want to be explicit,
// you can use the WithHeaderType option with the HeaderTypeSnake value:
//
//	people, err := csv2structs.Parse[Person](r, csv2structs.WithHeaderType(csv2structs.HeaderTypeSnake))
package csv2structs

import (
	"io"
)

// Parse parses a CSV and returns a slice of structs
func Parse[T any](reader io.Reader, opts ...Option) ([]*T, error) {
	p, err := NewParser[T](reader, opts...)
	if err != nil {
		return nil, err
	}

	return p.ReadAll()
}
