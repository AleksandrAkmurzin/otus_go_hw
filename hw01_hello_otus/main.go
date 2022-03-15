package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	initialString := "Hello, OTUS!"
	reversedString := stringutil.Reverse(initialString)
	fmt.Println(reversedString)
}
