package runtime

import "context"

type Runtime interface {
	Start(ctx context.Context, command string, args ...string) error
	// Execute runs a command and returns the output and exit code.
	// This supports stateful execution if the runtime supports it.
	Execute(ctx context.Context, command string, args ...string) (string, int, error)

	Write(p []byte) (n int, err error)
	Read(p []byte) (n int, err error)
	Close() error

	// GetCwd returns the current working directory of the shell session
	GetCwd(ctx context.Context) (string, error)

	// CopyFileToContainer copies a local file to the runtime container
	CopyFileToContainer(ctx context.Context, hostPath string, containerPath string) error

	// CopyFileFromContainer copies a file from the runtime container to the local host
	CopyFileFromContainer(ctx context.Context, containerPath string, hostPath string) error

	// Future methods for parity with Python
	// RunAction(action Action) (Observation, error)
	// ReadFile(path string) ([]byte, error)
	// WriteFile(path string, content []byte) error
}
