package filestore

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/cavidyrm/instawall/config" // <-- Replace with your module name
)

// FileStore handles file upload operations.
type FileStore struct {
	client     *minio.Client
	bucketName string
}

// NewFileStore initializes a new MinIO client and filestore.
func NewFileStore(cfg config.MinIOConfig) (*FileStore, error) {
	ctx := context.Background()

	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	// Create the bucket if it doesn't exist.
	err = minioClient.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, cfg.BucketName)
		if errBucketExists == nil && exists {
			// Bucket already exists, which is fine.
		} else {
			return nil, err // A real error occurred.
		}
	}

	return &FileStore{
		client:     minioClient,
		bucketName: cfg.BucketName,
	}, nil
}

// UploadFile uploads a file to MinIO and returns its URL.
func (fs *FileStore) UploadFile(ctx context.Context, file io.Reader, fileSize int64, originalFilename string) (string, error) {
	// Generate a unique filename to prevent collisions.
	ext := filepath.Ext(originalFilename)
	uniqueFilename := fmt.Sprintf("%s-%s%s", time.Now().Format("20060102"), uuid.New().String(), ext)

	// Upload the file.
	_, err := fs.client.PutObject(ctx, fs.bucketName, uniqueFilename, file, fileSize, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}

	// Construct the URL. The URL format depends on your MinIO setup and domain.
	// This is a typical format for a local setup.
	url := fmt.Sprintf("http://%s/%s/%s", fs.client.EndpointURL().Host, fs.bucketName, uniqueFilename)

	return url, nil
}
