package handlers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"recipe-app/internal/logger"
)

type FileHandler struct {
	uploadDir    string
	maxFileSize  int64
	allowedTypes []string
}

func NewFileHandler() *FileHandler {
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}

	return &FileHandler{
		uploadDir:    uploadDir,
		maxFileSize:  10 * 1024 * 1024, // 10MB
		allowedTypes: []string{"image/jpeg", "image/png", "image/gif", "image/webp"},
	}
}

func (h *FileHandler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	logger.FromContext(r.Context()).Info("Processing file upload")

	// Parse multipart form
	err := r.ParseMultipartForm(h.maxFileSize)
	if err != nil {
		logger.LogError(r.Context(), err, "Failed to parse multipart form")
		http.Error(w, "File too large", http.StatusRequestEntityTooLarge)
		return
	}

	// Get the file from form
	file, header, err := r.FormFile("file")
	if err != nil {
		logger.LogError(r.Context(), err, "Failed to get file from form")
		http.Error(w, "No file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file type
	contentType := header.Header.Get("Content-Type")
	if !h.isAllowedType(contentType) {
		logger.FromContext(r.Context()).Info("Rejected file upload", "type", contentType)
		http.Error(w, "File type not allowed", http.StatusBadRequest)
		return
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	hash := md5.Sum([]byte(fmt.Sprintf("%s%d", header.Filename, time.Now().Unix())))
	filename := fmt.Sprintf("%x%s", hash, ext)

	// Create uploads directory if it doesn't exist
	if err := os.MkdirAll(h.uploadDir, 0755); err != nil {
		logger.LogError(r.Context(), err, "Failed to create upload directory")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create the file
	dst, err := os.Create(filepath.Join(h.uploadDir, filename))
	if err != nil {
		logger.LogError(r.Context(), err, "Failed to create file")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the destination
	_, err = io.Copy(dst, file)
	if err != nil {
		logger.LogError(r.Context(), err, "Failed to save file")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return success response with file info
	fileURL := fmt.Sprintf("/uploads/%s", filename)

	logger.FromContext(r.Context()).Info("File uploaded successfully", "filename", filename, "size", header.Size)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":     "File uploaded successfully",
		"filename":    filename,
		"url":         fileURL,
		"size":        header.Size,
		"contentType": contentType,
	})
}

func (h *FileHandler) HandleMultiUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	logger.FromContext(r.Context()).Info("Processing multiple file upload")

	err := r.ParseMultipartForm(h.maxFileSize * 5) // Allow 5 files max
	if err != nil {
		logger.LogError(r.Context(), err, "Failed to parse multipart form")
		http.Error(w, "Files too large", http.StatusRequestEntityTooLarge)
		return
	}

	// Get files from form
	files := r.MultipartForm.File
	if len(files) == 0 {
		http.Error(w, "No files provided", http.StatusBadRequest)
		return
	}

	var uploadedFiles []map[string]interface{}
	maxFiles := 5
	fileCount := 0

	for fieldName, fileHeaders := range files {
		for _, header := range fileHeaders {
			if fileCount >= maxFiles {
				break
			}

			file, err := header.Open()
			if err != nil {
				logger.LogError(r.Context(), err, "Failed to open file")
				continue
			}

			contentType := header.Header.Get("Content-Type")
			if !h.isAllowedType(contentType) {
				file.Close()
				continue
			}

			// Generate unique filename
			ext := filepath.Ext(header.Filename)
			hash := md5.Sum([]byte(fmt.Sprintf("%s%d%s", fieldName, time.Now().UnixNano(), header.Filename)))
			filename := fmt.Sprintf("%x%s", hash, ext)

			// Save file
			dstPath := filepath.Join(h.uploadDir, filename)
			dst, err := os.Create(dstPath)
			if err != nil {
				logger.LogError(r.Context(), err, "Failed to create file")
				file.Close()
				continue
			}

			_, err = io.Copy(dst, file)
			file.Close()
			dst.Close()

			if err != nil {
				logger.LogError(r.Context(), err, "Failed to save file")
				continue
			}

			uploadedFiles = append(uploadedFiles, map[string]interface{}{
				"fieldName":    fieldName,
				"filename":     filename,
				"originalName": header.Filename,
				"url":          fmt.Sprintf("/uploads/%s", filename),
				"size":         header.Size,
				"contentType":  contentType,
			})
			fileCount++
		}
	}

	logger.FromContext(r.Context()).Info("Files uploaded successfully", "count", len(uploadedFiles))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Files uploaded successfully",
		"files":   uploadedFiles,
		"count":   len(uploadedFiles),
	})
}

func (h *FileHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filename := chi.URLParam(r, "filename")
	if filename == "" {
		http.Error(w, "Filename is required", http.StatusBadRequest)
		return
	}

	// Security: Validate filename to prevent path traversal
	if strings.Contains(filename, "..") || strings.HasPrefix(filename, "/") {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(h.uploadDir, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Delete the file
	err := os.Remove(filePath)
	if err != nil {
		logger.LogError(r.Context(), err, "Failed to delete file")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.FromContext(r.Context()).Info("File deleted successfully", "filename", filename)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "File deleted successfully",
	})
}

func (h *FileHandler) ServeFile(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")
	if filename == "" {
		http.Error(w, "Filename is required", http.StatusBadRequest)
		return
	}

	// Security: Validate filename
	if strings.Contains(filename, "..") || strings.HasPrefix(filename, "/") {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(h.uploadDir, filename)

	// Check if file exists
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	} else if err != nil {
		logger.LogError(r.Context(), err, "Failed to check file")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Determine content type
	ext := strings.ToLower(filepath.Ext(filePath))
	contentType := "application/octet-stream"
	switch ext {
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	case ".gif":
		contentType = "image/gif"
	case ".webp":
		contentType = "image/webp"
	}

	// Serve the file
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	w.Header().Set("Cache-Control", "public, max-age=31536000") // Cache for 1 year

	http.ServeFile(w, r, filePath)
}

func (h *FileHandler) isAllowedType(contentType string) bool {
	for _, allowedType := range h.allowedTypes {
		if contentType == allowedType {
			return true
		}
	}
	return false
}
