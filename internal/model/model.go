package model

import (
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID         uuid.UUID
	Path       string // path of the book file
	Title      string
	PageCount  int
	LastPage   int
	LastOpened *time.Time // last modified
	CoverPath  string
}
