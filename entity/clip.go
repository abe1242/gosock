package entity

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
)

func clipGet() io.ReadCloser {
	var cmd *exec.Cmd
	if runtime.GOOS == "android" {
		cmd = exec.Command("termux-clipboard-get")
	} else if runtime.GOOS == "linux" {
		cmd = exec.Command("xclip", "-sel", "clip", "-o")
	} else {
		fmt.Fprintln(os.Stderr, "Your OS doesn't support this option (--clipboard)")
		os.Exit(1)
	}
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil
	}
	if err = cmd.Start(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil
	}
	return stdout
}

func clipSet() (*exec.Cmd, io.WriteCloser) {
	var cmd *exec.Cmd
	if runtime.GOOS == "android" {
		cmd = exec.Command("termux-clipboard-set")
	} else if runtime.GOOS == "linux" {
		cmd = exec.Command("xclip", "-sel", "clip")
	} else {
		fmt.Fprintln(os.Stderr, "Your OS doesn't support this option (--clipboard)")
		os.Exit(1)
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, nil
	}
	if err = cmd.Start(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, nil
	}
	return cmd, stdin
}
