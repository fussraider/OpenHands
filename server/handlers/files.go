package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var workspaceDir = "workspace"

// isPathSafe checks if the path is within the workspaceDir
func isPathSafe(base, path string) bool {
	absBase, err := filepath.Abs(base)
	if err != nil {
		return false
	}
	absPath, err := filepath.Abs(filepath.Join(base, path))
	if err != nil {
		return false
	}

	// Ensure the base path ends with a separator to prevent partial matches
	// e.g. /workspace matching /workspace_secrets
	if !strings.HasSuffix(absBase, string(os.PathSeparator)) {
		absBase += string(os.PathSeparator)
	}

	return strings.HasPrefix(absPath, absBase) || absPath == strings.TrimSuffix(absBase, string(os.PathSeparator))
}

// ListFilesHandler lists files in the workspace
func ListFilesHandler(w http.ResponseWriter, r *http.Request) {
	// r.PathValue("id") is conversation ID, but we share workspace for now
	path := r.URL.Query().Get("path")
	if path == "" {
		path = "."
	}

	if !isPathSafe(workspaceDir, path) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	fullPath := filepath.Join(workspaceDir, path)

	files, err := os.ReadDir(fullPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fileList := make([]string, 0, len(files))
	for _, file := range files {
		// Only list files, or directories if path is "."?
		// The python implementation returns relative paths
		name := file.Name()
		if path != "." {
			name = filepath.Join(path, name)
		}

		if file.IsDir() {
			name += "/"
		}
		fileList = append(fileList, name)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fileList)
}

// SelectFileHandler returns file content
func SelectFileHandler(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Query().Get("file")
	if file == "" {
		http.Error(w, "file parameter is required", http.StatusBadRequest)
		return
	}

	if !isPathSafe(workspaceDir, file) {
		http.Error(w, "invalid file path", http.StatusBadRequest)
		return
	}

	fullPath := filepath.Join(workspaceDir, file)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "file not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"code": string(content)})
}
