package main

import (
	"CamMonitoring/src/service/comunication"
	"CamMonitoring/src/service/middleware"
	"CamMonitoring/src/utility"
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"

	"CamMonitoring/src/handler"
	"CamMonitoring/src/service"
)

func main() {
	// Load configuration using config package
	config := service.LoadConfig()

	// Initialize BrokerManagement service
	brokerManagement := initBrokerManagement(config)
	// Test user
	_ = brokerManagement.CreateUser("test", "test")
	_ = brokerManagement.SetPermission("test")
	_ = brokerManagement.SetTopicPermission("test", []string{"command.*"})

	// Create Gin router
	router := setupRouter(config)

	// Initialize MQTT client
	mqttClient := initMqttClient(config.MQTT)

	// Create CommandHandler
	commandHandler := handler.NewCommandHandler(mqttClient)

	// Set up routes
	setUpRoutes(router, commandHandler)

	// Run the server
	runServer(router, config)

}

func initBrokerManagement(config service.Config) (brokerManagement *comunication.BrokerManagement) {
	brokerManagement = comunication.NewBrokerManagement(
		config.Broker.Url, config.Broker.Username, config.Broker.Password,
	)
	return
}

func setupRouter(config service.Config) (router *gin.Engine) {
	router = gin.New()
	// Middleware for authentication
	router.Use(
		middleware.AuthMiddleware(
			config.Server.AuthEndpoint, []string{},
		), /*,service.AuthorizedIPMiddleware(config.Server.AllowedIPs) */
		gin.Recovery(), gin.LoggerWithFormatter(middleware.CustomLogFormatter),
	)
	return
}

func initMqttClient(config service.MQTTConfig) (mqttClient *comunication.MqttClient) {

	mqttClient, err := comunication.NewMqttClient(config.BrokerURL, config.ClientID, config.Username, config.Password)
	if err != nil {
		utility.ErrorLog().Fatalf("Error initializing MQTT client: %v", err)
	}
	return
}

func setUpRoutes(router *gin.Engine, commandHandler *handler.CommandHandler) {
	router.POST("/commands/:id", commandHandler.Post)
}

func runServer(router *gin.Engine, config service.Config) {
	srv := &http.Server{
		Addr:    config.Server.HttpPort, /*config.Server.Address*/
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			utility.ErrorLog().Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	utility.InfoLog().Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), config.Server.ShutdownTimeout*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		utility.ErrorLog().Fatal("Server forced to shutdown:", err)
	}

	utility.InfoLog().Println("Server exiting")
}
