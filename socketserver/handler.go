package socketserver

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// serveHome serves the home page and static files
func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)

	// Get the web directory path (relative to where the server is run)
	webDir := "web"

	// Determine which file to serve
	var filePath string
	if r.URL.Path == "/" {
		filePath = filepath.Join(webDir, "index.html")
	} else {
		filePath = filepath.Join(webDir, r.URL.Path)
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	// Serve the file
	http.ServeFile(w, r, filePath)
}
