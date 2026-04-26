package storage

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/badersalis/gidana_backend/internal/config"
	"github.com/google/uuid"
)

// UploadFile uploads a multipart file to Supabase Storage under the given folder
// and returns its public URL. The bucket must already be set to public in Supabase.
func UploadFile(fh *multipart.FileHeader, folder string) (string, error) {
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	objectPath := fmt.Sprintf("%s/%s%s", folder, uuid.New().String(), ext)

	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	src, err := fh.Open()
	if err != nil {
		return "", fmt.Errorf("open: %w", err)
	}
	defer src.Close()

	body, err := io.ReadAll(src)
	if err != nil {
		return "", fmt.Errorf("read file: %w", err)
	}

	uploadURL := fmt.Sprintf("%s/storage/v1/object/%s/%s",
		config.App.SupabaseURL, config.App.SupabaseBucket, objectPath)

	req, err := http.NewRequest(http.MethodPost, uploadURL, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+config.App.SupabaseKey)
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("upload: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		msg, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("supabase upload failed (%d): %s", resp.StatusCode, string(msg))
	}

	return fmt.Sprintf("%s/storage/v1/object/public/%s/%s",
		config.App.SupabaseURL, config.App.SupabaseBucket, objectPath), nil
}
