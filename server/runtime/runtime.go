package runtime

import (
	"context"
	"io"
	"os"
	"os/exec"

	"github.com/creack/pty"
)

type LocalRuntime struct {
	cmd *exec.Cmd
	pty *os.File
}

func NewLocalRuntime() *LocalRuntime {
	return &LocalRuntime{}
}

func (r *LocalRuntime) Start(ctx context.Context, command string, args ...string) error {
	r.cmd = exec.CommandContext(ctx, command, args...)

	// Start the command with a PTY
	ptmx, err := pty.Start(r.cmd)
	if err != nil {
		return err
	}
	r.pty = ptmx

	return nil
}

func (r *LocalRuntime) Write(p []byte) (n int, err error) {
	if r.pty == nil {
		return 0, io.ErrClosedPipe
	}
	return r.pty.Write(p)
}

func (r *LocalRuntime) Read(p []byte) (n int, err error) {
	if r.pty == nil {
		return 0, io.ErrClosedPipe
	}
	return r.pty.Read(p)
}

func (r *LocalRuntime) Close() error {
	if r.pty != nil {
		r.pty.Close()
	}
	if r.cmd != nil && r.cmd.Process != nil {
		return r.cmd.Process.Kill()
	}
	return nil
}
