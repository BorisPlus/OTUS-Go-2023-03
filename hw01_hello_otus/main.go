package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func Echo(a string) string {
	return a
}

func main() {
	row := "Hello, OTUS!"
	reversedRow := stringutil.Reverse(row)
	fmt.Println(reversedRow)
}
