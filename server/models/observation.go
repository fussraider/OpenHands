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
}
