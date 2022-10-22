package main

import (
	"flag"
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
