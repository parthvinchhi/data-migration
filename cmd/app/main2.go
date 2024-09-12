package main

import (
	"github.com/gin-gonic/gin"
	"github.com/parthvinchhi/data-migration/pkg/routes"
	"github.com/parthvinchhi/data-migration/pkg/services"
)

func main() {
	// Initialize Gin router
	router := gin.Default()

	// Serve static files
	router.Static("/static", "./pkg/static")

	// Initialize services
	storageService := &services.MockStorageService{} // or use your real implementation

	// Initialize handlers with injected dependencies
	handler := routes.NewHandler(storageService)

	// Load routes with handlers
	router.GET("/", handler.ServeIndex)
	router.POST("/migrate", handler.HandleMigrate)

	// Start the server
	router.Run(":8080")
}
