package microagent

import (
	"time"
)

type MicroagentType string

const (
	MicroagentTypeKnowledge     MicroagentType = "knowledge"
	MicroagentTypeRepoKnowledge MicroagentType = "repo"
	MicroagentTypeTask          MicroagentType = "task"
)

type InputMetadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type MicroagentMetadata struct {
	Name     string          `json:"name"`
	Type     MicroagentType  `json:"type"`
	Version  string          `json:"version"`
	Agent    string          `json:"agent"`
	Triggers []string        `json:"triggers,omitempty"`
	Inputs   []InputMetadata `json:"inputs,omitempty"`
	// MCPTools not supported yet in Go struct
}

type MicroagentResponse struct {
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"created_at"`
}

type MicroagentContentResponse struct {
	Content     string   `json:"content"`
	Path        string   `json:"path"`
	Triggers    []string `json:"triggers"`
	GitProvider string   `json:"git_provider,omitempty"`
}
