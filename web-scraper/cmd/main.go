package main

import (
	"log"

	"github.com/mcghieb/kart-stats/web-scraper/internal/cache"
	"github.com/mcghieb/kart-stats/web-scraper/internal/scripts"
)

func main() {
	log.Printf("[BEGIN]: Downloading Race Data HTML Pages\n")
	scripts.DownloadHTMLpages()
	log.Printf("[BEGIN]: Downloading Race Data HTML Pages\n")

	store := cache.NewCache()

	log.Printf("[BEGIN]: Parsing Race Data Pages\n")
	scripts.ParseRacePages(store)
	log.Printf("[END]: Parsing Race Data Pages\n")

	log.Printf("[BEGIN]: Parsing Driver Pages\n")
	scripts.ParseDriverPages(store)
	log.Printf("[END]: Parsing Driver Pages\n")
}
