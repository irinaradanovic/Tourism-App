package model

import "time"

type Blog struct {
	ID          string    `json:"id"`
	AuthorId    string    `json:"author_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Images      []string  `json:"images,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	Likes       int       `json:"likes,omitempty"`

	//dodati posle komentare
}
