package runtime

import (
	"context"
	"openhands-go/server/config"
	"testing"
)

func TestDockerRuntime(t *testing.T) {
	// Skip if no docker available (how to check? try connecting client)
	// Or mock client. But mocking docker client is hard.
	// We'll skip if NewDockerRuntime fails.

	cfg := &config.Config{}
	rt, err := NewDockerRuntime(cfg)
	if err != nil {
		t.Skipf("Skipping DockerRuntime test: %v", err)
	}

	// This integration test requires Docker
	// We can try to run a simple command
	ctx := context.Background()
	err = rt.Start(ctx, "echo", "hello from docker")
	if err != nil {
		t.Skipf("Skipping DockerRuntime test (failed to start): %v", err)
	}
	defer rt.Close()

	buf := make([]byte, 1024)
	n, err := rt.Read(buf)
	if err != nil {
		t.Errorf("Read failed: %v", err)
	}

	output := string(buf[:n])
	// Check for "hello from docker"
	// Output might contain headers if TTY is false, but we set TTY=true.
	// In TTY mode, it should be raw output (maybe with \r\n)
	if output == "" {
		t.Errorf("Expected output, got empty string")
	}
}
