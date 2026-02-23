package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "OK"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestModelsHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/options/models", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(modelsHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check content type
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v",
			contentType, "application/json")
	}
}

func TestSPAHandler(t *testing.T) {
	// Create a temporary directory for static files
	tmpDir, err := os.MkdirTemp("", "static")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create index.html
	indexContent := "<html>index</html>"
	err = os.WriteFile(filepath.Join(tmpDir, "index.html"), []byte(indexContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Create another file
	fileContent := "some content"
	err = os.WriteFile(filepath.Join(tmpDir, "file.txt"), []byte(fileContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	fs := http.FileServer(http.Dir(tmpDir))
	handler := spaHandler(tmpDir, fs)

	// Test 1: Request existing file
	req, _ := http.NewRequest("GET", "/file.txt", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Test 1: expected 200, got %d", rr.Code)
	}
	if rr.Body.String() != fileContent {
		t.Errorf("Test 1: expected %q, got %q", fileContent, rr.Body.String())
	}

	// Test 2: Request non-existent file (should serve index.html)
	req, _ = http.NewRequest("GET", "/unknown", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Test 2: expected 200, got %d", rr.Code)
	}
	if rr.Body.String() != indexContent {
		t.Errorf("Test 2: expected %q, got %q", indexContent, rr.Body.String())
	}

	// Test 3: Request API path (should 404)
	req, _ = http.NewRequest("GET", "/api/unknown", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Errorf("Test 3: expected 404, got %d", rr.Code)
	}
}
