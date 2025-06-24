package handlers

import (
	"net/http"
	"strconv"

	"go-flv/database"
	"go-flv/middleware"
	"go-flv/models"

	"github.com/gin-gonic/gin"
)

// @Summary     Get admin page
// @Description Get the admin management page
// @Tags        admin
// @Produce     html
// @Success     200 {string} string "HTML page"
// @Router      /admin [get]
func GetAdminPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin.html", gin.H{
		"title": "FLV 视频管理",
	})
}

// @Summary     Get login page
// @Description Renders the admin login page
// @Tags        admin
// @Produce     html
// @Success     200 {string} string "HTML page"
// @Router      /admin/login [get]
func GetLoginPage(c *gin.Context) {
	// Check if already logged in
	if middleware.RequireLogin(c) {
		c.Redirect(http.StatusSeeOther, "/admin")
		return
	}

	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "管理员登录",
	})
}

// @Summary     Admin login
// @Description Authenticate admin user and create session
// @Tags        admin
// @Accept      application/x-www-form-urlencoded
// @Param       username formData string true "Username"
// @Param       password formData string true "Password"
// @Success     302 {string} string "Redirect to admin page"
// @Failure     401 {string} string "Login page with error"
// @Router      /admin/login [post]
func PostLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if middleware.Login(c, username, password) {
		// Login successful - use 303 to force GET method
		c.Redirect(http.StatusSeeOther, "/admin")
		return
	}

	// Login failed
	c.HTML(http.StatusUnauthorized, "login.html", gin.H{
		"title":    "管理员登录",
		"error":    "用户名或密码错误，请重试",
		"username": username,
	})
}

// @Summary     Admin logout
// @Description Destroy admin session and redirect to login
// @Tags        admin
// @Success     302 {string} string "Redirect to login page"
// @Router      /admin/logout [post]
func PostLogout(c *gin.Context) {
	middleware.Logout(c)
	c.Redirect(http.StatusSeeOther, "/admin/login")
}

// @Summary     Admin logout (GET)
// @Description Destroy admin session and redirect to login
// @Tags        admin
// @Success     302 {string} string "Redirect to login page"
// @Router      /admin/logout [get]
func GetLogout(c *gin.Context) {
	middleware.Logout(c)
	c.Redirect(http.StatusSeeOther, "/admin/login")
}

// @Summary     Get all FLV videos
// @Description Get list of all FLV videos
// @Tags        admin
// @Produce     json
// @Success     200 {array} models.FlvVideo
// @Router      /api/videos [get]
func GetVideos(c *gin.Context) {
	videos, err := database.GetAllVideos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve videos"})
		return
	}
	c.JSON(http.StatusOK, videos)
}

// @Summary     Create a new FLV video
// @Description Create a new FLV video entry
// @Tags        admin
// @Accept      json
// @Produce     json
// @Param       video body models.CreateVideoRequest true "Video data"
// @Success     201 {object} models.FlvVideo
// @Router      /api/videos [post]
func CreateVideo(c *gin.Context) {
	var req models.CreateVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video data: " + err.Error()})
		return
	}

	video := models.FlvVideo{
		Name:        req.Name,
		URL:         req.URL,
		Description: req.Description,
		Status:      req.Status,
	}

	if err := database.CreateVideo(&video); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create video: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, video)
}

// @Summary     Update a FLV video
// @Description Update an existing FLV video entry
// @Tags        admin
// @Accept      json
// @Produce     json
// @Param       id path int true "Video ID"
// @Param       video body models.CreateVideoRequest true "Video data"
// @Success     200 {object} models.FlvVideo
// @Router      /api/videos/{id} [put]
func UpdateVideo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	video, err := database.GetVideoByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}

	var req models.CreateVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video data: " + err.Error()})
		return
	}

	video.Name = req.Name
	video.URL = req.URL
	video.Description = req.Description
	video.Status = req.Status

	if err := database.UpdateVideo(video); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update video: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, video)
}

// @Summary     Delete a FLV video
// @Description Delete a FLV video entry
// @Tags        admin
// @Produce     json
// @Param       id path int true "Video ID"
// @Success     204
// @Router      /api/videos/{id} [delete]
func DeleteVideo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	if err := database.DeleteVideo(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete video: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
