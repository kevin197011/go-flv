package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

//go:embed templates/*
var templatesFS embed.FS

//go:embed static/*
var staticFS embed.FS

// ==================== Models ====================

// FlvVideo represents a FLV video entry
type FlvVideo struct {
	ID          uint       `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Name        string     `json:"name"`
	URL         string     `json:"url"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

// CreateVideoRequest represents the request body for creating a video
type CreateVideoRequest struct {
	Name        string `json:"name" binding:"required"`
	URL         string `json:"url" binding:"required"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// ==================== Database ====================

var DB *sql.DB

func initDB() error {
	var err error

	// 获取数据库文件路径，支持环境变量配置
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		// 优先使用 data 目录，如果不存在则创建
		dataDir := "./data"
		if _, err := os.Stat(dataDir); os.IsNotExist(err) {
			os.MkdirAll(dataDir, 0755)
		}
		dbPath = filepath.Join(dataDir, "flv_videos.db")
	}

	// 确保数据库目录存在
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	// 打开数据库连接
	DB, err = sql.Open("sqlite", dbPath+"?_journal_mode=WAL&_foreign_keys=1")
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// 配置连接池
	DB.SetMaxOpenConns(25)                  // 最大打开连接数
	DB.SetMaxIdleConns(5)                   // 最大空闲连接数
	DB.SetConnMaxLifetime(5 * time.Minute)  // 连接最大生命周期
	DB.SetConnMaxIdleTime(10 * time.Minute) // 连接最大空闲时间

	// 测试连接
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// 创建表
	if err := createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}

func createTables() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS flv_videos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		name TEXT NOT NULL,
		url TEXT NOT NULL,
		description TEXT,
		status TEXT DEFAULT 'active',
		deleted_at DATETIME NULL
	);

	CREATE INDEX IF NOT EXISTS idx_flv_videos_status ON flv_videos(status);
	CREATE INDEX IF NOT EXISTS idx_flv_videos_deleted_at ON flv_videos(deleted_at);
	`

	_, err := DB.Exec(createTableSQL)
	return err
}

func dbCreateVideo(video *FlvVideo) error {
	query := `
		INSERT INTO flv_videos (name, url, description, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	if video.Status == "" {
		video.Status = "active"
	}

	result, err := DB.Exec(query, video.Name, video.URL, video.Description, video.Status, now, now)
	if err != nil {
		return fmt.Errorf("failed to create video: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	video.ID = uint(id)
	video.CreatedAt = now
	video.UpdatedAt = now

	return nil
}

func getAllVideos() ([]FlvVideo, error) {
	query := `
		SELECT id, created_at, updated_at, name, url, description, status
		FROM flv_videos
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query videos: %w", err)
	}
	defer rows.Close()

	var videos []FlvVideo
	for rows.Next() {
		var video FlvVideo
		err := rows.Scan(
			&video.ID,
			&video.CreatedAt,
			&video.UpdatedAt,
			&video.Name,
			&video.URL,
			&video.Description,
			&video.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan video: %w", err)
		}
		videos = append(videos, video)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	if videos == nil {
		videos = []FlvVideo{}
	}

	return videos, nil
}

func getVideoByID(id uint) (*FlvVideo, error) {
	query := `
		SELECT id, created_at, updated_at, name, url, description, status
		FROM flv_videos
		WHERE id = ? AND deleted_at IS NULL
	`

	var video FlvVideo
	err := DB.QueryRow(query, id).Scan(
		&video.ID,
		&video.CreatedAt,
		&video.UpdatedAt,
		&video.Name,
		&video.URL,
		&video.Description,
		&video.Status,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("video not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get video: %w", err)
	}

	return &video, nil
}

func updateVideo(video *FlvVideo) error {
	query := `
		UPDATE flv_videos
		SET name = ?, url = ?, description = ?, status = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	now := time.Now()
	result, err := DB.Exec(query, video.Name, video.URL, video.Description, video.Status, now, video.ID)
	if err != nil {
		return fmt.Errorf("failed to update video: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("video not found")
	}

	video.UpdatedAt = now
	return nil
}

func deleteVideo(id uint) error {
	query := `
		UPDATE flv_videos
		SET deleted_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := DB.Exec(query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete video: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("video not found")
	}

	return nil
}

// ==================== Middleware ====================

type AuthConfig struct {
	Username string
	Password string
}

func getAuthConfig() AuthConfig {
	username := os.Getenv("ADMIN_USERNAME")
	password := os.Getenv("ADMIN_PASSWORD")

	if username == "" {
		username = "admin"
	}
	if password == "" {
		password = "admin123"
	}

	return AuthConfig{
		Username: username,
		Password: password,
	}
}

func sessionMiddleware(secret string) gin.HandlerFunc {
	if secret == "" {
		secret = os.Getenv("SESSION_SECRET")
		if secret == "" {
			secret = "go-flv-secret-key-change-in-production"
		}
	}
	store := cookie.NewStore([]byte(secret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 24, // 24 hours
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
	})
	return sessions.Sessions("go-flv-session", store)
}

func isJSONRequest(c *gin.Context) bool {
	return c.GetHeader("Content-Type") == "application/json" ||
		c.GetHeader("Accept") == "application/json" ||
		c.GetHeader("X-Requested-With") == "XMLHttpRequest"
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		authenticated := session.Get("authenticated")

		if authenticated != true {
			if isJSONRequest(c) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error":     "Authentication required",
					"login_url": "/admin/login",
				})
			} else {
				c.Redirect(http.StatusSeeOther, "/admin/login")
			}
			c.Abort()
			return
		}

		c.Next()
	}
}

