package main

import (
	"Library/internal/server"
	database "Library/pkg/db"
	"Library/pkg/systray"
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	pathToDb := filepath.Join(dir, "bookmarks.db")
	db, err := database.InitDB(ctx, pathToDb)

	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("Using database file at: %s", dbPath)
	defer db.Close()

	store := database.NewStore(db)

	handler := server.NewHandler(store)

	addr := ":3333"
	log.Printf("Serving on http://localhost%s", addr)
	go func() {
		log.Fatal(http.ListenAndServe(addr, handler))
	}()

	systray.Start()

	// wait for quit signal from tray
	<-systray.GetQuitChan()
	log.Println("shutting down...")
}
