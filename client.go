package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func client() {
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	b, err := io.ReadAll(conn)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
