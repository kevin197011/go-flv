package database

import (
	"database/sql"
	"fmt"
	"go-flv/models"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() error {
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
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

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

// Video CRUD operations
func CreateVideo(video *models.FlvVideo) error {
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

func GetAllVideos() ([]models.FlvVideo, error) {
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

	var videos []models.FlvVideo
	for rows.Next() {
		var video models.FlvVideo
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

	// 确保返回空数组而不是 nil
	if videos == nil {
		videos = []models.FlvVideo{}
	}

	return videos, nil
}

func GetVideoByID(id uint) (*models.FlvVideo, error) {
	query := `
		SELECT id, created_at, updated_at, name, url, description, status
		FROM flv_videos
		WHERE id = ? AND deleted_at IS NULL
	`

	var video models.FlvVideo
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

func UpdateVideo(video *models.FlvVideo) error {
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

func DeleteVideo(id uint) error {
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
