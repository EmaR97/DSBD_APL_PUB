package storage

import (
	"CamMonitoring/src/utility"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"strings"
	"time"
)

// NewMinioClient creates a new MinioClient instance.
func NewMinioClient(endpoint, id, secret, bucketName string, useSSL bool) (*MinioClient, error) {
	client, err := minio.New(
		endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(id, secret, ""),
			Secure: useSSL,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Minio client: %v", err)
	}
	return &MinioClient{client: *client, bucketName: bucketName, endpoint: endpoint}, nil
}

// MinioClient is a client for interacting with the Minio storage service.
type MinioClient struct {
	client     minio.Client
	bucketName string
	endpoint   string
}

// GeneratePresignedURL generates a presigned URL for the specified object.
func (mc *MinioClient) GeneratePresignedURL(objectName string, expiration time.Duration) (string, error) {
	// Set the expiration time for the presigned URL

	// Generate a presigned URL for the file
	URL, err := mc.client.PresignedGetObject(
		context.Background(), mc.bucketName, objectName, expiration, nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %v", err)
	}

	minioToReplace := mc.endpoint
	minioReplacement := "localhost"
	utility.InfoLog().Println("URL: ", URL.String())

	// Perform the replacements
	imageUrl := strings.ReplaceAll(URL.String(), minioToReplace, minioReplacement)
	utility.InfoLog().Println("URL modified: ", imageUrl)

	return imageUrl, nil
}

// DeleteObject deletes the specified object from the Minio bucket.
func (mc *MinioClient) DeleteObject(objectName string) error {
	err := mc.client.RemoveObject(context.Background(), mc.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object: %v", err)
	}
	return nil
}
