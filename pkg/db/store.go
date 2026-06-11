package db

import (
	"Library/internal/model"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

// UpSertBook update or insert the book if it doesn't exists
func (s *Store) UpSertBook(book *model.Book) error {
	book.LastOpened = time.Now().UTC().Round(time.Second)

	_, err := s.db.Exec(`
		INSERT INTO books (id, path, title, page_count, last_page, last_opened)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(path) DO UPDATE SET
			id = excluded.id,
			title = excluded.title,
			page_count = excluded.page_count,
			last_opened = excluded.last_opened
	`, book.ID.String(), book.Path, book.Title, book.PageCount, book.LastPage, book.LastOpened)
	if err != nil {
		return err
	}
	var idStr string
	err = s.db.QueryRow(`SELECT id FROM books WHERE path = ?`, book.Path).Scan(&idStr)
	//log.Printf("DEBUG UpSertBook: read-back id = %s", idStr)
	if err != nil {
		return err
	}
	book.ID, err = uuid.Parse(idStr)
	return err
}

// UpdatePage updates the last page for a given book ID
func (s *Store) UpdatePage(id uuid.UUID, page int) error {
	_, err := s.db.Exec(`UPDATE books SET last_page = ?, last_opened = ? WHERE id = ?`,
		page, time.Now().UTC().Round(time.Second), id.String())
	return err
}

// GetBook retrieves a single book by ID.
func (s *Store) GetBook(id uuid.UUID) (*model.Book, error) {
	row := s.db.QueryRow(`SELECT id, path, title, page_count, last_page, last_opened, cover_path FROM books WHERE id = ?`, id.String())
	b := &model.Book{}
	var idStr, lastOpenedStr string
	var coverPath sql.NullString
	err := row.Scan(&idStr, &b.Path, &b.Title, &b.PageCount, &b.LastPage, &lastOpenedStr, &coverPath)
	if err != nil {
		return nil, err
	}
	b.CoverPath = coverPath.String
	b.LastOpened, err = parseTime(lastOpenedStr)
	if err != nil {
		return nil, err
	}
	b.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// AllBooks returns all books ordered by last opened.
func (s *Store) AllBooks() ([]model.Book, error) {
	rows, err := s.db.Query(`
	SELECT id, path, title, page_count, last_page, last_opened, cover_path from books ORDER BY last_opened DESC;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var books []model.Book
	for rows.Next() {
		var b model.Book
		var idStr, lastOpenedStr string
		var coverPath sql.NullString
		if err := rows.Scan(&idStr, &b.Path, &b.Title, &b.PageCount, &b.LastPage, &lastOpenedStr, &coverPath); err != nil {
			return nil, err
		}
		b.CoverPath = coverPath.String
		b.LastOpened, err = parseTime(lastOpenedStr)

		if err != nil {
			return nil, err
		}
		b.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}
		books = append(books, b)

	}
	return books, nil
}

// Add a helper at the top of store.go
func parseTime(s string) (time.Time, error) {
	formats := []string{
		"2006-01-02 15:04:05 +0000 UTC", // ← add this as first (most common now)
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05.999999999-07:00",
		"2006-01-02 15:04:05Z07:00",
		"2006-01-02 15:04:05",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("cannot parse time: %q", s)
}
