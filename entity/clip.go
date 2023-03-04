package entity

import (
	"os/exec"
	"io"
)

func clipGet() io.ReadCloser {
    cmd := exec.Command("termux-clipboard-get")
    stdout, err := cmd.StdoutPipe()
    cmd.Stderr = cmd.Stdout
    if err != nil {
        return nil
    }
    if err = cmd.Start(); err != nil {
        return nil
    }
    return stdout
}

func clipSet() io.WriteCloser {
    cmd := exec.Command("termux-clipboard-set")
    stdin, err := cmd.StdinPipe()
    if err != nil {
        return nil
    }
    if err = cmd.Start(); err != nil {
        return nil
    }
    return stdin
}
