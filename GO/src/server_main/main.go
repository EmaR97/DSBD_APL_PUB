package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"CamMonitoring/src/_interface"
	"CamMonitoring/src/entity"
	"CamMonitoring/src/handler"
	"CamMonitoring/src/repository"
	"CamMonitoring/src/service"
	"CamMonitoring/src/service/comunication"
	"CamMonitoring/src/service/middleware"
	"CamMonitoring/src/service/storage"
	"CamMonitoring/src/utility"
)

func main() {

	// Load configuration using config package
	config := service.LoadConfig()
	utility.InfoLog().Println("config.Kafka.Brokers: ", config.Kafka.Brokers)

	// Initialize BrokerManagement service
	initBrokerManagement(config)

	// Initialize MongoDB service
	mongoDBService := initMongoDBService(config)
	defer func(mongoDBService *storage.MongoDBService) {
		_ = mongoDBService.Close()
	}(mongoDBService)

	// Create Gin router
	router := setupRouter(config)

	// Create receivedFrames repository
	cameraRepository := repository.NewMongoDBRepository[entity.Camera](mongoDBService, config.Database.Name, nil)

	// Initialize Minio client
	minioClient, err := storage.NewMinioClient(
		config.Minio.Endpoint, config.Minio.Username, config.Minio.Password, config.Minio.BucketName, false,
	)
	if err != nil {
		utility.ErrorLog().Fatalf("Failed to initialize Minio client: %v", err)
	}
	// Initialize services with dependency injection
	cameraHandler := handler.CameraExernalRequestHandler{
		Handler:           *_interface.NewHandler[entity.Camera](cameraRepository),
		GenerateImageUrl_: minioClient.GeneratePresignedURL,
	}

	videoFeedHandler := handler.NewVideoFeedHandler(
		&repository.CameraRepository{
			MongoDBRepository: *cameraRepository,
		},
	)

	// Initialize cleanup service
	cleanupService := storage.NewCleanupService(
		cameraRepository, time.Duration(config.Minio.CleanUpInterval)*time.Hour,
		time.Duration(config.Minio.CleanUpOlderThan)*time.Hour, minioClient,
	)
	go cleanupService.Start()
	defer cleanupService.Stop()

	// Initialize NewKafkaProducer service
	kafkaProducer := comunication.NewKafkaProducer(config.Kafka.Brokers)
	defer kafkaProducer.Close()
	sendNotification := func(jsonString string, key string) error {
		return kafkaProducer.Publish(config.Kafka.TopicOut, jsonString, key)
	}
	sendResponse := func(response, topic string) error {
		return kafkaProducer.Publish(topic, response, "")
	}
	generateImageUrl := func(id string, timestamp int64) string {
		return fmt.Sprintf("http://localhost%s/api/videoFeed/%s?lastseen=%d", config.Server.HttpPort, id, timestamp)
	}
	// Create and start FrameInfoStoring service
	frameInfoStoring := handler.NewFrameHandler(
		&repository.CameraRepository{MongoDBRepository: *cameraRepository}, sendNotification, generateImageUrl,
		sendResponse,
	)
	topicCallbacks := map[string]func(string2 string) error{
		config.Kafka.TopicProcessAndStoreFrame: frameInfoStoring.ProcessAndStoreFrame,
		config.Kafka.TopicGetCamIds:            frameInfoStoring.GetCamIds,
	}
	// Initialize Kafka consumer
	consumer := comunication.NewKafkaConsumer(config.Kafka.Brokers, config.Kafka.GroupID, topicCallbacks)
	go consumer.Start()

	// Set up routes
	setUpRoutes(router, cameraHandler, videoFeedHandler, config.Server.TemplateDir)

	// Run the server
	utility.RunServer(router, config)
	// To close in the right order
	consumer.Stop()
	consumer.Close()

}

func initBrokerManagement(config service.Config) {
	brokerManagement := comunication.NewBrokerManagement(
		config.Broker.Url, config.Broker.Username, config.Broker.Password,
	)
	// Test user
	_ = brokerManagement.CreateUser("test", "test")
	_ = brokerManagement.SetPermission("test")
	_ = brokerManagement.SetTopicPermission("test", []string{"command.*"})
}

func initMongoDBService(config service.Config) *storage.MongoDBService {
	databaseURL := config.Database.URL
	mongoDBService, err := storage.NewMongoDBService(databaseURL)
	if err != nil {
		utility.ErrorLog().Fatalf("Failed to initialize MongoDB service: %v", err)
	}
	return mongoDBService
}

func setupRouter(config service.Config) *gin.Engine {
	router := gin.New()

	// Middleware for authentication
	router.Use(
		//middleware.AuthorizedIPMiddleware(config.Server.AllowedIPs),
		middleware.AuthMiddleware(config.Server.AuthEndpoint, []string{"/api/camera/login"}), gin.Recovery(),
		gin.LoggerWithFormatter(middleware.CustomLogFormatter),
	)
	return router
}

func setUpRoutes(
	router *gin.Engine, receivedFramesHandler handler.CameraExernalRequestHandler,
	videoHandler *handler.VideoFeedHandler, templateDir string,
) {
	api := router.Group("/api")
	{
		api.GET("/camera/:id/:lastSeen", receivedFramesHandler.Next)
		api.POST("/camera", receivedFramesHandler.Create)
		api.POST("/camera/login", receivedFramesHandler.CameraLogin)

		api.GET("/videoFeed/:id", videoHandler.Cam)
		api.GET("/videoFeed", videoHandler.Index) //http://localhost:8080/api/videoFeed
	}

	router.LoadHTMLGlob(templateDir)
}
