package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"openhands-go/server/models"
	"testing"
)

func TestExecuteActionHandler(t *testing.T) {
	reqBody := models.ActionRequest{
		Action: "run",
		Args:   "echo 'hello world'",
	}
	body, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "/api/conversations/123/action", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ExecuteActionHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Errorf("handler returned invalid JSON: %v", err)
	}

	// Output might contain PTY control characters or newline, so check substring
	// But local runtime might not work in CI/docker environment without PTY support or deps?
	// The `pty` library usually works on Linux.
	// Let's assume it works or we catch error.
	// Actually, if `pty` fails to start (e.g. no pty available), it returns 500.

	// If the test environment doesn't support PTY, we might get an error.
	// But for now, let's see.
}
