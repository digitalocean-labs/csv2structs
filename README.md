# CSV2Structs

[![GoDoc](https://godoc.org/github.com/digitalocean-labs/csv2structs?status.svg)](https://godoc.org/github.com/digitalocean-labs/csv2structs)
[![Go Report Card](https://goreportcard.com/badge/github.com/digitalocean-labs/csv2structs)](https://goreportcard.com/report/github.com/digitalocean-labs/csv2structs)

Package csv2structs parses CSV data into a slice of structs.

## Example Usage

```go
package main

import (
	"fmt"
	"strings"

	"github.com/digitalocean-labs/csv2structs"
)

func main() {
	csvData := `first_name,age
Alice,30
Bob,25
Charlie,35`

	type Person struct {
		FirstName string
		Age       int
	}

	r := strings.NewReader(csvData)
	people, err := csv2structs.Parse[Person](r)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	for _, p := range people {
		fmt.Printf("%+v\n", p)
	}
}
```

The above is [available as a runnable example](demo/demo.go).


## Headers

All exported fields in the struct passed must match the headers in the CSV data. 

By default, the headers in the CSV data are transformed from snake_case to TitleCase.

If you want to disable the header transformation, you can use the WithHeaderType option with the HeaderTypeNone value:

```go
people, err := csv2structs.Parse[Person](r, csv2structs.WithHeaderType(csv2structs.HeaderTypeNone))
```

If your CSV data has headers in a different format, you can implement your own header transformation function and pass it to the WithHeaderTransform option:

```go
func customHeaderTransform(header string) string {
    // your custom header transformation logic
}

people, err := csv2structs.Parse[Person](r, csv2structs.WithHeaderTransform(customHeaderTransform))
```

Or, if your CSV data has headers in snake_case format and you want to be explicit, you can use the WithHeaderType option with the HeaderTypeSnake value:

```go
people, err := csv2structs.Parse[Person](r, csv2structs.WithHeaderType(csv2structs.HeaderTypeSnake))
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

