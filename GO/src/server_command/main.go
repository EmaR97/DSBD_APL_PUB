package main

import (
	"CamMonitoring/src/service/comunication"
	"CamMonitoring/src/service/middleware"
	"CamMonitoring/src/utility"
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
	utility.RunServer(router, config)

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
