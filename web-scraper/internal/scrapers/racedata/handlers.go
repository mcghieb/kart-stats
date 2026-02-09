package racedata

import (
	"fmt"
	"regexp"
	"strconv"
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
	c.OnHTML("table.LapTimes", func(e *colly.HTMLElement) {
		h := parse.HeatNum(e)
		handleTimeTable(e, cache, h)
	})
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

// TODO: update DB with cache.Race
func top3Row(i int, e *colly.HTMLElement, c *cache.Cache, heatNum string) {
	idx := (i % 3) + 1 // Determines which row parsing strategy to use
	switch idx {
	case 1:
		cacheTop3RowOne(e, c, heatNum)
	case 2:
		cacheTop3RowTwo(e, c, heatNum)
	case 3:
		cacheTop3RowThree(e, c, heatNum)
	default:
		fmt.Println(fmt.Errorf("ERROR: top3Row route does not exist"))
	}

	if idx == 3 {
		// TODO: update DB with cache.Race

		fmt.Println(c.Race.Get(parse.HeatNum(e))) // TODO: remove

		// clear cache.Race
		c.Race = cache.NewRaceCache()
	}
}

func cacheTop3RowOne(e *colly.HTMLElement, c *cache.Cache, heatNum string) {
	errs := make([]error, 0)
	if err := cache.UpdateCachedDriverID(e, c, heatNum, parse.DriverAlias); err != nil {
		errs = append(errs, err)
	}
	if err := cache.UpdateCachedRaceAttribute(e, c, heatNum, parse.Top3Position, func(r *models.Race, v int) { r.Position = v }); err != nil {
		errs = append(errs, err)
	}

	logErrors(errs)
}

func cacheTop3RowTwo(e *colly.HTMLElement, c *cache.Cache, heatNum string) {
	errs := make([]error, 0)
	if err := cache.UpdateCachedRaceAttribute(e, c, heatNum, parse.BestLaptime, func(r *models.Race, v float64) { r.BestLaptime = v }); err != nil {
		errs = append(errs, err)
	}
	if err := cache.UpdateCachedRaceAttribute(e, c, heatNum, parse.NumLaps, func(r *models.Race, v int) { r.NumLaps = v }); err != nil {
		errs = append(errs, err)
	}
	if err := cache.UpdateCachedRaceAttribute(e, c, heatNum, parse.LeaderGap, func(r *models.Race, v float64) { r.GapFromLeader = v }); err != nil {
		errs = append(errs, err)
	}
	if err := cache.UpdateCachedRaceAttribute(e, c, heatNum, parse.AvgLaptime, func(r *models.Race, v float64) { r.AvgLaptime = v }); err != nil {
		errs = append(errs, err)
	}

	logErrors(errs)
}

func cacheTop3RowThree(e *colly.HTMLElement, c *cache.Cache, heatNum string) {
	errs := make([]error, 0)
	if err := cache.UpdateCachedRaceAttribute(e, c, heatNum, parse.Proskill, func(r *models.Race, v int) { r.SnapshotProskillRating = v }); err != nil {
		errs = append(errs, err)
	}

	// update proskill driver cache
	proskill := c.Race.Get(heatNum).SnapshotProskillRating
	driverID := c.Race.Get(heatNum).DriverID
	if err := c.Driver.UpdateProskill(driverID, proskill); err != nil {
		errs = append(errs, err)
	}

	logErrors(errs)
}

// FIXME: persist to db and clear cache.Race
func regularRow(e *colly.HTMLElement, c *cache.Cache, heatNum string) {
	cacheRegularRow(e, c, heatNum)

	// Print completed race result
	fmt.Println(c.Race.Get(parse.HeatNum(e))) // TODO: remove this

	// FIXME: persist to db and clear cache.Race
	c.Race = cache.NewRaceCache()
}

func cacheRegularRow(e *colly.HTMLElement, c *cache.Cache, heatNum string) {
	errs := make([]error, 0)
	if err := cache.UpdateCachedDriverID(e, c, heatNum, parse.DriverAlias); err != nil {
		errs = append(errs, err)
	}
	if err := cache.UpdateCachedRaceAttribute(e, c, heatNum, parse.Position, func(r *models.Race, v int) { r.Position = v }); err != nil {
		errs = append(errs, err)
	}
	if err := cache.UpdateCachedRaceAttribute(e, c, heatNum, parse.BestLaptime, func(r *models.Race, v float64) { r.BestLaptime = v }); err != nil {
		errs = append(errs, err)
	}
	if err := cache.UpdateCachedRaceAttribute(e, c, heatNum, parse.NumLaps, func(r *models.Race, v int) { r.NumLaps = v }); err != nil {
		errs = append(errs, err)
	}
	if err := cache.UpdateCachedRaceAttribute(e, c, heatNum, parse.LeaderGap, func(r *models.Race, v float64) { r.GapFromLeader = v }); err != nil {
		errs = append(errs, err)
	}
	if err := cache.UpdateCachedRaceAttribute(e, c, heatNum, parse.AvgLaptime, func(r *models.Race, v float64) { r.AvgLaptime = v }); err != nil {
		errs = append(errs, err)
	}
	if err := cache.UpdateCachedRaceAttribute(e, c, heatNum, parse.Proskill, func(r *models.Race, v int) { r.SnapshotProskillRating = v }); err != nil {
		errs = append(errs, err)
	}
	if err := cache.UpdateCachedRaceAttribute(e, c, heatNum, parse.NumLaps, func(r *models.Race, v int) { r.NumLaps = v }); err != nil {
		errs = append(errs, err)
	}
	if err := cache.UpdateCachedRaceAttribute(e, c, heatNum, parse.NumLaps, func(r *models.Race, v int) { r.NumLaps = v }); err != nil {
		errs = append(errs, err)
	}
	logErrors(errs)
}

func handleTimeTable(e *colly.HTMLElement, c *cache.Cache, heatNum string) {
	driverAlias := e.ChildText("thead > tr > th")
	fmt.Printf("Time Table Found: %s\n", driverAlias)

	rawPen := e.ChildText("tbody > tr:first-child > td")
	re := regexp.MustCompile(`[0-9]*`)
	penString := re.FindString(rawPen)
	numPenalties, err := strconv.Atoi(penString)
	if err != nil {
		println("ERROR: parsing int from string in handleTimeTable")
	}

	cache.UpdateCachedRace(c, heatNum, func(r *models.Race) {
		r.Penalties = numPenalties
	})

	e.ForEach("tr[class^=LapTimes]", func(i int, lap *colly.HTMLElement) {
		recordLap(driverAlias, lap)
	})
}

func recordLap(driverAlias string, lap *colly.HTMLElement) {
	cells := lap.DOM.Find("td")

	lapNumber := cells.Eq(0).Text()
	rawLaptime := cells.Eq(1).Text()
	re := regexp.MustCompile(`[0-9]*\.[0-9]*`)
	laptimeString := re.FindString(rawLaptime)

	fmt.Printf("\t%s's lap %s: %s\n", driverAlias, lapNumber, laptimeString)
}

func logErrors(errs []error) {
	for _, err := range errs {
		fmt.Println(fmt.Errorf("ERROR: %w", err))
	}
}
