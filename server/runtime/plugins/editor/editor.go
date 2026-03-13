package editor

import (
	"context"
	"encoding/json"
	"fmt"
	"openhands-go/server/runtime"
	"strings"

	"github.com/tmc/langchaingo/llms"
)

type EditorPlugin struct {
	runtime runtime.Runtime
}

func NewEditorPlugin() *EditorPlugin {
	return &EditorPlugin{}
}

func (p *EditorPlugin) Name() string {
	return "str_replace_editor"
}

func (p *EditorPlugin) Init(ctx context.Context, rt runtime.Runtime) error {
	p.runtime = rt
	return nil
}

func (p *EditorPlugin) Tools() []llms.Tool {
	return []llms.Tool{
		{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        "str_replace_editor",
				Description: "Custom editing tool for viewing, creating and editing files in plain-text format. Allowed commands: view, create, str_replace, insert, undo_edit.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"command": map[string]interface{}{
							"type":        "string",
							"description": "The commands to run. Allowed options are: view, create, str_replace, insert, undo_edit.",
						},
						"path": map[string]interface{}{
							"type":        "string",
							"description": "Absolute path to file or directory, e.g. /workspace/file.py.",
						},
						"file_text": map[string]interface{}{
							"type":        "string",
							"description": "Required parameter of create command, with the content of the file to be created.",
						},
						"old_str": map[string]interface{}{
							"type":        "string",
							"description": "Required parameter of str_replace command containing the string in path to replace.",
						},
						"new_str": map[string]interface{}{
							"type":        "string",
							"description": "Optional parameter of str_replace command containing the new string. Required parameter of insert command.",
						},
						"insert_line": map[string]interface{}{
							"type":        "integer",
							"description": "Required parameter of insert command. The new_str will be inserted AFTER the line insert_line of path.",
						},
						"view_range": map[string]interface{}{
							"type":        "array",
							"items":       map[string]interface{}{"type": "integer"},
							"description": "Optional parameter of view command when path points to a file. Defaults to full file. E.g. [11, 12] shows lines 11 and 12.",
						},
					},
					"required": []string{"command", "path"},
				},
			},
		},
	}
}

func (p *EditorPlugin) HandleToolCall(ctx context.Context, name string, args string) (string, bool, error) {
	if name != "str_replace_editor" {
		return "", false, nil
	}

	var params struct {
		Command    string `json:"command"`
		Path       string `json:"path"`
		FileText   string `json:"file_text"`
		OldStr     string `json:"old_str"`
		NewStr     string `json:"new_str"`
		InsertLine *int   `json:"insert_line"`
	}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return fmt.Sprintf("Error unmarshalling args: %v", err), true, nil
	}

	switch params.Command {
	case "view":
		out, _, err := p.runtime.Execute(ctx, "cat", "-n", params.Path)
		if err != nil {
			return fmt.Sprintf("File not found or unreadable: %s\n%s", params.Path, out), true, nil
		}
		return out, true, nil

	case "create":
		escaped := strings.ReplaceAll(params.FileText, "`", "\\`")
		out, _, err := p.runtime.Execute(ctx, "bash", "-c", fmt.Sprintf("cat <<'EOF' > %s\n%s\nEOF", params.Path, escaped))
		if err != nil {
			return fmt.Sprintf("Error creating file: %s\n%s", err.Error(), out), true, nil
		}
		return fmt.Sprintf("File created successfully at: %s", params.Path), true, nil

	case "str_replace":
		// Read file
		content, _, err := p.runtime.Execute(ctx, "cat", params.Path)
		if err != nil {
			return fmt.Sprintf("File not found: %s", params.Path), true, nil
		}

		// Check uniqueness
		count := strings.Count(content, params.OldStr)
		if count == 0 {
			return "Error: old_str not found in file. Make sure it matches exactly.", true, nil
		}
		if count > 1 {
			return "Error: old_str found multiple times. Please provide more context to make it unique.", true, nil
		}

		// Replace
		newContent := strings.Replace(content, params.OldStr, params.NewStr, 1)

		// Write back
		escaped := strings.ReplaceAll(newContent, "`", "\\`")
		out, _, err := p.runtime.Execute(ctx, "bash", "-c", fmt.Sprintf("cat <<'EOF' > %s\n%s\nEOF", params.Path, escaped))
		if err != nil {
			return fmt.Sprintf("Error writing file: %s\n%s", err.Error(), out), true, nil
		}
		return fmt.Sprintf("Successfully replaced string in %s", params.Path), true, nil

	case "insert":
		return "Insert command not currently implemented in simplified backend plugin.", true, nil

	case "undo_edit":
		return "Undo command not currently implemented in simplified backend plugin.", true, nil

	default:
		return fmt.Sprintf("Unknown command: %s", params.Command), true, nil
	}
}
