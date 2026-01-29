package scraper

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/mcghieb/kart-stats/web-scraper/internal/models"
	"github.com/mcghieb/kart-stats/web-scraper/internal/scraper/cache"
)

// AttachHandlers attaches the HEAT SCRAPING HANDLERS.
// Name will change if a different handler set is needed
// for other scraping.
func AttachHandlers(c *colly.Collector, cache *cache.Cache) {
	c.OnHTML("a[href^='RacerHistory']", func(e *colly.HTMLElement) {
		handleCreateUser(e, cache.Driver)
	})

	c.OnHTML("table.RaceResults > tbody", func(e *colly.HTMLElement) {
		handleResultsTable(e, cache)
	})

	// FIXME:
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

func handleResultsTable(e *colly.HTMLElement, c *cache.Cache) {
	// FIXME: find a way to add heatnumber to cached result

	e.ForEach("tr[class^='Top3Winners']", func(i int, e *colly.HTMLElement) {
		top3Row(i, e, c)
	})

	e.ForEach("tr[class^='RegularRow']", func(i int, e *colly.HTMLElement) {
		regularRow(i, e, c)
	})
}

func top3Row(i int, e *colly.HTMLElement, c *cache.Cache) {
	idx := (i % 3) + 1 // Determines which row parsing strategy to use
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

	if idx == 3 {
		// TODO: update DB with cache.Race

		fmt.Println(c.Race.Get(getHeatNum(e)))
		// clear cache.Race
		c.Race = cache.NewRaceCache()
	}
}

func top3RowOne(e *colly.HTMLElement, c *cache.Cache) {
	// set position
	p := e.ChildText("td.Position")
	position := 0
	switch p {
	case "Heat Winner:":
		position = 1
	case "2nd Place:":
		position = 2
	case "3rd Place:":
		position = 3
	}
	updateResultInCache(e, c, func(r *models.RaceResult) {
		r.Position = position
	})

	if err := setDriverID(e, c); err != nil {
		// TODO: error handling
	}
}

func top3RowTwo(e *colly.HTMLElement, c *cache.Cache) {
	if err := setBestLaptime(e, c); err != nil {
		// TODO: error handling
	}

	if err := setNumLaps(e, c); err != nil {
		// TODO: error handling
	}

	if err := setLeaderGap(e, c); err != nil {
		// TODO: error handling
	}

	if err := setAvgLaptime(e, c); err != nil {
		// TODO: error handling
	}
}

func top3RowThree(e *colly.HTMLElement, c *cache.Cache) {
	if err := setProskill(e, c); err != nil {
		// TODO: error
	}
}

func regularRow(i int, e *colly.HTMLElement, c *cache.Cache) {
	if err := setDriverID(e, c); err != nil {
		// TODO: error handling
	}

	if err := setPosition(e, c); err != nil {
		// TODO: error handling
	}

	if err := setBestLaptime(e, c); err != nil {
		// TODO: error handling
	}

	if err := setNumLaps(e, c); err != nil {
		// TODO: error handling
	}

	if err := setLeaderGap(e, c); err != nil {
		// TODO: error handling
	}

	if err := setAvgLaptime(e, c); err != nil {
		// TODO: error handling
	}

	if err := setProskill(e, c); err != nil {
		// TODO: error
	}

	// Print completed race result
	fmt.Println(c.Race.Get(getHeatNum(e)))

	// FIXME: persist to db and clear cache.Race
}

// FIXME: Time Table Handler
func handleTimeTable(e *colly.HTMLElement) {
	fmt.Println("Time Table Found")
}
