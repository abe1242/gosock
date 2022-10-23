package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Parse()

	arg := flag.Arg(0)

	if arg == "c" {
		client()
	} else {
		server()
	}
}

func check(err error) {
	if err != nil {
		// panic(err)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
