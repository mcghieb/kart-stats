package racedata

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/mcghieb/kart-stats/web-scraper/internal/cache"
	"github.com/mcghieb/kart-stats/web-scraper/internal/models"
	"github.com/mcghieb/kart-stats/web-scraper/internal/parse"
)

// AttachHandlers attaches the HEAT SCRAPING HANDLERS.
// Name will change if a different handler set is needed
// for other scraping.
func AttachHandlers(c *colly.Collector, cache *cache.Cache) {
	c.OnHTML("a[href^='RacerHistory']", func(e *colly.HTMLElement) {
		handleCreateUser(e, cache.Driver)
	})

	c.OnHTML("table.RaceResults > tbody", func(e *colly.HTMLElement) {
		h := parse.HeatNum(e)
		handleResultsTable(e, cache, h)
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
			ProskillRating: 0, // THIS GETS UPDATED IN A DIFFERENT HANDLER
		}
		c.Put(d)
	}
}

func handleResultsTable(e *colly.HTMLElement, c *cache.Cache, heatNum string) {
	e.ForEach("tr[class^='Top3Winners']", func(i int, e *colly.HTMLElement) {
		top3Row(i, e, c, heatNum)
	})

	e.ForEach("tr[class^='RegularRow']", func(i int, e *colly.HTMLElement) {
		regularRow(e, c, heatNum)
	})
}

func top3Row(i int, e *colly.HTMLElement, c *cache.Cache, heatNum string) {
	idx := (i % 3) + 1 // Determines which row parsing strategy to use
	switch idx {
	case 1:
		top3RowOne(e, c, heatNum)
	case 2:
		top3RowTwo(e, c, heatNum)
	case 3:
		top3RowThree(e, c, heatNum)
	default:
		// TODO: ERROR
	}

	if idx == 3 {
		// TODO: update DB with cache.Race

		fmt.Println(c.Race.Get(parse.HeatNum(e))) // TODO: remove

		// clear cache.Race
		c.Race = cache.NewRaceCache()
	}
}

// FIXME: swallowing error handling
func top3RowOne(e *colly.HTMLElement, c *cache.Cache, heatNum string) {
	cache.UpdateCachedDriverID(e, c, heatNum, parse.DriverAlias)
	cache.UpdateCachedRaceAttribute(e, c, heatNum, parse.Top3Position, func(r *models.Race, v int) { r.Position = v })
}

// FIXME: swallowing error handling
func top3RowTwo(e *colly.HTMLElement, c *cache.Cache, heatNum string) {
	cache.UpdateCachedRaceAttribute(e, c, heatNum, parse.BestLaptime, func(r *models.Race, v float64) { r.BestLaptime = v })
	cache.UpdateCachedRaceAttribute(e, c, heatNum, parse.NumLaps, func(r *models.Race, v int) { r.NumLaps = v })
	cache.UpdateCachedRaceAttribute(e, c, heatNum, parse.LeaderGap, func(r *models.Race, v float64) { r.GapFromLeader = v })
	cache.UpdateCachedRaceAttribute(e, c, heatNum, parse.AvgLaptime, func(r *models.Race, v float64) { r.AvgLaptime = v })
}

// FIXME: swallowing error handling
func top3RowThree(e *colly.HTMLElement, c *cache.Cache, heatNum string) {
	cache.UpdateCachedRaceAttribute(e, c, heatNum, parse.Proskill, func(r *models.Race, v int) { r.ProskillRating = v })

	// update proskill driver cache
	proskill := c.Race.Get(heatNum).ProskillRating
	driverID := c.Race.Get(heatNum).DriverID
	if err := c.Driver.UpdateProskill(driverID, proskill); err != nil {
		// TODO: error
	}
}

// FIXME: swallowing error handling
func regularRow(e *colly.HTMLElement, c *cache.Cache, heatNum string) {
	cache.UpdateCachedDriverID(e, c, heatNum, parse.DriverAlias)
	cache.UpdateCachedRaceAttribute(e, c, heatNum, parse.Position, func(r *models.Race, v int) { r.Position = v })
	cache.UpdateCachedRaceAttribute(e, c, heatNum, parse.BestLaptime, func(r *models.Race, v float64) { r.BestLaptime = v })
	cache.UpdateCachedRaceAttribute(e, c, heatNum, parse.NumLaps, func(r *models.Race, v int) { r.NumLaps = v })
	cache.UpdateCachedRaceAttribute(e, c, heatNum, parse.LeaderGap, func(r *models.Race, v float64) { r.GapFromLeader = v })
	cache.UpdateCachedRaceAttribute(e, c, heatNum, parse.AvgLaptime, func(r *models.Race, v float64) { r.AvgLaptime = v })
	cache.UpdateCachedRaceAttribute(e, c, heatNum, parse.Proskill, func(r *models.Race, v int) { r.ProskillRating = v })
	cache.UpdateCachedRaceAttribute(e, c, heatNum, parse.NumLaps, func(r *models.Race, v int) { r.NumLaps = v })
	cache.UpdateCachedRaceAttribute(e, c, heatNum, parse.NumLaps, func(r *models.Race, v int) { r.NumLaps = v })

	// Print completed race result TODO: remove this
	fmt.Println(c.Race.Get(parse.HeatNum(e)))

	// FIXME: persist to db and clear cache.Race
}

// FIXME: Time Table Handler
func handleTimeTable(e *colly.HTMLElement) {
	fmt.Println("Time Table Found")
}
