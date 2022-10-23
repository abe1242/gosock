package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path"
)

func server() {
	const (
		FILEPATH string = "go.mod"
		HOST            = "0.0.0.0"
		PORT            = "8888"
		BUFSIZE  int    = 1024
	)

	s, err := net.Listen("tcp", HOST+":"+PORT)
	check(err)
	fmt.Printf("Listening for connections at %v:%v\n", HOST, PORT)
	defer s.Close()

	for {
		conn, err := s.Accept()
		check(err)
		fmt.Printf("Connection from (%s)\n", conn.RemoteAddr())

		f, err := os.Open(FILEPATH)
		check(err)

		///////////////////
		fi, err := f.Stat()
		check(err)
		var FileSize int64 = fi.Size()
		var StartFrom int64 = 0
		var FileName string = path.Base(FILEPATH)
		var FileNameLen uint16 = uint16(len([]byte(FileName)))

		binary.Write(conn, binary.BigEndian, FileSize)
		binary.Write(conn, binary.BigEndian, StartFrom)
		binary.Write(conn, binary.BigEndian, FileNameLen)

		conn.Write([]byte(FileName))
		//////////////////////

		buf := make([]byte, BUFSIZE)
		for {
			n, err := f.Read(buf)
			if err == io.EOF {
				break
			}
			check(err)

			conn.Write(buf[:n])
		}

		f.Close()
		conn.Close()
	}

}
