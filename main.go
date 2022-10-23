package main

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Error: Invalid number of arguments\n")
	} else if os.Args[1] == "c" {
		client()
	} else {
		server(os.Args[1])
	}
}

func check(err error) {
	if err != nil {
		// panic(err)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
