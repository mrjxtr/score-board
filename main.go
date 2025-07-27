package main

import (
	"log"
	"net/http"

	"github.com/mrjxtr-dev/score-board/internal/config"
	"github.com/mrjxtr-dev/score-board/internal/routes"
	"github.com/mrjxtr-dev/score-board/internal/store"
)

func main() {
	db := store.LoadDB()
	cfg := config.LoadConfig()
	r := routes.SetupRoutes(cfg, db)

	server := &http.Server{
		Addr:    ":" + cfg.PORT,
		Handler: r,
	}

	log.Println("Starting server on port " + cfg.PORT)
	log.Printf("Test connection at http://localhost:%s/ping", cfg.PORT)
	server.ListenAndServe()
}
