package main

import (
	"github.com/mcghieb/kart-stats/web-scraper/internal/cache"
	"github.com/mcghieb/kart-stats/web-scraper/internal/scripts"
)

func main() {
	// scripts.DownloadHTMLpages()
	store := cache.NewCache()
	scripts.ParseRaceData(store)
}
