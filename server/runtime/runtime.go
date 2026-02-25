package runtime

import (
	"context"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/creack/pty"
)

type LocalRuntime struct {
	cmd     *exec.Cmd
	pty     *os.File
	workDir string
	bash    *BashSession
}

func NewLocalRuntime() *LocalRuntime {
	// Default to workspace/ in current directory, or env var
	wd := os.Getenv("WORKSPACE_BASE")
	if wd == "" {
		wd = "workspace"
	}
	// Ensure directory exists
	os.MkdirAll(wd, 0755)

	bash, _ := NewBashSession(wd)

	return &LocalRuntime{
		workDir: wd,
		bash:    bash,
	}
}

func (r *LocalRuntime) Start(ctx context.Context, command string, args ...string) error {
	r.cmd = exec.CommandContext(ctx, command, args...)
	r.cmd.Dir = r.workDir
	// Set safe environment variables?
	// For now inherit but verify workDir is set.

	// Start the command with a PTY
	ptmx, err := pty.Start(r.cmd)
	if err != nil {
		return err
	}
	r.pty = ptmx

	return nil
}

func (r *LocalRuntime) Execute(ctx context.Context, command string, args ...string) (string, int, error) {
	// Check if command is bash -c which is common from ActionService
	cmdStr := command
	if command == "bash" && len(args) >= 2 && args[0] == "-c" {
		cmdStr = args[1]
	} else if len(args) > 0 {
		// Attempt to reconstruct command string
		cmdStr = command + " " + strings.Join(args, " ")
	}

	return r.bash.Execute(ctx, cmdStr)
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
	var firstErr error
	if r.bash != nil {
		if err := r.bash.Close(); err != nil {
			firstErr = err
		}
	}

	if r.pty != nil {
		r.pty.Close()
	}
	if r.cmd != nil && r.cmd.Process != nil {
		if err := r.cmd.Process.Kill(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
