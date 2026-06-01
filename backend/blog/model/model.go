package model

import (
	"time"

	"github.com/lib/pq"
)

type Blog struct {
	ID             string         `gorm:"primaryKey" json:"id"`
	AuthorId       string         `json:"author_id"`
	Title          string         `json:"title"`
	Description    string         `json:"description"`
	Images         pq.StringArray `gorm:"type:text[]" json:"images,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	Likes          int64          `gorm:"-" json:"likes"`
	AuthorUsername string         `gorm:"-" json:"authorUsername,omitempty"`
}

type Like struct {
	ID     uint   `gorm:"primaryKey"`
	UserId string `gorm:"index:idx_user_blog,unique"`
	BlogId string `gorm:"index:idx_user_blog,unique"`
}
type Comment struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	BlogID    string    `json:"blog_id"`
	AuthorID  string    `json:"author_id"`
	AuthorUsername string    `gorm:"-" json:"authorUsername,omitempty"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	EditedAt  time.Time `json:"edited_at,omitempty"`
}
