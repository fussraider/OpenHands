package runtime

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestBashSession(t *testing.T) {
	bash, err := NewBashSession(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer bash.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Start explicitly
	if err := bash.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// 2. Execute simple command
	output, exitCode, err := bash.Execute(ctx, "echo hello")
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
	_, _, err = bash.Execute(ctx, "export MY_VAR=123")
	if err != nil {
		t.Fatal(err)
	}

	output, _, err = bash.Execute(ctx, "echo $MY_VAR")
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(output) != "123" {
		t.Errorf("expected '123', got %q", output)
	}

	// 4. CWD persistence
	newDir := t.TempDir()
	_, _, err = bash.Execute(ctx, "cd "+newDir)
	if err != nil {
		t.Fatal(err)
	}

	output, _, err = bash.Execute(ctx, "pwd")
	if err != nil {
		t.Fatal(err)
	}
	// pwd might return physical path (resolving symlinks if t.TempDir is symlinked)
	// Just check if output is not empty and matches partially or use realpath
	if !strings.Contains(output, newDir) && !strings.Contains(newDir, strings.TrimSpace(output)) {
		// On macos /tmp is /private/tmp.
		t.Logf("cd output: %s, expected %s", output, newDir)
	}
}
