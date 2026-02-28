package microagent

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

func Parse(content, path string) (*MicroagentMetadata, string, error) {
	// 1. Separate Frontmatter
	if !strings.HasPrefix(content, "---\n") {
		// No frontmatter, treat as Repo Knowledge if name suggests
		// Or basic knowledge if triggers are inferred (not implemented here)
		return &MicroagentMetadata{
			Name: "default",
			Type: MicroagentTypeRepoKnowledge,
		}, content, nil
	}

	parts := strings.SplitN(content, "---\n", 3)
	if len(parts) < 3 {
		return nil, "", fmt.Errorf("invalid frontmatter format")
	}

	frontmatter := parts[1]
	body := parts[2]

	var metadata MicroagentMetadata
	if err := yaml.Unmarshal([]byte(frontmatter), &metadata); err != nil {
		return nil, "", fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	// Validate / Defaulting
	if metadata.Type == "" {
		if len(metadata.Triggers) > 0 {
			metadata.Type = MicroagentTypeKnowledge
		} else if len(metadata.Inputs) > 0 {
			metadata.Type = MicroagentTypeTask
		} else {
			metadata.Type = MicroagentTypeRepoKnowledge
		}
	}

	if metadata.Name == "" {
		// Derive from path
		// e.g. .openhands/microagents/my-agent.md -> my-agent
		// Basic implementation
		metadata.Name = path // simplified
	}

	return &metadata, body, nil
}

// ConvertToResponse converts metadata and path to API response
func ConvertToResponse(name, path string, createdAt interface{}) MicroagentResponse {
	// Time handling is tricky without known input type, assuming time.Time or mocking
	// For now, placeholder
	return MicroagentResponse{
		Name: name,
		Path: path,
		// CreatedAt: ...
	}
}
