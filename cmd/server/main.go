package main

import (
	"log"

	"urlshortener/internal/adapter/primary/http"
	"urlshortener/internal/adapter/secondary/sqlite"
	"urlshortener/internal/domain/service"
)

func main() {
	db, err := sqlite.InitDB(":memory:")
	if err != nil {
		log.Fatalf("failed to init DB: %v", err)
	}
	defer db.Close()

	repo := sqlite.NewURLRepository(db)
	urlService := service.NewURLService(repo)

	handler := http.NewURLHandler(urlService, "http://localhost:8080")
	router := http.SetupRouter(handler)

	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
