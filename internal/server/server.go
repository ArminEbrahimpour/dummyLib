package server

import (
	"Library/pkg/db"
	"Library/web"
	"net/http"
)

func NewHandler(store *db.Store) http.Handler {
	mux := http.NewServeMux()

	// api
	api := NewAPI(store)
	api.RegisterRoutes(mux)

	//  static files (pdf.js , html , css , pdfs)
	fs := http.FileServer(http.FS(web.WebFS))
	mux.Handle("/", fs)
	return mux
}
