package main

import (
	"io/fs"
	"log"
	"net/http"

	"github.com/mrjxtr-dev/score-board/internal/config"
	"github.com/mrjxtr-dev/score-board/internal/routes"
	"github.com/mrjxtr-dev/score-board/internal/store"
)

func main() {
	db := store.LoadDB()
	cfg := config.LoadConfig()

	// Mount embedded static filesystem if available (from assets.go)
	var staticFS http.FileSystem
	if sub, err := fs.Sub(embeddedStatic, "static"); err == nil {
		staticFS = http.FS(sub)
	}

	r := routes.SetupRoutes(cfg, db, staticFS)

	server := &http.Server{
		Addr:    ":" + cfg.PORT,
		Handler: r,
	}

	log.Println("Starting server on port " + cfg.PORT)
	log.Printf("Test connection at http://localhost:%s/ping", cfg.PORT)
	server.ListenAndServe()
}
