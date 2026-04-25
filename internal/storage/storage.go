package storage

import (
	"context"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	firebase "firebase.google.com/go/v4"
	cloudstorage "cloud.google.com/go/storage"
	"google.golang.org/api/option"

	"github.com/badersalis/gidana_backend/internal/config"
	"github.com/google/uuid"
)

var bucket *cloudstorage.BucketHandle

// Init initialises the Firebase Storage client. Must be called after config.Load().
func Init() error {
	ctx := context.Background()

	var opt option.ClientOption
	if config.App.FirebaseCredJSON != "" {
		opt = option.WithCredentialsJSON([]byte(config.App.FirebaseCredJSON))
	} else if config.App.FirebaseCredPath != "" {
		opt = option.WithCredentialsFile(config.App.FirebaseCredPath)
	} else {
		return fmt.Errorf("no Firebase credentials provided (set FIREBASE_CREDENTIALS_JSON or FIREBASE_CREDENTIALS_PATH)")
	}

	app, err := firebase.NewApp(ctx, &firebase.Config{
		StorageBucket: config.App.FirebaseBucket,
	}, opt)
	if err != nil {
		return fmt.Errorf("firebase.NewApp: %w", err)
	}

	sc, err := app.Storage(ctx)
	if err != nil {
		return fmt.Errorf("app.Storage: %w", err)
	}

	bucket, err = sc.DefaultBucket()
	if err != nil {
		return fmt.Errorf("DefaultBucket: %w", err)
	}

	return nil
}

// UploadFile uploads a multipart file to Firebase Storage under the given folder
// and returns its public URL. The file is made publicly readable.
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

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	obj := bucket.Object(objectPath)
	wc := obj.NewWriter(ctx)
	wc.ContentType = contentType

	if _, err := io.Copy(wc, src); err != nil {
		_ = wc.Close()
		return "", fmt.Errorf("upload copy: %w", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("upload close: %w", err)
	}

	// Make publicly readable so the URL works without a token.
	if err := obj.ACL().Set(ctx, cloudstorage.AllUsers, cloudstorage.RoleReader); err != nil {
		return "", fmt.Errorf("set ACL: %w", err)
	}

	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", config.App.FirebaseBucket, objectPath), nil
}
