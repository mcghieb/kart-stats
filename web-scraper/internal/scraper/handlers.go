package scraper

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

// AttachHandlers attaches the HEAT SCRAPING HANDLERS. Name will change if a different handler set is needed for other scraping.
func AttachHandlers(c *colly.Collector) {
	c.OnHTML("a[href^='RacerHistory'", handleCreateUser)

	c.OnHTML("table[class='RaceResults']", handleResultsTable)
	c.OnHTML("table[class='LapTimes']", handleTimeTable)
}

// FIXME: Create User Handler
func handleCreateUser(e *colly.HTMLElement) {
	// Add name to cache to not check again in this scraping process
	// Grab the CustID from the href and check against database to see if this user exists already
	// if not: create a new row in the drivers table
	fmt.Println("Create User Hit")
}

func handleResultsTable(e *colly.HTMLElement) {
	fmt.Println("Main Results Table Found")
}

// FIXME: Time Table Handler
func handleTimeTable(e *colly.HTMLElement) {
	fmt.Println("Time Table Found")
}