func requireLogin(c *gin.Context) bool {
	session := sessions.Default(c)
	authenticated := session.Get("authenticated")
	return authenticated == true
}

func login(c *gin.Context, username, password string) bool {
	config := getAuthConfig()

	if username == config.Username && password == config.Password {
		session := sessions.Default(c)
		session.Set("authenticated", true)
		session.Set("username", username)
		session.Save()
		return true
	}

	return false
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许的来源（可以通过环境变量配置）
		origin := c.GetHeader("Origin")
		allowOrigin := os.Getenv("CORS_ORIGIN")
		if allowOrigin == "" {
			allowOrigin = "*" // 默认允许所有来源
		} else if origin != "" && origin == allowOrigin {
			allowOrigin = origin // 允许特定来源
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400") // 24小时

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

// requestLoggerMiddleware 记录请求日志
func requestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 记录日志
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		if statusCode >= 400 {
			log.Printf("[%s] %s %s %d %v %s %s",
				"ERROR",
				method,
				path,
				statusCode,
				latency,
				clientIP,
				errorMessage,
			)
		} else {
			log.Printf("[%s] %s %s %d %v %s",
				"INFO",
				method,
				path,
				statusCode,
				latency,
				clientIP,
			)
		}
	}
}

// recoveryMiddleware 错误恢复中间件
func recoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[PANIC] %s %s - %v", c.Request.Method, c.Request.URL.Path, err)

				if isJSONRequest(c) {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": "Internal server error",
					})
				} else {
					c.HTML(http.StatusInternalServerError, "error.html", gin.H{
						"message": "Internal server error",
					})
				}
				c.Abort()
			}
		}()

		c.Next()
	}
}

// ==================== Handlers ====================

func getPlayerPage(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"URL":       "",
		"HideInput": false,
	})
}

func getVideoPage(c *gin.Context) {
	url := c.Query("url")
	c.HTML(http.StatusOK, "index.html", gin.H{
		"URL":       url,
		"HideInput": url != "",
	})
}

func getMonitorPage(c *gin.Context) {
	c.HTML(http.StatusOK, "monitor.html", gin.H{
		"title": "视频监控页面",
	})
}

func getAdminPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin.html", gin.H{
		"title": "FLV 视频管理",
	})
}

func getLoginPage(c *gin.Context) {
	if requireLogin(c) {
		c.Redirect(http.StatusSeeOther, "/admin")
		return
	}

	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "管理员登录",
	})
}

func postLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if login(c, username, password) {
		c.Redirect(http.StatusSeeOther, "/admin")
		return
	}

	c.HTML(http.StatusUnauthorized, "login.html", gin.H{
		"title":    "管理员登录",
		"error":    "用户名或密码错误，请重试",
		"username": username,
	})
}

func getLogout(c *gin.Context) {
	logout(c)
	c.Redirect(http.StatusSeeOther, "/admin/login")
}

func postLogout(c *gin.Context) {
	logout(c)
	c.Redirect(http.StatusSeeOther, "/admin/login")
}

func getVideos(c *gin.Context) {
	videos, err := getAllVideos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve videos"})
		return
	}
	c.JSON(http.StatusOK, videos)
}

func createVideo(c *gin.Context) {
	var req CreateVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video data: " + err.Error()})
		return
	}

	video := FlvVideo{
		Name:        req.Name,
		URL:         req.URL,
		Description: req.Description,
		Status:      req.Status,
	}

	if err := dbCreateVideo(&video); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create video: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, video)
}

