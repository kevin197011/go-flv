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

	// Add session middleware
	r.Use(middleware.SessionMiddleware(""))
	r.Use(middleware.CORSMiddleware())

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public API for monitor page (no authentication required) - different path to avoid conflicts
	r.GET("/public/videos", handlers.GetVideos)

	// Admin routes
	adminGroup := r.Group("/admin")
	{
		// Public routes (no auth)
		adminGroup.GET("/login", handlers.GetLoginPage)
		adminGroup.POST("/login", handlers.PostLogin)
		adminGroup.GET("/logout", handlers.GetLogout)
		adminGroup.POST("/logout", handlers.PostLogout)
	}

	// Protected admin routes - use individual routes instead of group
	r.GET("/admin", middleware.AuthMiddleware(), handlers.GetAdminPage)
	r.GET("/admin/", middleware.AuthMiddleware(), handlers.GetAdminPage)

	// Handle POST requests to /admin by serving the same admin page (for browser redirects after login)
	r.POST("/admin", middleware.AuthMiddleware(), handlers.GetAdminPage)
	r.POST("/admin/", middleware.AuthMiddleware(), handlers.GetAdminPage)

	// Protected API routes (require authentication)
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// Video management API
		api.GET("/videos", handlers.GetVideos)
		api.POST("/videos", handlers.CreateVideo)
		api.PUT("/videos/:id", handlers.UpdateVideo)
		api.DELETE("/videos/:id", handlers.DeleteVideo)
	}

	// API routes (v1) - also protected
	apiV1 := r.Group("/api/v1")
	apiV1.Use(middleware.AuthMiddleware())
	{
		apiV1.GET("/", handlers.GetPlayerPage)
		apiV1.GET("/video", handlers.GetVideoPage)

		// Video management API
		apiV1.GET("/videos", handlers.GetVideos)
		apiV1.POST("/videos", handlers.CreateVideo)
		apiV1.PUT("/videos/:id", handlers.UpdateVideo)
		apiV1.DELETE("/videos/:id", handlers.DeleteVideo)
	}

	// Public web routes (video player)
	r.GET("/", handlers.GetPlayerPage)
	r.GET("/video", handlers.GetVideoPage)
	r.GET("/monitor", handlers.GetMonitorPage)

	return r
}
