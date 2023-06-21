package entity

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/fs"
	"net"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

func Server(fpath, host, port string, clip bool) {
	// Should the file be readfrom stdin
	rstdin := fpath == "-"

	// Check if filename is provided if 'clip' option is not set
	if !clip && fpath == "" {
		fmt.Fprintln(os.Stderr, "Error: filename not provided")
		os.Exit(1)
	}

	// Listening for connections
	s, err := net.Listen("tcp", host+":"+port)
	checkExit(err)
	fmt.Printf("gosock listening on %v at port %v\n", host, port)
	fmt.Println("----------------")
	printIPs()
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
		var f *os.File
		var clipReader io.ReadCloser
		if !rstdin && !clip {
			f, err = os.Open(fpath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				continue
			}
		} else if clip {
			clipReader = clipGet()
			if clipReader == nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", "could not read from clipboard")
				os.Exit(1)
			}
		} else {
			f = os.Stdin
		}

		// Setting up header variables
		var fileinfo fs.FileInfo
		if !rstdin && !clip {
			fileinfo, err = f.Stat()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				continue
			}
		}

		var (
			IsStdinRead bool = true
			FileSize    int64
			FileName    string = filepath.Base(fpath)
			FileNameLen uint16 = uint16(len([]byte(FileName)))
			StartFrom   int64
		)
		if fileinfo != nil {
			IsStdinRead = false
			FileSize = fileinfo.Size()
		}
		if clip {
			FileName = "clip"
			FileNameLen = uint16(len([]byte(FileName)))
		}

		// Sending the header data
		binary.Write(conn, binary.BigEndian, IsStdinRead)
		binary.Write(conn, binary.BigEndian, FileSize)
		binary.Write(conn, binary.BigEndian, FileNameLen)
		conn.Write([]byte(FileName))

		// Read the start byte and set seek
		binary.Read(conn, binary.BigEndian, &StartFrom)
		if !rstdin && !clip {
			f.Seek(StartFrom, io.SeekStart)
		}

		// Send the file
		var maxBytes int64 = -1
		if !rstdin && !clip {
			maxBytes = FileSize - StartFrom
		}
		bar := progressbar.DefaultBytes(
			maxBytes,
			"Sending",
		)
		if clip {
			_, err = io.Copy(io.MultiWriter(conn, bar), clipReader)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				continue
			}
		} else {
			_, err = io.Copy(io.MultiWriter(conn, bar), f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				continue
			}
		}
		fmt.Printf("File '%v' sent to %v successfully\n", FileName, conn.RemoteAddr())

		f.Close()
		conn.Close()

		if rstdin {
			os.Exit(1)
		}
	}

}

// prints all the local ipv4 addresses
func printIPs() {
	ifs, err := net.Interfaces()
	checkExit(err)
	for _, i := range ifs {
		addrs, err := i.Addrs()
		checkExit(err)

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip.To4() != nil && !ip.IsLoopback() {
				fmt.Printf("%v - %v\n", ip, i.Name)
			}
		}
	}
	fmt.Println()
}
