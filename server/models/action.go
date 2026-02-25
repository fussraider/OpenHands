package models

// ActionType defines the type of action an agent can take.
type ActionType string

const (
	ActionTypeCmdRun       ActionType = "run"
	ActionTypeIPythonRun   ActionType = "run_ipython"
	ActionTypeAgentFinish  ActionType = "finish"
	ActionTypeMessage      ActionType = "message"
	ActionTypeThink        ActionType = "think"
	ActionTypeDelegate     ActionType = "delegate"
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
	Action ActionType `json:"action"` // "run_ipython"
	Code   string     `json:"code"`
	Thought string    `json:"thought,omitempty"`
}

// AgentFinishAction represents the agent completing its task.
type AgentFinishAction struct {
	Action  ActionType `json:"action"` // "finish"
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
	Action  ActionType             `json:"action"` // "delegate"
	Agent   string                 `json:"agent"`
	Inputs  map[string]interface{} `json:"inputs,omitempty"`
	Thought string                 `json:"thought,omitempty"`
}
