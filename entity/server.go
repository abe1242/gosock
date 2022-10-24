package entity

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path"
)

func Server(filepath, host, port string) {
	// Listening for connections
	s, err := net.Listen("tcp", host+":"+port)
	check(err)
	fmt.Printf("Listening for connections at %v:%v\n", host, port)
	defer s.Close()

	for {
		// Accepting connection
		conn, err := s.Accept()
		check(err)
		fmt.Printf("Connection from (%s)\n", conn.RemoteAddr())

		// Opening file to read from
		f, err := os.Open(filepath)
		check(err)

		// Setting up header variables
		fileinfo, err := f.Stat()
		check(err)
		var (
			FileSize    int64  = fileinfo.Size()
			FileName    string = path.Base(filepath)
			FileNameLen uint16 = uint16(len([]byte(FileName)))
			StartFrom   int64
		)

		// Sending the header data
		binary.Write(conn, binary.BigEndian, FileSize)
		binary.Write(conn, binary.BigEndian, FileNameLen)
		conn.Write([]byte(FileName))

		// Read the start byte and set seek
		binary.Read(conn, binary.BigEndian, &StartFrom)
		f.Seek(StartFrom, io.SeekStart)

		// Send the file
		_, err = io.Copy(conn, f)
		check(err)

		f.Close()
		conn.Close()
	}

}
