package main

import (
	"fmt"
	"log"
	"net"
)

func server() {
	s, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	for {
		conn, err := s.Accept()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Connection from (%s)\n", conn.RemoteAddr())

		conn.Write([]byte("Hello"))

		conn.Close()
	}

}
