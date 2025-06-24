package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary     Get player page
// @Description Get the main player page
// @Tags        pages
// @Produce     html
// @Success     200 {string} string "HTML page"
// @Router      / [get]
func GetPlayerPage(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"URL":       "",
		"HideInput": false,
	})
}

// @Summary     Get video page
// @Description Get the video page with URL parameter
// @Tags        web
// @Produce     html
// @Param       url query string false "Video URL"
// @Success     200 {string} string "HTML page"
// @Router      /video [get]
func GetVideoPage(c *gin.Context) {
	url := c.Query("url")
	c.HTML(http.StatusOK, "index.html", gin.H{
		"URL":       url,
		"HideInput": url != "",
	})
}

// @Summary     Get monitor page
// @Description Get the monitor page for displaying multiple videos in grid layout
// @Tags        web
// @Produce     html
// @Success     200 {string} string "HTML page"
// @Router      /monitor [get]
func GetMonitorPage(c *gin.Context) {
	c.HTML(http.StatusOK, "monitor.html", gin.H{
		"title": "视频监控页面",
	})
}
