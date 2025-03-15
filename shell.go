package benchtune

import (
	"context"
	"fmt"
	"io"
	"os/exec"
)

func NewShell(ctx context.Context) (error, io.WriteCloser, io.ReadCloser) {
	cmd := exec.CommandContext(ctx, "/bin/sh", "-i")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("NewShell: %w", err), nil, nil
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("NewShell: %w", err), nil, nil
	}

	cmd.Stderr = cmd.Stdout

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("NewShell: %w", err), nil, nil
	}

	buf := make([]byte, 1024)
	n, err := stdout.Read(buf)
	if err != nil {
		return fmt.Errorf("NewShell: %w", err), nil, nil
	}

	if string(buf[:n]) != "$ " {
		return fmt.Errorf("NewShell: no prompt read"), nil, nil
	}

	return nil, stdin, stdout
}
