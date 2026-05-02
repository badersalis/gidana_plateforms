package handlers

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/badersalis/gidana_backend/internal/config"
	"github.com/badersalis/gidana_backend/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var allowedImageExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
}

// deleteStorageFile removes a file from Supabase or local disk.
// Errors are logged but not returned — a missing file must not block DB cleanup.
func deleteStorageFile(fileURL string) {
	if fileURL == "" {
		return
	}
	if config.App.SupabaseURL != "" {
		if err := storage.DeleteFile(fileURL); err != nil {
			log.Printf("storage: failed to delete %s: %v", fileURL, err)
		}
		return
	}
	// Local path: "/uploads/properties/uuid.ext" → "./uploads/properties/uuid.ext"
	localPath := filepath.Join(".", fileURL)
	if err := os.Remove(localPath); err != nil && !os.IsNotExist(err) {
		log.Printf("storage: failed to delete local file %s: %v", localPath, err)
	}
}

// saveFile stores fh in Supabase Storage (if configured) or local disk.
// It validates the file extension and returns the URL/path to persist in the DB.
func saveFile(c *gin.Context, fh *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	if !allowedImageExts[ext] {
		return "", fmt.Errorf("invalid file type: only jpg, jpeg, png, gif, webp are allowed")
	}

	if config.App.SupabaseURL != "" {
		return storage.UploadFile(fh, "properties")
	}

	filename := uuid.New().String() + ext
	savePath := filepath.Join(config.App.UploadDir, filename)
	if err := c.SaveUploadedFile(fh, savePath); err != nil {
		return "", fmt.Errorf("failed to save file locally: %w", err)
	}
	return "/uploads/properties/" + filename, nil
}
