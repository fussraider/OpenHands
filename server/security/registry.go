package security

import (
	"sync"
)

var (
	analyzerRegistryMu sync.RWMutex
	analyzerRegistry   = make(map[string]bool)
)

// RegisterAnalyzer allows a security analyzer implementation to register its name at init time.
func RegisterAnalyzer(name string) {
	analyzerRegistryMu.Lock()
	defer analyzerRegistryMu.Unlock()
	analyzerRegistry[name] = true
}

// GetAvailableAnalyzers returns a list of all dynamically registered analyzer names.
func GetAvailableAnalyzers() []string {
	analyzerRegistryMu.RLock()
	defer analyzerRegistryMu.RUnlock()

	analyzers := make([]string, 0, len(analyzerRegistry))
	for name := range analyzerRegistry {
		analyzers = append(analyzers, name)
	}
	return analyzers
}
