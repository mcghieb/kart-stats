package scraper

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/mcghieb/kart-stats/web-scraper/internal/models"
	"github.com/mcghieb/kart-stats/web-scraper/internal/scraper/cache"
)

// AttachHandlers attaches the HEAT SCRAPING HANDLERS. Name will change if a different handler set is needed for other scraping.
func AttachHandlers(c *colly.Collector, cache *cache.Cache) {
	c.OnHTML("a[href^='RacerHistory']", func(e *colly.HTMLElement) {
		handleCreateUser(e, cache.Driver)
	})

	c.OnHTML("table.RaceResults > tbody", func(e *colly.HTMLElement) {
		handleResultsTable(e, cache)
	})

	c.OnHTML("table.LapTimes", handleTimeTable)
}

func handleCreateUser(e *colly.HTMLElement, c *cache.Driver) {
	href := e.Attr("href")
	id := strings.Split(href, "=")[1]

	var d models.Driver
	if !c.Has(id) { // persist new driver in cache
		d = models.Driver{
			ID:             id,
			Alias:          e.Text,
			ProskillRating: 0, // THIS SHOULD GET UPDATED IN A DIFFERENT HANDLER
		}
		c.Put(d)
	}
	cache.Put(d)

	fmt.Printf("Create User Hit: %s\n\tCustID: %s\n", e.Text, id)
}

func handleResultsTable(e *colly.HTMLElement) {
	fmt.Println("Main Results Table Found")
}

// FIXME: Time Table Handler
func handleTimeTable(e *colly.HTMLElement) {
	fmt.Println("Time Table Found")
}
