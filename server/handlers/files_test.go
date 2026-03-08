package handlers

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
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

func TestUploadFilesHandler(t *testing.T) {
	// Setup test workspace
	tmpDir, err := os.MkdirTemp("", "openhands_test_workspace_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	oldWorkspace := workspaceDir
	workspaceDir = tmpDir
	defer func() { workspaceDir = oldWorkspace }()

	// Create a multipart form buffer
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("files", "upload.txt")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte("uploaded content"))
	writer.Close()

	req, err := http.NewRequest("POST", "/api/conversations/123/upload-files", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UploadFilesHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var response map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&response)

	uploaded := response["uploaded_files"].([]interface{})
	if len(uploaded) != 1 || uploaded[0].(string) != "upload.txt" {
		t.Errorf("unexpected uploaded files response: %v", response)
	}

	// Verify file exists on disk
	content, err := os.ReadFile(filepath.Join(tmpDir, "upload.txt"))
	if err != nil {
		t.Errorf("failed to read uploaded file: %v", err)
	}
	if string(content) != "uploaded content" {
		t.Errorf("expected file content 'uploaded content', got '%s'", string(content))
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
