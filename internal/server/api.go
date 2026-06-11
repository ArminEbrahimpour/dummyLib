package server

import (
	"Library/internal/model"
	"Library/pkg/db"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type API struct {
	Store *db.Store
}

func NewAPI(s *db.Store) *API {
	return &API{Store: s}
}

func (a *API) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/books", a.handleBooks)
	mux.HandleFunc("/api/books/", a.handleBookByID)
	mux.HandleFunc("/api/open", a.handleOpen)
}

// List al books
func (a *API) handleBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		books, err := a.Store.AllBooks()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(books)
		return
	}
	http.Error(w, "method not allowed", 405)
}

// handleBookByID handles for GET and PUT
func (a *API) handleBookByID(w http.ResponseWriter, r *http.Request) {
	trimmed := strings.TrimPrefix(r.URL.Path, "/api/books/")
	if trimmed == "" {
		http.Error(w, "missing book id", 400)
		return
	}

	// Check if this is a /pdf request
	if strings.HasSuffix(trimmed, "/pdf") {
		idStr := strings.TrimSuffix(trimmed, "/pdf")
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, "invalid id", 400)
			return
		}
		a.handleServePDF(w, r, id)
		return
	}
	// extract ID from url:
	idStr := strings.TrimPrefix(r.URL.Path, "/api/books/")
	if idStr == "" {
		http.Error(w, "missing book id", 400)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", 400)
		return
	}
	switch r.Method {
	case http.MethodGet:
		book, err := a.Store.GetBook(id)
		if err != nil {
			http.Error(w, "not found", 404)
			return
		}
		json.NewEncoder(w).Encode(book)
	case http.MethodPut, http.MethodPost:
		var payload struct {
			Page int `json:"page"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "invalid json", 400)
			return
		}
		if err := a.Store.UpdatePage(id, payload.Page); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "method not allowed", 405)
	}
}

// handleOpen accepts a path to a PDF, registers int in the db, and redirects to the viewer with the book id
func (a *API) handleOpen(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "missing path", 400)
		return
	}
	// verity existence of file
	if _, err := os.Stat(path); err != nil {
		http.Error(w, "file doesn't exist", 404)
		return
	}
	if !strings.HasSuffix(strings.ToLower(path), ".pdf") {
		http.Error(w, "only pdf files are supported", 400)
		return
	}

	// create book id
	id := uuid.New()
	title := filepath.Base(path)

	book := &model.Book{
		ID:    id,
		Path:  path,
		Title: title,
	}
	if err := a.Store.UpSertBook(book); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/view.html?book="+book.ID.String(), http.StatusFound)
}
func (a *API) handleServePDF(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	book, err := a.Store.GetBook(id)
	if err != nil {
		http.Error(w, "book not found", 404)
		return
	}
	// Verify the file still exists on disk
	if _, err := os.Stat(book.Path); err != nil {
		http.Error(w, "file not found on disk", 404)
		return
	}
	w.Header().Set("Content-Type", "application/pdf")
	http.ServeFile(w, r, book.Path)
}
