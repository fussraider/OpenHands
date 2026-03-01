package security

import (
	"context"
	"fmt"
	"log/slog"
	"openhands-go/server/llm"
	"openhands-go/server/models"
	"strings"

)

type LLMSecurityAnalyzer struct {
	llmService *llm.LLMService
}

func NewLLMSecurityAnalyzer(service *llm.LLMService) *LLMSecurityAnalyzer {
	return &LLMSecurityAnalyzer{
		llmService: service,
	}
}

func (a *LLMSecurityAnalyzer) Analyze(ctx context.Context, action interface{}) (SecurityRisk, string, error) {
	cmdAction, ok := action.(models.CmdRunAction)
	if !ok {
		return RiskLow, "", nil
	}

	prompt := fmt.Sprintf(`You are a security analyzer. Evaluate the following bash command for security risks.
Command: %s

Risk Levels:
- LOW: Safe, read-only, or standard operations (ls, cat, echo, grep).
- MEDIUM: Network access, file modification in workspace (curl, wget, touch, rm file).
- HIGH: Destructive system-wide changes, infinite loops, fork bombs, accessing sensitive files outside workspace (rm -rf /, mkfs, :(){ :|:& };:).

Respond in the following format:
RISK: [LOW/MEDIUM/HIGH]
REASON: [Short explanation]
`, cmdAction.Command)

	msgs := []llm.Message{
		{Role: "user", Content: prompt},
	}

	// We use a separate context/span for security check?
	resp, err := a.llmService.Complete(ctx, msgs)
	if err != nil {
		slog.Error("Security analysis failed", "error", err)
		return RiskUnknown, "Analysis failed", err
	}

	return a.parseResponse(resp)
}

func (a *LLMSecurityAnalyzer) parseResponse(content string) (SecurityRisk, string, error) {
	lines := strings.Split(content, "\n")
	risk := RiskUnknown
	reason := ""

	for _, line := range lines {
		if strings.HasPrefix(strings.ToUpper(line), "RISK:") {
			level := strings.TrimSpace(strings.TrimPrefix(strings.ToUpper(line), "RISK:"))
			switch level {
			case "LOW":
				risk = RiskLow
			case "MEDIUM":
				risk = RiskMedium
			case "HIGH":
				risk = RiskHigh
			}
		}
		if strings.HasPrefix(strings.ToUpper(line), "REASON:") {
			// Trim prefix case-insensitively by length
			reason = strings.TrimSpace(line[7:])
		}
	}

	if risk == RiskUnknown {
		// Fallback parsing if format is loose
		upper := strings.ToUpper(content)
		if strings.Contains(upper, "RISK: HIGH") || strings.Contains(upper, "HIGH RISK") {
			risk = RiskHigh
		} else if strings.Contains(upper, "RISK: MEDIUM") || strings.Contains(upper, "MEDIUM RISK") {
			risk = RiskMedium
		} else if strings.Contains(upper, "RISK: LOW") || strings.Contains(upper, "LOW RISK") {
			risk = RiskLow
		}
	}

	return risk, reason, nil
}

func init() {
	RegisterAnalyzer("llm")
}
