package security

import (
	"context"
	"openhands-go/server/models"
	"testing"
)

func TestBasicAnalyzer(t *testing.T) {
	analyzer := NewBasicAnalyzer()
	ctx := context.Background()

	tests := []struct {
		cmd          string
		expectedRisk SecurityRisk
	}{
		{"echo hello", RiskLow},
		{"ls -la", RiskLow},
		{"curl http://example.com", RiskMedium},
		{"wget http://example.com", RiskMedium},
		{"rm -rf /", RiskHigh},
		{"rm -rf /*", RiskHigh},
		{"sudo rm -rf /", RiskHigh}, // contains check
		{":(){ :|:& };:", RiskHigh},
	}

	for _, tt := range tests {
		action := models.CmdRunAction{Command: tt.cmd}
		risk, _, err := analyzer.Analyze(ctx, action)
		if err != nil {
			t.Errorf("Analyze failed: %v", err)
		}
		if risk != tt.expectedRisk {
			t.Errorf("Cmd: %s, Expected %d, Got %d", tt.cmd, tt.expectedRisk, risk)
		}
	}
}
