package handlers

import (
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
// @Description Get the video player page with a specific URL
// @Tags        pages
// @Produce     html
// @Param       url  query    string  true  "Video URL"
// @Success     200 {string} string "HTML page"
// @Router      /video [get]
func GetVideoPage(c *gin.Context) {
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
