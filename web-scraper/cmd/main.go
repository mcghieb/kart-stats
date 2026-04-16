package main

import (
	"log"

	"github.com/mcghieb/kart-stats/web-scraper/internal/cache"
	"github.com/mcghieb/kart-stats/web-scraper/internal/config"
	"github.com/mcghieb/kart-stats/web-scraper/internal/repository/providers"
	"github.com/mcghieb/kart-stats/web-scraper/internal/scripts"
)

func main() {
	cfg := config.LoadConfig()

	log.Printf("[BEGIN]: Downloading Race Data HTML Pages\n")
	scripts.DownloadHTMLpages()
	log.Printf("[END]: Downloading Race Data HTML Pages\n")

	store := cache.NewCache()

	log.Printf("[BEGIN]: Parsing Race Data Pages\n")
	scripts.ParseRacePages(store)
	log.Printf("[END]: Parsing Race Data Pages\n")

	log.Printf("[BEGIN]: Parsing Driver Pages\n")
	scripts.ParseDriverPages(store)
	log.Printf("[END]: Parsing Driver Pages\n")

	repo, err := providers.NewPostgresRepository(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer repo.Close()

	if err := repo.Migrate(); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	log.Printf("[BEGIN]: Uploading Values To Database\n")
	if err := scripts.UploadToDatabase(store, repo); err != nil {
		log.Fatalf("failed to upload to database: %v", err)
	}
	log.Printf("[END]: Uploading Values To Database\n")
}
