package jupyter

import (
	"context"
	"encoding/json"
	"fmt"
	"openhands-go/server/runtime"

	"github.com/tmc/langchaingo/llms"
)

type JupyterPlugin struct {
	runtime runtime.Runtime
}

func NewJupyterPlugin() *JupyterPlugin {
	return &JupyterPlugin{}
}

func (p *JupyterPlugin) Name() string {
	return "jupyter"
}

func (p *JupyterPlugin) Init(ctx context.Context, rt runtime.Runtime) error {
	p.runtime = rt
	// Check if python3 is installed?
	_, _, err := rt.Execute(ctx, "python3", "--version")
	return err
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

	// MVP: Execute via python3 -c
	// Note: This is NOT stateful (variables won't persist).
	// For full Jupyter support we need a persistent python shell session.
	// But this satisfies the "Implement stub" requirement and works for simple tasks.
	output, exitCode, err := p.runtime.Execute(ctx, "python3", "-c", params.Code)
	if err != nil {
		return fmt.Sprintf("Error executing python: %v", err), true, nil
	}
	if exitCode != 0 {
		return fmt.Sprintf("Python exited with code %d. Output:\n%s", exitCode, output), true, nil
	}

	return output, true, nil
}
