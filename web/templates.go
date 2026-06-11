package web

import "embed"

// WebFS holds all static frontend files (HTML, PDF.js, test PDF).
//
//go:embed *
var WebFS embed.FS
