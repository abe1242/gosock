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

func Client(host, port string, contnue, clip bool, outfile string) {
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
		IsStdinRead bool
		FileSize    int64
		FileNameLen uint16
		FileName    string
		StartFrom   int64 = 0
	)

	// Recieve header variables
	binary.Read(conn, binary.BigEndian, &IsStdinRead)
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
	if outfile != "" {
		FileName = outfile
	}

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

	// Setup progress bar
	var maxBytes int64 = -1
	if !IsStdinRead {
		maxBytes = FileSize - StartFrom
	}
	bar := progressbar.DefaultBytes(
		maxBytes,
		"Downloading",
	)

	// Check whether the clip option is present
	// If yes, then write to clipboard then exit.
	var clipWriter io.WriteCloser
	clipCmd, clipWriter := clipSet()
	if clip {
		_, err = io.Copy(io.MultiWriter(clipWriter, bar), conn)
		checkExit(err)
		clipWriter.Close()
		clipCmd.Wait()
		fmt.Fprintf(os.Stderr, "File '%v' copied to clipboard\n", FileName)
		os.Exit(1)
	}

	// Open file
	if IsStdinRead && FileName == "-" && outfile == "" {
		fmt.Fprintln(os.Stderr, "Data from stdin is being sent!\nPlease specify an output filename with -o option")
		os.Exit(1)
	}
	var f *os.File
	if outfile != "-" {
		openflags := os.O_CREATE | os.O_WRONLY
		if !contnue {
			openflags |= os.O_TRUNC
		}
		f, err = os.OpenFile(FileName, openflags, 0644)
		checkExit(err)
		s, err := f.Seek(StartFrom, io.SeekStart)
		checkExit(err)
		if contnue {
			fmt.Fprintf(os.Stderr, "Resuming download from %v bytes\n", s)
		}
		defer f.Close()
	} else {
		f = os.Stdout
	}

	// Copying the data to file
	_, err = io.Copy(io.MultiWriter(f, bar), conn)
	checkExit(err)
	fmt.Fprintf(os.Stderr, "File '%v' downloaded successfully\n", FileName)
}
