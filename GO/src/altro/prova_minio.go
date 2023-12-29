package main

import (
	"CamMonitoring/src/service/storage"
	"fmt"
	"log"
)

func main() {
	endpoint := "0.0.0.0:9000"
	bucketName := "your-bucket"
	minioClient, err := storage.NewMinioClient(endpoint, "username", "password", bucketName, false)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Set the bucket name and file name
	fileName := "6579e33559a6784fbf6c3bf1/1702969445.150402148.jpg"

	// Download the file from Minio
	presignedURL, err := minioClient.GeneratePresignedURL(fileName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Access Url: %s", presignedURL)

}
