package scraper

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/mcghieb/kart-stats/web-scraper/internal/models"
)

// AttachHandlers attaches the HEAT SCRAPING HANDLERS. Name will change if a different handler set is needed for other scraping.
func AttachHandlers(c *colly.Collector, cache *DriverCache) {
	c.OnHTML("a[href^='RacerHistory']", func(e *colly.HTMLElement) {
		handleCreateUser(e, cache)
	})
	c.OnHTML("table.RaceResults > tbody", handleResultsTable)
	c.OnHTML("table.LapTimes", handleTimeTable)
}

// FIXME: Create User Handler
func handleCreateUser(e *colly.HTMLElement, cache *DriverCache) {
	// Add name to cache to not check again in this scraping process
	// Grab the CustID from the href and check against database to see if this user exists already
	// if not: create a new row in the drivers table

	href := e.Attr("href")
	id := strings.Split(href, "=")[1]

	var d models.Driver
	if !cache.Has(id) { // persist new driver in cache
		// TODO: create driver
		d = models.Driver{
			ID:             id,
			Name:           "", // FIXME:
			Alias:          "", // FIXME:
			ProskillRating: 0,  // FIXME:
		}
	} else { // modify driver's proskill rating in the cache
		// FIXME:
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
