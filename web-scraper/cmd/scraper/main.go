package main

import (
	"fmt"

	"github.com/mcghieb/kart-stats/web-scraper/internal/scraper"
)

func main() {
	// TODO: get a basic heat scrape script up and running

	// TODO: control script (this one is the one that does CDC and all other heats)
	// determine where to start
	// determine when to stop (could use the avg number of races an hour compared to the starting point)
	// Async visit every url between start and stop
	// scrape from start_heatno to stop_heatno

	fmt.Println("BEGINNING")

	c := scraper.NewCollector()
	scraper.AttachHandlers(c)

	// HAVE IT VISIT A DUMMY PAGE
	if err := c.Visit("https://rrorem.clubspeed.com/sp_center/HeatDetails.aspx?HeatNo=39299"); err != nil {
		fmt.Println("Scraper encountered an error on heat 39299")
	}

	fmt.Println("END")
}
