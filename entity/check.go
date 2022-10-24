package entity

import (
	"fmt"
	"os"
)

func check(err error) {
	if err != nil {
		// panic(err)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
