package agent

import (
	"sync"
)

var (
	agentRegistryMu sync.RWMutex
	agentRegistry   = make(map[string]bool)
)

// RegisterAgent allows an agent implementation to register its name at init time.
// This mimics the Python subclass introspection registry pattern.
func RegisterAgent(name string) {
	agentRegistryMu.Lock()
	defer agentRegistryMu.Unlock()
	agentRegistry[name] = true
}

// GetAvailableAgents returns a list of all dynamically registered agent names.
func GetAvailableAgents() []string {
	agentRegistryMu.RLock()
	defer agentRegistryMu.RUnlock()

	agents := make([]string, 0, len(agentRegistry))
	for name := range agentRegistry {
		agents = append(agents, name)
	}
	return agents
}
