package racedata

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/mcghieb/kart-stats/web-scraper/internal/cache"
	"github.com/mcghieb/kart-stats/web-scraper/internal/models"
	"github.com/mcghieb/kart-stats/web-scraper/internal/parse"
)

// AttachHandlers attaches the HEAT SCRAPING HANDLERS.
// Name will change if a different handler set is needed
// for other scraping.
func AttachHandlers(c *colly.Collector, appCache *cache.Cache) {
	driverHandlers(c, appCache)
	raceDataHandlers(c, appCache)
}

func driverHandlers(c *colly.Collector, appCache *cache.Cache) {
	// set kartID in driver's race result
	// early returns on missing context (heatNumber or driverID)
	// this should only run after driver creation step in raceDataHandlers
	c.OnHTML("a[href^='HeatDetails']", func(e *colly.HTMLElement) {
		handleUpdateKartID(e, appCache)
	})

	// set driver's most current Proskill
	// this should only run after driver creation step in raceDataHandlers
	c.OnHTML("span#lblSpeedLimit", func(e *colly.HTMLElement) {
		driverID := e.Request.Ctx.Get("driverID")

		// convert proskill to int
		proskill, err := strconv.Atoi(e.Text)
		if err != nil {
			log.Printf("failed to parse proskill from [%s]", e.Text)
			return
		}

		appCache.Driver.UpdateProskill(driverID, proskill)
	})
}

func raceDataHandlers(c *colly.Collector, appCache *cache.Cache) {
	// Create users first to avoid cache breaks
	c.OnHTML("a[href^='RacerHistory']", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		driverID := strings.Split(href, "=")[1]

		handleCreateUser(e, appCache.Driver)

		h := parse.HeatNum(e)

		// Set context for the driver history page visit
		e.Request.Ctx.Put("heatNumber", h)
		e.Request.Ctx.Put("driverID", driverID)
		e.Request.Visit(href)
	})

	c.OnHTML("span#lblRaceType", func(e *colly.HTMLElement) {
		h := parse.HeatNum(e)
		cache.UpdateCachedRace(appCache, h, func(r *models.Race) {
			r.Track = e.Text
		})
	})

	c.OnHTML("span#lblWinnerBy", func(e *colly.HTMLElement) {
		h := parse.HeatNum(e)
		cache.UpdateCachedRace(appCache, h, func(r *models.Race) {
			r.WinBy = e.Text
		})
	})

	c.OnHTML("span#lblDate", func(e *colly.HTMLElement) {
		h := parse.HeatNum(e)

		dateString := e.Text
		layout := "1/2/2006 3:04 PM"
		date, err := time.Parse(layout, dateString)
		if err != nil {
			log.Printf("failed to parse date for heat %s: %v", h, err)
			return
		}

		cache.UpdateCachedRace(appCache, h, func(r *models.Race) {
			r.Timestamp = date
		})
	})

	c.OnHTML("table.RaceResults > tbody", func(e *colly.HTMLElement) {
		h := parse.HeatNum(e)
		handleResultsTable(e, appCache, h)
	})

	c.OnHTML("table.LapTimes", func(e *colly.HTMLElement) {
		h := parse.HeatNum(e)
		handleTimeTable(e, appCache, h)
	})
}

