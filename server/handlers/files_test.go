package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestListFilesHandler(t *testing.T) {
	// Create a temporary workspace
	tmpDir, err := os.MkdirTemp("", "workspace")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a file
	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("content"), 0644)

	// Override workspaceDir for test
	oldWorkspaceDir := workspaceDir
	workspaceDir = tmpDir
	defer func() { workspaceDir = oldWorkspaceDir }()

	req, err := http.NewRequest("GET", "/api/conversations/123/list-files", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ListFilesHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var files []string
	if err := json.NewDecoder(rr.Body).Decode(&files); err != nil {
		t.Errorf("handler returned invalid JSON: %v", err)
	}

	found := false
	for _, f := range files {
		if f == "test.txt" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("test.txt not found in file list: %v", files)
	}
}

func TestSelectFileHandler(t *testing.T) {
	// Create a temporary workspace
	tmpDir, err := os.MkdirTemp("", "workspace")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a file
	content := "file content"
	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte(content), 0644)

	// Override workspaceDir for test
	oldWorkspaceDir := workspaceDir
	workspaceDir = tmpDir
	defer func() { workspaceDir = oldWorkspaceDir }()

	req, err := http.NewRequest("GET", "/api/conversations/123/select-file?file=test.txt", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SelectFileHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Errorf("handler returned invalid JSON: %v", err)
	}

	if resp["code"] != content {
		t.Errorf("handler returned unexpected content: got %v want %v",
			resp["code"], content)
	}
}

func TestPathTraversal(t *testing.T) {
	// Create a temporary workspace
	tmpDir, err := os.MkdirTemp("", "workspace")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Override workspaceDir for test
	oldWorkspaceDir := workspaceDir
	workspaceDir = tmpDir
	defer func() { workspaceDir = oldWorkspaceDir }()

	// Try to access parent directory
	req, err := http.NewRequest("GET", "/api/conversations/123/select-file?file=../test.txt", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SelectFileHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler should block path traversal: got %v want %v",
			status, http.StatusBadRequest)
	}
}
