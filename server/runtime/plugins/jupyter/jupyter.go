package jupyter

import (
	"context"
	"encoding/json"
	"fmt"
	"openhands-go/server/runtime"
	"strings"
	"time"

	"github.com/tmc/langchaingo/llms"
)

// JupyterPlugin implements a persistent Python REPL simulation
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

	// Wrapper code to persist state using pickle
	// We load state at start, execute code, then save state.
	// This simulates a persistent session without a background server.

	wrappedCode := fmt.Sprintf(`
import pickle
import os
import sys
import traceback

STATE_FILE = "/tmp/openhands_python_state.pkl"

# Helper to load state
global_scope = {}
if os.path.exists(STATE_FILE):
    try:
        with open(STATE_FILE, 'rb') as f:
            global_scope = pickle.load(f)
    except Exception:
        pass

# Execute User Code
try:
    exec("""%s""", global_scope)
except Exception:
    traceback.print_exc()

# Helper to save state (filter out unpicklable)
to_save = {}
for k, v in global_scope.items():
    if k.startswith('__') or type(v).__name__ == 'module':
        continue
    try:
        pickle.dumps(v)
        to_save[k] = v
    except:
        pass

with open(STATE_FILE, 'wb') as f:
    pickle.dump(to_save, f)
`, strings.ReplaceAll(params.Code, `"""`, `\"\"\"`)) // Basic escaping

	// We can't write file easily via Execute (echo might fail on complex chars).
	// But we can try to use python -c with the whole blob if it fits argument limits.
	// Alternatively, rely on stateless execution if code is too complex.

	// For MVP reliability, let's stick to the stateless version if stateful wrapper is risky.
	// BUT the user asked for "Refine Jupyter Kernel (Stateful)".

	// Let's use a simpler approach:
	// Instead of wrapping every call, just acknowledge we are stateless for now but structure it clearly.
	// OR: Try to implement the file writing via `cat <<EOF`.

	// Let's try writing the file.
	fileName := fmt.Sprintf("/tmp/exec_%d.py", time.Now().UnixNano())

	// Escape backticks for shell heredoc
	safeContent := strings.ReplaceAll(wrappedCode, "`", "\\`")

	writeCmd := fmt.Sprintf("cat <<'EOF' > %s\n%s\nEOF", fileName, safeContent)
	_, _, err := p.runtime.Execute(ctx, "bash", "-c", writeCmd)
	if err != nil {
		// Fallback to stateless direct execution
		return p.executeStateless(ctx, params.Code)
	}

	output, exitCode, err := p.runtime.Execute(ctx, "python3", fileName)

	// Cleanup
	p.runtime.Execute(ctx, "rm", fileName)

	if err != nil {
		return fmt.Sprintf("Error executing python: %v", err), true, nil
	}

	// Format output
	result := output
	if exitCode != 0 {
		result = fmt.Sprintf("Exit Code: %d\nOutput:\n%s", exitCode, output)
	}

	if strings.TrimSpace(result) == "" {
		result = "(No output)"
	}

	return result, true, nil
}

func (p *JupyterPlugin) executeStateless(ctx context.Context, code string) (string, bool, error) {
	output, exitCode, err := p.runtime.Execute(ctx, "python3", "-c", code)
	if err != nil {
		return fmt.Sprintf("Error executing python: %v", err), true, nil
	}
	result := output
	if exitCode != 0 {
		result = fmt.Sprintf("Exit Code: %d\nOutput:\n%s", exitCode, output)
	}
	return result, true, nil
}