func handleUpdateKartID(e *colly.HTMLElement, c *cache.Cache) {
	href := e.Attr("href")

	// get heatNumber from context
	heatNumber := e.Request.Ctx.Get("heatNumber")
	// get driverID from context
	driverID := e.Request.Ctx.Get("driverID")
	if heatNumber == "" || driverID == "" {
		log.Printf("failed to fetch heatNumber or driverID from request context when visiting [%s]", href)
		return
	}

	if !strings.Contains(href, heatNumber) { // check if correct element
		// no need to log as this is not an error
		return
	}

	// parse kart id
	re := regexp.MustCompile(`Kart [0-9]+`)
	kartString := re.FindString(e.Text)
	kart := strings.Split(kartString, " ")[1]
	kartID, err := strconv.Atoi(kart)
	if err != nil {
		log.Printf("failed to parse kart number from [%s] for heat %s: %v", kartString, heatNumber, err)
		return
	}

	// cache kart id in driver's race result
	if err := cache.UpdateCachedResultAttribute(e, c, heatNumber, driverID,
		func(e *colly.HTMLElement) (int, error) {
			return kartID, nil
		},
		func(r *models.Result, val int) {
			r.KartID = val
		},
	); err != nil {
		log.Printf("failed to cache kart number for [%s] for heat %s: %v", kartString, heatNumber, err)
		return
	}
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
	var currentDriverID string

	e.ForEach("tr[class^='Top3Winners']", func(i int, e *colly.HTMLElement) {
		// Driver ID is only in the first row of each 3-row group
		if i%3 == 0 {
			var err error
			currentDriverID, err = cache.UpdateCachedDriverID(e, c, heatNum, parse.DriverAlias)
			if err != nil {
				log.Printf("failed to get driver ID for top3 row: %v", err)
				return
			}
		}

		top3Row(i, e, c, heatNum, currentDriverID)
	})

	e.ForEach("tr[class^='RegularRow']", func(i int, e *colly.HTMLElement) {
		driverID, err := cache.UpdateCachedDriverID(e, c, heatNum, parse.DriverAlias)
		if err != nil {
			log.Printf("failed to get driver ID for regular row: %v", err)
			return
		}

		cacheRegularRow(e, c, heatNum, driverID)
	})
}

// TODO: update DB with cache.Race
func top3Row(i int, e *colly.HTMLElement, c *cache.Cache, heatNum, driverID string) {
	idx := (i % 3) + 1 // Determines which row parsing strategy to use
	switch idx {
	case 1:
		cacheTop3RowOne(e, c, heatNum, driverID)
	case 2:
		cacheTop3RowTwo(e, c, heatNum, driverID)
	case 3:
		cacheTop3RowThree(e, c, heatNum, driverID)
	default:
		log.Printf("invalid top3 row index: %d", idx)
	}

	if idx == 3 {
		// Cache printing and clearing moved to OnScraped callback
	}
}

func cacheTop3RowOne(e *colly.HTMLElement, c *cache.Cache, heatNum, driverID string) {
	errs := make([]error, 0)
	if err := cache.UpdateCachedResultAttribute(e, c, heatNum, driverID, parse.Top3Position, func(r *models.Result, v int) { r.Position = v }); err != nil {
		errs = append(errs, err)
	}

	logErrors(errs)
}

func cacheTop3RowTwo(e *colly.HTMLElement, c *cache.Cache, heatNum, driverID string) {
	errs := make([]error, 0)

	if err := cache.UpdateCachedResultAttribute(e, c, heatNum, driverID, parse.BestLaptime,
		func(r *models.Result, v float64) { r.BestLaptime = v },
	); err != nil {
		errs = append(errs, err)
	}

	if err := cache.UpdateCachedResultAttribute(e, c, heatNum, driverID, parse.NumLaps,
		func(r *models.Result, v int) { r.NumLaps = v },
	); err != nil {
		errs = append(errs, err)
	}

	if err := cache.UpdateCachedResultAttribute(e, c, heatNum, driverID, parse.LeaderGap,
		func(r *models.Result, v float64) { r.GapFromLeader = v },
	); err != nil {
		errs = append(errs, err)
	}

	if err := cache.UpdateCachedResultAttribute(e, c, heatNum, driverID, parse.AvgLaptime,
		func(r *models.Result, v float64) { r.AvgLaptime = v },
	); err != nil {
		errs = append(errs, err)
	}

	logErrors(errs)
}

