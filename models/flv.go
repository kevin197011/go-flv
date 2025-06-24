package models

import (
	"time"
)

// FlvVideo represents a FLV video entry
type FlvVideo struct {
	ID          uint       `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Name        string     `json:"name" example:"示例视频"`
	URL         string     `json:"url" example:"http://example.com/video.flv"`
	Description string     `json:"description" example:"这是一个示例视频"`
	Status      string     `json:"status" example:"active"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

// CreateVideoRequest represents the request body for creating a video
type CreateVideoRequest struct {
	Name        string `json:"name" binding:"required" example:"示例视频"`
	URL         string `json:"url" binding:"required" example:"http://example.com/video.flv"`
	Description string `json:"description" example:"这是一个示例视频"`
	Status      string `json:"status" example:"active"`
}

func (FlvVideo) TableName() string {
	return "flv_videos"
}
