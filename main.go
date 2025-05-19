package main

import (
	_ "go-flv/docs"
	"go-flv/routes"
)

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
	r := routes.SetupRouter()
	r.LoadHTMLGlob("templates/index.html")
	r.Static("/static", "./static")
	r.Run(":8080")
}
