package utility

import (
	"CamMonitoring/src/service"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// RunServer initializes and runs an HTTP server using the provided Gin router and configuration.
// It also handles graceful shutdown when an interrupt signal is received.
func RunServer(router *gin.Engine, config service.Config) {
	// Create an HTTP server with the specified address and Gin router as the handler
	srv := &http.Server{
		Addr:    config.Server.HttpPort,
		Handler: router,
	}
	// Start the server in a goroutine
	go func() {
		// Attempt to start the server and log an error if it fails
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			ErrorLog().Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt) // Notify the quit channel upon receiving an interrupt signal
	<-quit                            // Wait for the interrupt signal
	InfoLog().Println("Shutting down server...")

	// Create a context with a timeout for the server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), config.Server.ShutdownTimeout*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := srv.Shutdown(ctx); err != nil {
		ErrorLog().Printf("Server forced to shutdown: %v", err)
	}

	InfoLog().Println("Server exiting")
}
