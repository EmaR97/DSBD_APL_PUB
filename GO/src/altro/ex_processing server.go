package main

//
//import (
//	"CamMonitoring/src/service"
//	"bytes"
//	"github.com/minio/minio-go/v7"
//	"github.com/minio/minio-go/v7/pkg/credentials"
//	"gocv.io/x/gocv"
//	"image/color"
//	"log"
//)
//
//func detectPedestrians(img *gocv.Mat, classifier *gocv.CascadeClassifier) {
//	rects := classifier.DetectMultiScale(*img)
//	for _, r := range rects {
//		gocv.Rectangle(img, r, color.RGBA{0, 255, 0, 0}, 3)
//	}
//}
//
//func uploadImageToMinIO_(data []byte) error {
//	// Initialize MinIO client
//	minioClient, err := initializeMinIOClient()
//	if err != nil {
//		log.Println("Error initializing MinIO client:", err)
//		return err
//	}
//
//	// Specify your MinIO bucket and object key
//	bucketName := "your-bucket-name"
//	objectKey := "result.jpg"
//
//	// Upload the result image to MinIO
//	err = uploadImageToMinIO(minioClient, bucketName, objectKey, data)
//	if err != nil {
//		log.Println("Error uploading result image to MinIO:", err)
//		return err
//	}
//
//	log.Println("Result image uploaded to MinIO successfully")
//	return nil
//}
//
//func initializeMinIOClient() (*minio.Client, error) {
//	endpoint := "your-minio-endpoint"
//	accessKeyID := "your-access-key-id"
//	secretAccessKey := "your-secret-access-key"
//	useSSL := false // Change to true if your MinIO server uses SSL
//
//	return minio.New(
//		endpoint, &minio.Options{
//			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
//			Secure: useSSL,
//		},
//	)
//}
//
//func uploadImageToMinIO(minioClient *minio.Client, bucketName, objectKey string, data []byte) error {
//	_, err := minioClient.PutObject(
//		context.Background(), bucketName, objectKey, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{},
//	)
//	return err
//}
//
//func publishMessagesToKafka(producer *service.KafkaProducer) error {
//	// Publish MinIO retrieval message
//	minioMessage := "MinIO retrieval message content"
//	err := producer.Publish("minio-topic", minioMessage)
//	if err != nil {
//		return err
//	}
//
//	// Publish person detection notification message
//	notificationMessage := "Person detection notification content"
//	err = producer.Publish("notification-topic", notificationMessage)
//	return err
//}
