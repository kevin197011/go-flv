package main

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed templates
var templatesFS embed.FS

//go:embed static
var staticFS embed.FS

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

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"URL":       "",
			"HideInput": false,
		})
	})

	r.GET("/video", func(c *gin.Context) {
		url := c.Query("url")
		if url == "" {
			c.Redirect(302, "/")
			return
		}

		c.HTML(200, "index.html", gin.H{
			"URL":       url,
			"HideInput": true,
		})
	})

	r.Run(":8080")
}
