package runtime

import "context"

type Runtime interface {
	Start(ctx context.Context, command string, args ...string) error
	Write(p []byte) (n int, err error)
	Read(p []byte) (n int, err error)
	Close() error

	// Future methods for parity with Python
	// RunAction(action Action) (Observation, error)
	// ReadFile(path string) ([]byte, error)
	// WriteFile(path string, content []byte) error
}
