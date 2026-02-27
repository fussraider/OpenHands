package agent

import "context"

// Delegator defines the interface for delegating tasks to other agents.
type Delegator interface {
	// Delegate starts a sub-agent to handle the task defined by inputs.
	// Returns the outputs from the sub-agent's AgentFinishAction.
	Delegate(ctx context.Context, agentName string, inputs map[string]interface{}) (outputs map[string]interface{}, err error)
}
