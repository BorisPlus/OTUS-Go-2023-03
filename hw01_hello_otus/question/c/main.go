package main

import (
	"fmt"

	"github.com/BorisPlus/stringutil"
)

func main() {
	row := "Hello, OTUS!"
	reversedRow := stringutil.Reverse(row)
	fmt.Println(reversedRow)
}
