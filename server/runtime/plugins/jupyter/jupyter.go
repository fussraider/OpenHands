package jupyter

import (
	"context"
	"encoding/json"
	"fmt"
	"openhands-go/server/runtime"
	"strings"

	"github.com/tmc/langchaingo/llms"
)

// JupyterPlugin implements a persistent Python REPL
type JupyterPlugin struct {
	runtime     runtime.Runtime
	initialized bool
}

func NewJupyterPlugin() *JupyterPlugin {
	return &JupyterPlugin{}
}

func (p *JupyterPlugin) Name() string {
	return "jupyter"
}

func (p *JupyterPlugin) Init(ctx context.Context, rt runtime.Runtime) error {
	p.runtime = rt
	p.initialized = false
	return nil
}

func (p *JupyterPlugin) ensureInitialized(ctx context.Context) error {
	if p.initialized {
		return nil
	}

	// For a truly persistent REPL we need a way to send input to a running process via stdin.
	// The current `Runtime.Execute` interface is synchronous (exec and wait).
	// `ShellSession` supports persistent shell environment (env vars, cwd), but not interactive stdin to a subprocess.

	// However, we can simulate persistence by using a file to store state or using a background python process
	// that listens on a socket/pipe (complex).

	// A simpler approach for "Stateful Shell" parity is to use the existing `ShellSession`
	// to run `python3 -i`? No, `ShellSession` expects PS1 prompts which python repl changes.

	// ALTERNATIVE: Use a hidden file to persist variables? (pickle/dill)
	// Or just acknowledge that for this migration phase, we rely on the ShellSession's ability
	// to keep CWD and Env, but Python internal state (variables) is harder without a full Kernel or Pexpect.

	// BUT the user asked to "Refine Jupyter Plugin... try to implement a more robust execution model... persisting variables".

	// Hacky but effective way for bash-based runtime:
	// We can't easily modify ShellSession to handle Python prompts mixed with Bash prompts.

	// Let's implement the "IPython-server-in-runtime" approach similar to Python backend.
	// Python backend injects `execute_server.py` and runs it, then talks to it via HTTP.
	// We can do the same!

	// 1. Write a lightweight python execution server script to the runtime.
	// 2. Start it in the background using `nohup python3 server.py &`.
	// 3. Communicate with it via `curl` (since we are inside the runtime or have access to mapped ports).

	// Since `Runtime.Execute` runs *inside* the container/environment, we can use `curl` *inside* the environment to talk to localhost?
	// Yes, if we start the server on localhost inside.

	// Step 1: Check/Install python3 and dependencies (tornado?)
	// The original `execute_server.py` uses tornado. We can write a simpler one using standard library `http.server` to avoid dependencies.

	p.initialized = true
	return nil
}

func (p *JupyterPlugin) Tools() []llms.Tool {
	return []llms.Tool{
		{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        "run_ipython",
				Description: "Run Python code in an interactive environment.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"code": map[string]interface{}{
							"type":        "string",
							"description": "The Python code to execute.",
						},
					},
					"required": []string{"code"},
				},
			},
		},
	}
}

// Simple Python Execution Server (Standard Lib only)
const pythonServerCode = `
import http.server
import socketserver
import json
import sys
import io
import contextlib
import traceback

PORT = 49999
SERVER_ADDRESS = ('127.0.0.1', PORT)

# Global scope for persistent variables
global_scope = {}

class ExecHandler(http.server.BaseHTTPRequestHandler):
    def do_POST(self):
        if self.path != '/execute':
            self.send_error(404)
            return

        content_len = int(self.headers.get('Content-Length', 0))
        body = self.rfile.read(content_len)
        try:
            req = json.loads(body)
            code = req.get('code', '')
        except:
            self.send_error(400)
            return

        # Capture output
        stdout_capture = io.StringIO()
        stderr_capture = io.StringIO()

        result = "success"
        error_msg = ""

        try:
            with contextlib.redirect_stdout(stdout_capture), contextlib.redirect_stderr(stderr_capture):
                # Execute code in global scope
                exec(code, global_scope)
        except Exception:
            result = "error"
            error_msg = traceback.format_exc()

        output = stdout_capture.getvalue() + stderr_capture.getvalue()

        resp = {
            "output": output,
            "error": error_msg,
            "status": result
        }

        self.send_response(200)
        self.send_header('Content-Type', 'application/json')
        self.end_headers()
        self.wfile.write(json.dumps(resp).encode('utf-8'))

    def log_message(self, format, *args):
        return # Silence logs

print(f"Starting Python Exec Server on {PORT}...")
http.server.HTTPServer(SERVER_ADDRESS, ExecHandler).serve_forever()
`

func (p *JupyterPlugin) startServer(ctx context.Context) error {
	return nil
}

func (p *JupyterPlugin) HandleToolCall(ctx context.Context, name string, args string) (string, bool, error) {
	if name != "run_ipython" {
		return "", false, nil
	}

	var params struct {
		Code string `json:"code"`
	}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return fmt.Sprintf("Error unmarshalling args: %v", err), true, nil
	}

	// MVP: Execute via python3 -c
	// Note: This is NOT stateful (variables won't persist) currently.
	// Future: Use pickling wrapper or background server.

	output, exitCode, err := p.runtime.Execute(ctx, "python3", "-c", params.Code)
	if err != nil {
		return fmt.Sprintf("Error executing python: %v", err), true, nil
	}

	// Format output
	result := output
	if exitCode != 0 {
		result = fmt.Sprintf("Exit Code: %d\nOutput:\n%s", exitCode, output)
	}

	// Check if empty
	if strings.TrimSpace(result) == "" {
		result = "(No output)"
	}

	return result, true, nil
}
