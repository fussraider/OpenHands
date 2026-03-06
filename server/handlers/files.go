package handlers

import (
	"encoding/json"
	"log/slog"
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

// filesToIgnore contains list of files/directories to exclude from listing
var filesToIgnore = []string{
	".git",
	".DS_Store",
	"node_modules",
	"__pycache__",
	"lost+found",
	".vscode",
}

func shouldIgnore(name string) bool {
	for _, ignore := range filesToIgnore {
		if name == ignore || strings.HasPrefix(name, ignore+"/") {
			return true
		}
	}
	return false
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
		name := file.Name()
		if shouldIgnore(name) {
			continue
		}

		// Python implementation returns relative paths
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

// UploadFilesHandler handles file uploads to the workspace
func UploadFilesHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("File upload config: max_size=10MB, restrict_types=false, allowed_extensions=[.*]")

	// Parse multipart form, 10MB max memory
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	// Ensure workspace dir exists
	os.MkdirAll(workspaceDir, 0755)

	var uploadedFiles []string
	var skippedFiles []map[string]string

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			skippedFiles = append(skippedFiles, map[string]string{
				"name":   fileHeader.Filename,
				"reason": "Failed to open file: " + err.Error(),
			})
			continue
		}

		// Ensure safe path to prevent directory traversal in filename
		safeName := filepath.Base(fileHeader.Filename)
		destPath := filepath.Join(workspaceDir, safeName)

		destFile, err := os.Create(destPath)
		if err != nil {
			skippedFiles = append(skippedFiles, map[string]string{
				"name":   safeName,
				"reason": "Failed to create file: " + err.Error(),
			})
			file.Close()
			continue
		}

		_, err = destFile.ReadFrom(file)
		if err != nil {
			skippedFiles = append(skippedFiles, map[string]string{
				"name":   safeName,
				"reason": "Failed to write file: " + err.Error(),
			})
		} else {
			uploadedFiles = append(uploadedFiles, safeName)
		}

		destFile.Close()
		file.Close()
	}

	// For MVP, if we have a docker runtime, ideally we'd use CopyFileToContainer.
	// But since ListFilesHandler and SelectFileHandler read from workspaceDir directly,
	// and the Docker container mounts the workspaceDir, saving locally to workspaceDir is sufficient.

	if uploadedFiles == nil {
		uploadedFiles = []string{}
	}
	if skippedFiles == nil {
		skippedFiles = []map[string]string{}
	}

	response := map[string]interface{}{
		"uploaded_files": uploadedFiles,
		"skipped_files":  skippedFiles,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SelectFileHandler returns file content
// ZipWorkspaceHandler returns a zip of the workspace (placeholder)
func ZipWorkspaceHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Zipping workspace")
	http.Error(w, "Not implemented in MVP", http.StatusNotImplemented)
}

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
