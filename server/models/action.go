package models

// ActionType defines the type of action an agent can take.
type ActionType string

const (
	ActionTypeCmdRun      ActionType = "run"
	ActionTypeIPythonRun  ActionType = "run_ipython"
	ActionTypeAgentFinish ActionType = "finish"
	ActionTypeMessage     ActionType = "message"
	ActionTypeThink       ActionType = "think"
	ActionTypeDelegate    ActionType = "delegate"
	ActionTypeFileEdit    ActionType = "edit"
	ActionTypeFileRead    ActionType = "read"
)

// BaseAction contains common fields for all actions.
type BaseAction struct {
	Action ActionType `json:"action"`
	Args   any        `json:"args,omitempty"`
}

// CmdRunAction represents a command execution request.
type CmdRunAction struct {
	Action     ActionType `json:"action"` // "run"
	Command    string     `json:"command"`
	Thought    string     `json:"thought,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
}

// IPythonRunCellAction represents an IPython code execution request.
type IPythonRunCellAction struct {
	Action  ActionType `json:"action"` // "run_ipython"
	Code    string     `json:"code"`
	Thought string     `json:"thought,omitempty"`
}

// AgentFinishAction represents the agent completing its task.
type AgentFinishAction struct {
	Action  ActionType        `json:"action"` // "finish"
	Outputs map[string]string `json:"outputs,omitempty"`
	Thought string            `json:"thought,omitempty"`
}

// MessageAction represents a message to the user.
type MessageAction struct {
	Action  ActionType `json:"action"` // "message"
	Content string     `json:"content"`
	Thought string     `json:"thought,omitempty"`
}

// ThinkAction represents the agent's internal thought process.
type ThinkAction struct {
	Action  ActionType `json:"action"` // "think"
	Thought string     `json:"thought"`
}

// AgentDelegateAction represents delegating a task to another agent.
type AgentDelegateAction struct {
	Action     ActionType             `json:"action"` // "delegate"
	Agent      string                 `json:"agent"`
	Inputs     map[string]interface{} `json:"inputs,omitempty"`
	Thought    string                 `json:"thought,omitempty"`
	ToolCallID string                 `json:"tool_call_id,omitempty"`
}

// FileEditAction represents editing a file.
type FileEditAction struct {
	Action      ActionType `json:"action"` // "edit"
	Path        string     `json:"path"`
	Command     string     `json:"command,omitempty"`
	FileText    string     `json:"file_text,omitempty"`
	ViewRange   []int      `json:"view_range,omitempty"`
	OldStr      string     `json:"old_str,omitempty"`
	NewStr      string     `json:"new_str,omitempty"`
	InsertLine  *int       `json:"insert_line,omitempty"`
	Content     string     `json:"content,omitempty"`
	Start       *int       `json:"start,omitempty"`
	End         *int       `json:"end,omitempty"`
	Thought     string     `json:"thought,omitempty"`
	ImplSource  string     `json:"impl_source,omitempty"`
	ToolCallID  string     `json:"tool_call_id,omitempty"`
}

// FileReadAction represents reading a file.
type FileReadAction struct {
	Action     ActionType `json:"action"` // "read"
	Path       string     `json:"path"`
	Thought    string     `json:"thought,omitempty"`
	ImplSource string     `json:"impl_source,omitempty"`
	ViewRange  []int      `json:"view_range,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
}
