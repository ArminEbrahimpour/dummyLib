package model

import (
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID         uuid.UUID `json:"id"`
	Path       string    `json:"path"` // path of the book file
	Title      string    `json:"title"`
	PageCount  int       `json:"page_count"`
	LastPage   int       `json:"last_page"`
	LastOpened time.Time `json:"last_opened"` // last modified
	CoverPath  string    `json:"cover_path,omitempty"`
}
