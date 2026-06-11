package db

import (
	"Library/internal/model"
	"database/sql"
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
	book.LastOpened = time.Now()
	_, err := s.db.Exec(`
	INSERT INTO books (id, path, title, page_count, last_page, last_opened)
		VALUES(?, ?, ?, ?, ?, ?)
		ON CONFLICT(path) DO UPDATE SET
			title = excluded.title,
			page_count = exluded.page_count,
			last_page = excluded.last_page,
			last_opened = excluded.last_opened`, book.ID.String(), book.Path, book.Title, book.PageCount, book.LastPage, book.LastOpened)
	return err
}

// UpdatePage updates the last page for a given book ID
func (s *Store) UpdatePage(id uuid.UUID, page int) error {
	_, err := s.db.Exec(`UPDATE books SET last_page = ?, last_opened = ? WHERE id = ?`,
		page, time.Now(), id.String())
	return err
}

// GetBook retrieves a single book by ID.
func (s *Store) GetBook(id uuid.UUID) (*model.Book, error) {
	row := s.db.QueryRow(`SELECT id, path, title, page_count, last_page, last_opened, cover_path FROM books WHERE id = ?`, id.String())
	b := &model.Book{}
	var idStr string

	err := row.Scan(&idStr, &b.Path, b.Title, b.PageCount, b.LastPage, b.LastOpened, b.CoverPath)
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
		var idStr string
		if err := rows.Scan(&idStr, b.Path, b.Title, b.PageCount, b.LastPage, b.LastOpened, b.CoverPath); err != nil {
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
