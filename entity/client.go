package entity

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

func Client(host, port string, contnue bool) {
	// Establishing connection with server
    var conn net.Conn
    var err error
    for {
        conn, err = net.Dial("tcp", host+":"+port)
        if err == nil {
            break
        }
        time.Sleep(1 * time.Second)
    }
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
	checkExit(err)
	if n != int(FileNameLen) {
		fmt.Fprintf(os.Stderr, "Error: Filename not received fully\n")
	}
	FileName = string(buf)

	// Get the filesize and set start byte if continue flag is present
	if contnue {
		f, err := os.Open(FileName)
		checkExit(err)
		fi, err := f.Stat()
		checkExit(err)
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
	checkExit(err)
	s, err := f.Seek(StartFrom, io.SeekStart)
	checkExit(err)
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
	checkExit(err)
	fmt.Printf("File '%v' downloaded successfully\n", FileName)
}
