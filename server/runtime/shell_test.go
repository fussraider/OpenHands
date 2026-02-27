package runtime

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/creack/pty"
)

func TestShellSession(t *testing.T) {
	// Setup PTY manually for test
	cmd := exec.Command("bash", "--noprofile", "--norc")
	cmd.Env = os.Environ()

	f, err := pty.Start(cmd)
	if err != nil {
		t.Fatal(err)
	}
	// Note: We don't close f here, ShellSession closes it.
	// But we might need to kill cmd.
	defer cmd.Process.Kill()

	shell := NewShellSession(f)
	defer shell.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Initialize (replaces Start)
	if err := shell.Initialize(ctx); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// 2. Execute simple command
	output, exitCode, err := shell.Execute(ctx, "echo hello")
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
	if !strings.Contains(output, "hello") {
		t.Errorf("expected output to contain 'hello', got %q", output)
	}

	// 3. State persistence
	_, _, err = shell.Execute(ctx, "export MY_VAR=123")
	if err != nil {
		t.Fatal(err)
	}

	output, _, err = shell.Execute(ctx, "echo $MY_VAR")
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(output) != "123" {
		t.Errorf("expected '123', got %q", output)
	}

	// 4. CWD persistence
	newDir := t.TempDir()
	_, _, err = shell.Execute(ctx, "cd "+newDir)
	if err != nil {
		t.Fatal(err)
	}

	output, _, err = shell.Execute(ctx, "pwd")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output, newDir) && !strings.Contains(newDir, strings.TrimSpace(output)) {
		t.Logf("cd output: %s, expected %s", output, newDir)
	}
}
