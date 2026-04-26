package handlers

import (
	"fmt"
	"mime/multipart"
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

// saveFile stores fh in Firebase Storage (if enabled) or local disk.
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
