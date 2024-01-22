package main

import (
	"CamMonitoring/src/entity"
	"CamMonitoring/src/handler"
	"CamMonitoring/src/repository"
	"CamMonitoring/src/service"
	"CamMonitoring/src/service/middleware"
	"CamMonitoring/src/service/storage"
	"CamMonitoring/src/utility"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func main() {
	// Load configuration using config package
	config := service.LoadConfig()

	// Initialize MongoDB service
	mongoDBService := initMongoDBService(config)
	defer func(mongoDBService *storage.MongoDBService) {
		_ = mongoDBService.Close()
	}(mongoDBService)

	// Create user repository
	userRepository := repository.NewMongoDBRepository[entity.User](mongoDBService, config.Database.Name, nil)
	tokenRepository := repository.NewMongoDBRepository[entity.Token](mongoDBService, config.Database.Name, nil)

	// Create Gin router
	router := setupRouter()

	// Initialize services with dependency injection
	accessHandler := handler.NewAccessHandler(userRepository, tokenRepository)

	// Set up routes
	setUpRoutes(router, accessHandler, config.Server.TemplateDir)

	// Run the server
	utility.RunServer(router, config)
}

func initMongoDBService(config service.Config) (mongoDBService *storage.MongoDBService) {
	databaseURL := config.Database.URL
	mongoDBService, err := storage.NewMongoDBService(databaseURL)
	if err != nil {
		utility.ErrorLog().Printf("Failed to initialize MongoDB service: %v", err)
		os.Exit(1)
	}
	return mongoDBService
}
func setupRouter() (router *gin.Engine) {
	router = gin.New()
	// Middleware for authentication
	router.Use(gin.Recovery(), gin.LoggerWithFormatter(middleware.CustomLogFormatter))
	return
}
func setUpRoutes(router *gin.Engine, accessHandler *handler.AccessHandler, templateDir string) {

	access := router.Group("/access")
	{
		access.GET(
			"/login", func(c *gin.Context) {
				c.HTML(http.StatusOK, "login.html", nil)
			},
		)
		access.POST("/login", accessHandler.LoginPost)
		access.GET(
			"/signup", func(c *gin.Context) {
				c.HTML(http.StatusOK, "signup.html", nil)
			},
		)
		access.POST("/signup", accessHandler.SignupPost)
		access.GET("/logout", accessHandler.Logout)       //http://localhost:8080/access/logout
		access.POST("/verify", accessHandler.VerifyToken) //http://localhost:8080/access/verify
	}

	router.LoadHTMLGlob(templateDir)
}
