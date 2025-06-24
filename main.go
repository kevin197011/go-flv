package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"go-flv/database"
	_ "go-flv/docs"
	"go-flv/routes"
)

//go:embed templates/*
var templatesFS embed.FS

//go:embed static/*
var staticFS embed.FS

// @title           Go FLV Player API
// @version         1.0
// @description     A video player service that supports FLV format.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

func main() {
	// 初始化数据库
	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	r := routes.SetupRouter()

	// 加载嵌入的模板文件
	tmpl, err := template.ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		panic(err)
	}
	r.SetHTMLTemplate(tmpl)

	// 提供嵌入的静态文件
	staticSubFS, err := fs.Sub(staticFS, "static")
	if err != nil {
		panic(err)
	}
	r.StaticFS("/static", http.FS(staticSubFS))

	r.Run(":8080")
}
