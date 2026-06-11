package main

import (
	"Library/internal/server"
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
	handler := server.NewHandler()

	addr := ":3333"
	log.Printf("Serving on http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}
