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
