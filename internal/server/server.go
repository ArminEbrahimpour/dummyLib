package server

import (
	"Library/web"
	"net/http"
)

func NewHandler() http.Handler {
	mux := http.NewServeMux()

	fs := http.FileServer(http.FS(web.WebFS))
	mux.Handle("/", fs)
	return mux
}
