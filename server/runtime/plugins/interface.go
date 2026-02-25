package plugins

import (
	"context"
	"openhands-go/server/runtime"

	"github.com/tmc/langchaingo/llms"
)

// Plugin defines the interface for runtime plugins.
type Plugin interface {
	Name() string
	// Init initializes the plugin with the given runtime.
	Init(ctx context.Context, rt runtime.Runtime) error
	// Tools returns the list of tools provided by the plugin.
	Tools() []llms.Tool
	// HandleToolCall executes a tool call. Returns handled=true if the plugin handles this tool.
	HandleToolCall(ctx context.Context, name string, args string) (output string, handled bool, err error)
}
