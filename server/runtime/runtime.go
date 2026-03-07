package runtime

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/creack/pty"
)

type LocalRuntime struct {
	cmd     *exec.Cmd
	pty     *os.File
	workDir string
	shell   *ShellSession
}

func NewLocalRuntime() *LocalRuntime {
	// Default to workspace/ in current directory, or env var
	wd := os.Getenv("WORKSPACE_BASE")
	if wd == "" {
		wd = "workspace"
	}
	// Ensure directory exists
	os.MkdirAll(wd, 0755)

	return &LocalRuntime{
		workDir: wd,
	}
}

func (r *LocalRuntime) startLocalShell() error {
	if r.shell != nil {
		return nil
	}

	r.cmd = exec.Command("bash", "--noprofile", "--norc")
	r.cmd.Dir = r.workDir
	r.cmd.Env = os.Environ()

	f, err := pty.Start(r.cmd)
	if err != nil {
		return err
	}
	r.pty = f
	r.shell = NewShellSession(f)

	return nil
}

func (r *LocalRuntime) Start(ctx context.Context, command string, args ...string) error {
	// Legacy Start behavior: launch a command with PTY.
	// This is stateless (unless we assume this IS the shell).
	// If we use Start for One-Off commands, it conflicts with Persistent Shell.
	// We'll maintain existing behavior for now, but Execute is preferred.
	r.cmd = exec.CommandContext(ctx, command, args...)
	r.cmd.Dir = r.workDir

	ptmx, err := pty.Start(r.cmd)
	if err != nil {
		return err
	}
	r.pty = ptmx
	return nil
}

func (r *LocalRuntime) Execute(ctx context.Context, command string, args ...string) (string, int, error) {
	// Ensure shell is running
	if r.shell == nil {
		if err := r.startLocalShell(); err != nil {
			return "", -1, err
		}
	}

	// Check if command is bash -c which is common from ActionService
	cmdStr := command
	if command == "bash" && len(args) >= 2 && args[0] == "-c" {
		cmdStr = args[1]
	} else if len(args) > 0 {
		// Attempt to reconstruct command string
		cmdStr = command + " " + strings.Join(args, " ")
	}

	slog.Debug("Executing command:", "command", cmdStr)
	out, exitCode, err := r.shell.Execute(ctx, cmdStr)
	slog.Debug("Command finished", "exit_code", exitCode)

	return out, exitCode, err
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

func (r *LocalRuntime) GetCwd(ctx context.Context) (string, error) {
	if r.shell != nil {
		return r.shell.GetCwd(), nil
	}
	return r.workDir, nil
}

func (r *LocalRuntime) CopyFileToContainer(ctx context.Context, hostPath string, containerPath string) error {
	// For local runtime, it's just a local file copy
	sourceFileStat, err := os.Stat(hostPath)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", hostPath)
	}

	source, err := os.Open(hostPath)
	if err != nil {
		return err
	}
	defer source.Close()

	// Handle relative container paths using current working directory
	targetPath := containerPath
	if !filepath.IsAbs(containerPath) {
		cwd, _ := r.GetCwd(ctx)
		targetPath = filepath.Join(cwd, containerPath)
	}

	// Ensure the target directory exists
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	destination, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

func (r *LocalRuntime) CopyFileFromContainer(ctx context.Context, containerPath string, hostPath string) error {
	// For local runtime, it's just a local file copy

	// Handle relative container paths using current working directory
	sourcePath := containerPath
	if !filepath.IsAbs(containerPath) {
		cwd, _ := r.GetCwd(ctx)
		sourcePath = filepath.Join(cwd, containerPath)
	}

	sourceFileStat, err := os.Stat(sourcePath)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", sourcePath)
	}

	source, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(hostPath)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

func (r *LocalRuntime) GetVSCodeURL() *string {
	// Not natively supported in basic mock local runtime yet.
	return nil
}

func (r *LocalRuntime) GetWebHosts() map[string]interface{} {
	// Not natively supported in basic mock local runtime yet.
	return make(map[string]interface{})
}

func (r *LocalRuntime) Close() error {
	var firstErr error
	if r.shell != nil {
		if err := r.shell.Close(); err != nil {
			firstErr = err
		}
	} else if r.pty != nil {
		// If shell not initialized (legacy Start usage)
		r.pty.Close()
	}

	if r.cmd != nil && r.cmd.Process != nil {
		if err := r.cmd.Process.Kill(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
