package models

type CmdOutputMetadata struct {
	ExitCode      int    `json:"exit_code"`
	PID           int    `json:"pid"`
	Username      string `json:"username,omitempty"`
	Hostname      string `json:"hostname,omitempty"`
	WorkingDir    string `json:"working_dir,omitempty"`
	PyInterpreter string `json:"py_interpreter_path,omitempty"`
	Prefix        string `json:"prefix,omitempty"`
	Suffix        string `json:"suffix,omitempty"`
}

type CmdOutputObservation struct {
	Observation string            `json:"observation"` // "run"
	Content     string            `json:"content"`
	Metadata    CmdOutputMetadata `json:"metadata"`
	Command     string            `json:"command,omitempty"`
	Hidden      bool              `json:"hidden,omitempty"`
	ToolCallID  string            `json:"tool_call_id,omitempty"`
}

type TaskState string

const (
	TaskStateStarted    TaskState = "started"
	TaskStateRunning    TaskState = "running"
	TaskStateCompleted  TaskState = "completed"
	TaskStateFailed     TaskState = "failed"
	TaskStateDelegated  TaskState = "delegated"
)

// TaskTrackingObservation represents a change in task status or list
type TaskTrackingObservation struct {
	Observation string      `json:"observation"` // "task_tracking"
	Content     string      `json:"content"`     // Message describing the update
	TaskList    []TaskItem  `json:"task_list,omitempty"`
}

type TaskItem struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	State       TaskState `json:"state"`
	Children    []TaskItem `json:"children,omitempty"`
}

// LoopDetectionObservation represents detection of an infinite loop
type LoopDetectionObservation struct {
	Observation string `json:"observation"` // "loop_detection"
	Content     string `json:"content"`     // Warning message
}
