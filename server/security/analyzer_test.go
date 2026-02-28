package security

import (
	"context"
	"openhands-go/server/models"
	"testing"
)

// MockLLMService can't be easily mocked as LLMService is a struct, not interface.
// But we can check parseResponse logic at least.

func TestParseResponse(t *testing.T) {
	analyzer := &LLMSecurityAnalyzer{}

	tests := []struct {
		input      string
		wantRisk   SecurityRisk
		wantReason string
	}{
		{"RISK: HIGH\nREASON: rm -rf / is dangerous", RiskHigh, "rm -rf / is dangerous"},
		{"RISK: MEDIUM\nREASON: curl downloads file", RiskMedium, "curl downloads file"},
		{"RISK: LOW\nREASON: ls is safe", RiskLow, "ls is safe"},
		{"Risk: High\nReason: bad stuff", RiskHigh, "bad stuff"},
		{"Some noise\nRISK: HIGH\nMore noise", RiskHigh, ""}, // Basic parsing might fail reason if not standard
	}

	for _, tt := range tests {
		risk, reason, _ := analyzer.parseResponse(tt.input)
		if risk != tt.wantRisk {
			t.Errorf("parseResponse(%q) risk = %v, want %v", tt.input, risk, tt.wantRisk)
		}
		if tt.wantReason != "" && reason != tt.wantReason {
			t.Errorf("parseResponse(%q) reason = %q, want %q", tt.input, reason, tt.wantReason)
		}
	}
}

func TestBasicAnalyzer(t *testing.T) {
	a := NewBasicAnalyzer()

	tests := []struct {
		cmd      string
		wantRisk SecurityRisk
	}{
		{"ls -la", RiskLow},
		{"rm -rf /", RiskHigh},
		{"curl http://malware.com", RiskMedium},
		{"ssh user@host", RiskMedium},
	}

	for _, tt := range tests {
		act := models.CmdRunAction{Command: tt.cmd}
		risk, _, _ := a.Analyze(context.Background(), act)
		if risk != tt.wantRisk {
			t.Errorf("Analyze(%q) risk = %v, want %v", tt.cmd, risk, tt.wantRisk)
		}
	}
}
