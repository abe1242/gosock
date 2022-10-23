package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path"
)

func server(filepath string) {
	const (
		HOST = "0.0.0.0"
		PORT = "8888"
	)

	s, err := net.Listen("tcp", HOST+":"+PORT)
	check(err)
	fmt.Printf("Listening for connections at %v:%v\n", HOST, PORT)
	defer s.Close()

	for {
		conn, err := s.Accept()
		check(err)
		fmt.Printf("Connection from (%s)\n", conn.RemoteAddr())

		f, err := os.Open(filepath)
		check(err)

		///////////////////
		fi, err := f.Stat()
		check(err)
		var (
			FileSize    int64  = fi.Size()
			FileName    string = path.Base(filepath)
			FileNameLen uint16 = uint16(len([]byte(FileName)))
			StartFrom   int64
		)

		binary.Write(conn, binary.BigEndian, FileSize)
		binary.Write(conn, binary.BigEndian, FileNameLen)

		conn.Write([]byte(FileName))

		// Read the byte to start from and set seek
		binary.Read(conn, binary.BigEndian, &StartFrom)
		f.Seek(StartFrom, io.SeekStart)

		// Send the file
		io.Copy(conn, f)

		f.Close()
		conn.Close()
	}

}
