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
}

// tbody
func handleResultsTable(e *colly.HTMLElement, c *cache.Cache) {
	e.ForEach("tr[class^='Top3Winners']", func(i int, e *colly.HTMLElement) {
		top3Row(i, e, c)
	})
}

func top3Row(i int, e *colly.HTMLElement, c *cache.Cache) {
	idx := i + 1 // Determines which row parsing strategy to use
	fmt.Printf("i: %d\n", idx)
	switch idx {
	case 1:
		top3RowOne(e, c)
	case 2:
		top3RowTwo(e, c)
	case 3:
		top3RowThree(e, c)
	default:
		// TODO: ERROR
	}

	if idx%3 == 0 {
		// update DB with cache.Race
		// clear cache.Race
	}
}

func top3RowOne(e *colly.HTMLElement, c *cache.Cache) {
}

func top3RowTwo(e *colly.HTMLElement, c *cache.Cache) {
}

func top3RowThree(e *colly.HTMLElement, c *cache.Cache) {
}

// FIXME: Time Table Handler
func handleTimeTable(e *colly.HTMLElement) {
	fmt.Println("Time Table Found")
}
