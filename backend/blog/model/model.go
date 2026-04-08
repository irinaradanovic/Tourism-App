package model

import (
	"time"

	"github.com/lib/pq"
)

type Blog struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	AuthorId    string         `json:"author_id"` // ili int?
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Images      pq.StringArray `gorm:"type:text[]" json:"images,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	Likes       int            `json:"likes,omitempty"`

	//dodati posle komentare
}
