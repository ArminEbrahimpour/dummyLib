package main

import (
	"Library/internal/server"
	database "Library/pkg/db"
	"context"
	"log"
	"net/http"
	// "github.com/go-chi/chi/v5"
	// "github.com/go-chi/chi/v5/middleware"
)

func main() {
	// r := chi.NewRouter()
	// r.Use(middleware.Logger)
	// r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("welcome"))
	// })
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db, err := database.InitDB(ctx, "bookmarks.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := database.NewStore(db)

	handler := server.NewHandler(store)

	addr := ":3333"
	log.Printf("Serving on http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}
