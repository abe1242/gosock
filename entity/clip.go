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
