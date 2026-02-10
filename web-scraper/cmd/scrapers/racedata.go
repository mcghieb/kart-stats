package main

import (
	"fmt"

	"github.com/mcghieb/kart-stats/web-scraper/internal/cache"
	"github.com/mcghieb/kart-stats/web-scraper/internal/collector"
	"github.com/mcghieb/kart-stats/web-scraper/internal/scrapers/racedata"
)

func main() {
	// TODO: get a basic heat scrape script up and running

	// TODO: control script (this one is the one that does CDC and all other heats)
	// determine where to start
	// determine when to stop (could use the avg number of races an hour compared to the starting point)
	// Async visit every url between start and stop
	// scrape from start_heatno to stop_heatno

	fmt.Println("BEGIN")

	c := collector.NewCollector()
	cache := cache.NewCache()
	racedata.AttachHandlers(c, cache)

	// TODO: make url creation utility
	if err := c.Visit("https://rrorem.clubspeed.com/sp_center/HeatDetails.aspx?HeatNo=39299"); err != nil {
		fmt.Println("Scraper encountered an error on heat 39299")
	}

	// Print the cached race data for heat 39299
	race := cache.Race.Get("39299")
	fmt.Println(race)

	// print driver cache
	fmt.Println(cache.Driver)

	fmt.Println("END")

	// TODO: AT THE END OF THE SCRIPT:
	// go through all of the drivers in the cache
	// and update their proskills in the db. This
	// avoids us needing to make a bunch of updates
	// mid process when there could be a lot of
	// separate races of the same drivers
}