// FIXME: the proskill parsed here is probably behind by one race...
func cacheTop3RowThree(e *colly.HTMLElement, c *cache.Cache, heatNum, driverID string) {
	errs := make([]error, 0)
	if err := cache.UpdateCachedResultAttribute(e, c, heatNum, driverID, parse.Proskill,
		func(r *models.Result, v int) { r.SnapshotProskillRating = v },
	); err != nil {
		errs = append(errs, err)
	}

	logErrors(errs)
}

func cacheRegularRow(e *colly.HTMLElement, c *cache.Cache, heatNum string, driverID string) {
	errs := make([]error, 0)
	if err := cache.UpdateCachedResultAttribute(e, c, heatNum, driverID, parse.Position,
		func(r *models.Result, v int) { r.Position = v },
	); err != nil {
		errs = append(errs, err)
	}

	if err := cache.UpdateCachedResultAttribute(e, c, heatNum, driverID, parse.BestLaptime,
		func(r *models.Result, v float64) { r.BestLaptime = v },
	); err != nil {
		errs = append(errs, err)
	}

	if err := cache.UpdateCachedResultAttribute(e, c, heatNum, driverID, parse.NumLaps,
		func(r *models.Result, v int) { r.NumLaps = v },
	); err != nil {
		errs = append(errs, err)
	}

	if err := cache.UpdateCachedResultAttribute(e, c, heatNum, driverID, parse.LeaderGap,
		func(r *models.Result, v float64) { r.GapFromLeader = v },
	); err != nil {
		errs = append(errs, err)
	}

	if err := cache.UpdateCachedResultAttribute(e, c, heatNum, driverID, parse.AvgLaptime,
		func(r *models.Result, v float64) { r.AvgLaptime = v },
	); err != nil {
		errs = append(errs, err)
	}

	if err := cache.UpdateCachedResultAttribute(e, c, heatNum, driverID, parse.Proskill,
		func(r *models.Result, v int) { r.SnapshotProskillRating = v },
	); err != nil {
		errs = append(errs, err)
	}

	logErrors(errs)
}

func handleTimeTable(e *colly.HTMLElement, c *cache.Cache, heatNum string) {
	driverAlias := e.ChildText("thead > tr > th")

	// Look up driver by alias from cache
	driver, exists := c.Driver.ByAlias(driverAlias)
	if !exists {
		log.Printf("driver %s not found in cache for time table", driverAlias)
		return
	}
	driverID := driver.ID

	rawPen := e.ChildText("tbody > tr:first-child > td")
	re := regexp.MustCompile(`[0-9]+`)
	penString := re.FindString(rawPen)
	numPenalties := 0
	if penString != "" {
		var err error
		numPenalties, err = strconv.Atoi(penString)
		if err != nil {
			log.Printf("failed to parse penalties for heat %s: %v", heatNum, err)
		}
	}

	cache.UpdateCachedRace(c, heatNum, func(r *models.Race) {
		result := r.Results[driverID]
		result.Penalties = numPenalties
		r.Results[driverID] = result
	})

	e.ForEach("tr[class^=LapTimes]", func(i int, lap *colly.HTMLElement) {
		recordLap(c, lap, driverID, heatNum)
	})
}

func recordLap(c *cache.Cache, e *colly.HTMLElement, driverID, heatNum string) {
	cells := e.DOM.Find("td")

	ln := cells.Eq(0).Text()
	lapNum, err := strconv.Atoi(ln)
	if err != nil {
		log.Printf("failed to parse lap number for driver %s: %v", driverID, err)
		return
	}

	rawLaptime := cells.Eq(1).Text()
	re := regexp.MustCompile(`[0-9]+\.[0-9]+`)
	laptimeString := re.FindString(rawLaptime)

	// Skip empty lap cells (drivers who didn't complete all laps)
	if laptimeString == "" {
		return
	}

	laptime, err := strconv.ParseFloat(laptimeString, 64)
	if err != nil {
		log.Printf("failed to parse laptime for driver %s lap %d: %v", driverID, lapNum, err)
		return
	}

	cache.UpdateCachedRaceLaps(c, heatNum, driverID, lapNum, laptime)
}
