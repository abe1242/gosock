package entity

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

func Server(fpath, host, port string) {
	// Listening for connections
	s, err := net.Listen("tcp", host+":"+port)
	checkExit(err)
	fmt.Printf("Listening for connections at %v:%v\n", host, port)
	defer s.Close()

	for {
		// Accepting connection
		conn, err := s.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			continue
		}
		fmt.Printf("\nConnection from (%s)\n", conn.RemoteAddr())

		// Opening file to read from
		f, err := os.Open(fpath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			continue
		}

		// Setting up header variables
		fileinfo, err := f.Stat()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			continue
		}
		var (
			FileSize    int64  = fileinfo.Size()
			FileName    string = filepath.Base(fpath)
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
		bar := progressbar.DefaultBytes(
			FileSize-StartFrom,
			"Sending",
		)
		_, err = io.Copy(io.MultiWriter(conn, bar), f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			continue
		}
		fmt.Printf("File '%v' sent to %v successfully\n", FileName, conn.RemoteAddr())

		f.Close()
		conn.Close()
	}

}
