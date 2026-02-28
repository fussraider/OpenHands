package security

import (
	"context"
	"openhands-go/server/models"
	"strings"
)

type SecurityRisk int

const (
	RiskUnknown SecurityRisk = -1
	RiskLow     SecurityRisk = 0
	RiskMedium  SecurityRisk = 1
	RiskHigh    SecurityRisk = 2
)

type SecurityAnalyzer interface {
	Analyze(ctx context.Context, action interface{}) (SecurityRisk, string, error)
}

// BasicAnalyzer implements a simple keyword-based policy
func ListAnalyzers() []string {
	return []string{
		"basic",
		"llm",
	}
}

type BasicAnalyzer struct{}

func NewBasicAnalyzer() *BasicAnalyzer {
	return &BasicAnalyzer{}
}

func (a *BasicAnalyzer) Analyze(ctx context.Context, action interface{}) (SecurityRisk, string, error) {
	// Inspect action type
	switch act := action.(type) {
	case models.CmdRunAction:
		return a.analyzeCommand(act.Command)
	default:
		return RiskLow, "", nil
	}
}

func (a *BasicAnalyzer) analyzeCommand(cmd string) (SecurityRisk, string, error) {
	cmd = strings.ToLower(cmd)

	// High risk patterns
	if strings.Contains(cmd, "rm -rf /") || strings.Contains(cmd, "rm -rf /*") {
		return RiskHigh, "Destructive command detected (rm -rf /)", nil
	}
	if strings.Contains(cmd, "mkfs") || strings.Contains(cmd, "dd if=") {
		return RiskHigh, "Filesystem modification detected", nil
	}
	if strings.Contains(cmd, ":(){ :|:& };:") {
		return RiskHigh, "Fork bomb detected", nil
	}

	// Medium risk
	if strings.Contains(cmd, "curl") || strings.Contains(cmd, "wget") {
		return RiskMedium, "Network download detected", nil
	}
	if strings.Contains(cmd, "ssh ") {
		return RiskMedium, "SSH connection detected", nil
	}

	return RiskLow, "", nil
}
