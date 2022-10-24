package entity

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/schollz/progressbar/v3"
)

func Client(host, port string, contnue bool) {
	// Establishing connection with server
	conn, err := net.Dial("tcp", host+":"+port)
	check(err)
	defer conn.Close()

	// Define some header variables
	var (
		FileSize    int64
		FileNameLen uint16
		FileName    string
		StartFrom   int64 = 0
	)

	// Recieve header variables
	binary.Read(conn, binary.BigEndian, &FileSize)
	binary.Read(conn, binary.BigEndian, &FileNameLen)

	// Get the filename
	// Read as many bytes into buf as the FileNameLen variable
	buf := make([]byte, FileNameLen)
	n, err := io.ReadFull(conn, buf)
	check(err)
	if n != int(FileNameLen) {
		fmt.Fprintf(os.Stderr, "Error: Filename not received fully\n")
	}
	FileName = string(buf)

	// Get the filesize and set start byte if continue flag is present
	if contnue {
		f, err := os.Open(FileName)
		check(err)
		fi, err := f.Stat()
		check(err)
		StartFrom = fi.Size()
		f.Close()
	}

	// Send the position to start from
	binary.Write(conn, binary.BigEndian, StartFrom)

	// Open file
	openflags := os.O_CREATE | os.O_WRONLY
	if !contnue {
		openflags |= os.O_TRUNC
	}
	f, err := os.OpenFile(FileName, openflags, 0644)
	check(err)
	s, err := f.Seek(StartFrom, io.SeekStart)
	check(err)
	if contnue {
		fmt.Printf("Resuming download from %v bytes\n", s)
	}
	defer f.Close()

	// Copying the data to file
	bar := progressbar.DefaultBytes(
		FileSize-StartFrom,
		"Downloading",
	)
	_, err = io.Copy(io.MultiWriter(f, bar), conn)
	check(err)
}
