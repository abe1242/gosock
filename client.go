package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func client() {
	const (
		FILENAME string = "file.out"
		HOST            = "localhost"
		PORT            = "8888"
		BUFSIZE  int    = 1024
	)

	// Establishing connection with server
	conn, err := net.Dial("tcp", HOST+":"+PORT)
	check(err)
	defer conn.Close()

	///////////////////////
	var (
		FileSize    int64
		StartFrom   int64
		FileNameLen uint16
		FileName    string
	)

	binary.Read(conn, binary.BigEndian, &FileSize)
	binary.Read(conn, binary.BigEndian, &StartFrom)
	binary.Read(conn, binary.BigEndian, &FileNameLen)

	buf := make([]byte, FileNameLen)
	bytesrecieved := 0
	for {
		if bytesrecieved < int(FileNameLen) {
			n, err := conn.Read(buf[bytesrecieved:])
			if err == io.EOF {
				log.Fatal("This is wrang")
			}
			check(err)
			fmt.Printf("Received %v bytes\n", n)

			bytesrecieved += n
		} else {
			break
		}
	}
	FileName = string(buf)
	/////////////////////////

	f, err := os.OpenFile(FileName, os.O_CREATE|os.O_WRONLY, 0666)
	check(err)
	defer f.Close()

	// Getting the file data
	buf = make([]byte, BUFSIZE)
	for {
		n, err := conn.Read(buf)
		if err == io.EOF {
			break
		}
		check(err)
		f.Write(buf[:n])
	}

}
