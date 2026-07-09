package main

import (
	"Library/internal/server"
	database "Library/pkg/db"
	"Library/pkg/systray"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func getDbPath() string {

	homePth, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	confDir := filepath.Join(homePth, ".config", "dummylib")
	if err := os.MkdirAll(confDir, 0700); err != nil {
		log.Fatal(err)
	}
	log.Println("config directory created")
	return confDir

}

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	path := getDbPath()
	pathToDb := fmt.Sprintf("%s/bookmarks.db", path)
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
