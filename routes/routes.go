package routes

import (
	"go-flv/handlers"
	"go-flv/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter configures all the routes for the application
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := r.Group("/api/v1")
	{
		api.GET("/", handlers.GetPlayerPage)
		api.GET("/video", handlers.GetVideoPage)
	}

	// Web routes
	r.GET("/", handlers.GetPlayerPage)
	r.GET("/video", handlers.GetVideoPage)

	return r
}