func updateVideoHandler(c *gin.Context) {
	id, err := parseVideoID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	video, err := getVideoByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}

	var req CreateVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video data: " + err.Error()})
		return
	}

	video.Name = req.Name
	video.URL = req.URL
	video.Description = req.Description
	video.Status = req.Status

	if err := updateVideo(video); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update video: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, video)
}

func deleteVideoHandler(c *gin.Context) {
	id, err := parseVideoID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	if err := deleteVideo(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete video: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func parseVideoID(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// ==================== Health Check ====================

func healthCheck(c *gin.Context) {
	// 检查数据库连接
	if err := DB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "unhealthy",
			"message": "Database connection failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"message":   "Service is running",
		"timestamp": time.Now().Unix(),
	})
}

func readinessCheck(c *gin.Context) {
	// 检查数据库连接
	if err := DB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "not ready",
			"message": "Database connection failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ready",
		"message": "Service is ready",
	})
}

// ==================== Routes ====================

func setupRouter() *gin.Engine {
	// 根据环境变量设置Gin模式
	ginMode := os.Getenv("GIN_MODE")
	if ginMode != "" {
		gin.SetMode(ginMode)
	}

	r := gin.Default()

	// 设置请求大小限制（默认32MB，可通过环境变量配置）
	maxBodySize := 32 << 20 // 32MB
	if sizeEnv := os.Getenv("MAX_BODY_SIZE"); sizeEnv != "" {
		if size, err := strconv.ParseInt(sizeEnv, 10, 64); err == nil {
			maxBodySize = int(size)
		}
	}
	r.MaxMultipartMemory = int64(maxBodySize)

	// Add global middleware
	r.Use(sessionMiddleware(""))
	r.Use(corsMiddleware())
	r.Use(requestLoggerMiddleware()) // 请求日志中间件
	r.Use(recoveryMiddleware())      // 错误恢复中间件

	// Health check endpoints (public)
	r.GET("/health", healthCheck)
	r.GET("/ready", readinessCheck)

	// Public web routes (video player)
	r.GET("/", getPlayerPage)
	r.GET("/video", getVideoPage)
	r.GET("/monitor", getMonitorPage)

	// Public API for monitor page (no authentication required)
	r.GET("/public/videos", getVideos)

	// Admin routes
	adminGroup := r.Group("/admin")
	{
		// Public routes (no auth)
		adminGroup.GET("/login", getLoginPage)
		adminGroup.POST("/login", postLogin)
		adminGroup.GET("/logout", getLogout)
		adminGroup.POST("/logout", postLogout)

		// Protected admin page
		adminGroup.Any("", authMiddleware(), getAdminPage)
		adminGroup.Any("/", authMiddleware(), getAdminPage)
	}

	// Protected API routes (require authentication)
	setupVideoRoutes := func(group *gin.RouterGroup) {
		group.GET("/videos", getVideos)
		group.POST("/videos", createVideo)
		group.PUT("/videos/:id", updateVideoHandler)
		group.DELETE("/videos/:id", deleteVideoHandler)
	}

	api := r.Group("/api")
	api.Use(authMiddleware())
	setupVideoRoutes(api)

	// API routes (v1) - also protected
	apiV1 := r.Group("/api/v1")
	apiV1.Use(authMiddleware())
	{
		apiV1.GET("/", getPlayerPage)
		apiV1.GET("/video", getVideoPage)
	}
	setupVideoRoutes(apiV1)

	return r
}

// ==================== Main ====================

func main() {
	// 设置日志格式
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting go-flv server...")

	// 初始化数据库
	if err := initDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Database initialized successfully")

	r := setupRouter()

	// 加载嵌入的模板文件
	tmpl, err := template.ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		log.Fatal("Failed to parse templates:", err)
	}
	r.SetHTMLTemplate(tmpl)

	// 提供嵌入的静态文件
	staticSubFS, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatal("Failed to setup static files:", err)
	}
	r.StaticFS("/static", http.FS(staticSubFS))

	// 启动服务器
	port := ":8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = ":" + envPort
	}

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    port,
		Handler: r,
		// 优化服务器配置
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 在goroutine中启动服务器
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 优雅关闭，给30秒时间完成当前请求
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// 关闭数据库连接
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		} else {
			log.Println("Database connection closed")
		}
	}

	log.Println("Server exited")
}
