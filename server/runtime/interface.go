package runtime

import "context"

type Runtime interface {
	Start(ctx context.Context, command string, args ...string) error
	Write(p []byte) (n int, err error)
	Read(p []byte) (n int, err error)
	Close() error
}
