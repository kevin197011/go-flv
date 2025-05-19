package main

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	_ "go-flv/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//go:embed templates
var templatesFS embed.FS

//go:embed static
var staticFS embed.FS

// @title           FLV Player API
// @version         1.0
// @description     A simple FLV video player service
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

// @Summary     Get player page
// @Description Get the main player page
// @Tags        pages
// @Produce     html
// @Success     200 {string} string "HTML page"
// @Router      / [get]
func getPlayerPage(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"URL":       "",
		"HideInput": false,
	})
}

// @Summary     Get video page
// @Description Get the video player page with a specific URL
// @Tags        pages
// @Produce     html
// @Param       url  query    string  true  "Video URL"
// @Success     200 {string} string "HTML page"
// @Router      /video [get]
func getVideoPage(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.Redirect(302, "/")
		return
	}

	c.HTML(200, "index.html", gin.H{
		"URL":       url,
		"HideInput": true,
	})
}

func main() {
	r := gin.Default()
	r.Use(CORSMiddleware())

	// Load the HTML template from the embedded filesystem
	tmpl, err := template.ParseFS(templatesFS, "templates/index.html")
	if err != nil {
		panic(err)
	}
	r.SetHTMLTemplate(tmpl)

	// Serve static files from the embedded filesystem
	subFS, err := fs.Sub(staticFS, "static")
	if err != nil {
		panic(err)
	}
	r.StaticFS("/static", http.FS(subFS))

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := r.Group("/api/v1")
	{
		api.GET("/", getPlayerPage)
		api.GET("/video", getVideoPage)
	}

	// Web routes
	r.GET("/", getPlayerPage)
	r.GET("/video", getVideoPage)

	r.Run(":8080")
}
