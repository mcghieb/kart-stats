package driverhistory

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/mcghieb/kart-stats/web-scraper/internal/cache"
	"github.com/mcghieb/kart-stats/web-scraper/internal/models"
	"github.com/mcghieb/kart-stats/web-scraper/internal/parse"
)

// AttachHandlers attaches handlers for scraping driver history pages
// This scraper is used to get current proskill ratings for all drivers
func AttachHandlers(c *colly.Collector, store *cache.Cache) {
	// Parse current proskill rating
	c.OnHTML("span#lblSpeedLimit", func(e *colly.HTMLElement) {
		dID := e.Request.Ctx.Get("driverID")
		if dID == "" {
			log.Printf("missing driverID in context")
			return
		}

		fmt.Printf("PROSKILL for %s is %s\n", dID, e.Text)
		proskill, err := parseProskill(e.Text)
		if err != nil {
			log.Printf("failed to parse proskill for driver %s: %v", dID, err)
			return
		}

		if err := store.Driver.UpdateProskill(dID, proskill); err != nil {
			log.Printf("failed to store proskill for driver %s: %v", dID, err)
		}
	})

	// Update kart IDs for each heat in the cache from the driver history table
	c.OnHTML("table#dg", func(e *colly.HTMLElement) {
		dID := e.Request.Ctx.Get("driverID")
		if dID == "" {
			log.Printf("missing driverID in context")
			return
		}

		done := false
		e.ForEach("tr.Normal", func(_ int, row *colly.HTMLElement) {
			if done {
				return
			}

			link := row.DOM.Find("a[href*='HeatNo']")
			if link.Length() == 0 {
				return
			}

			href, exists := link.Attr("href")
			if !exists {
				return
			}

			heatNum := parse.HeatNumFromURL(href)

			if !store.Race.Has(heatNum) {
				done = true
				return
			}

			activityText := strings.TrimSpace(link.Text())
			kartID, err := parse.KartID(activityText)
			if err != nil {
				log.Printf("failed to parse kart ID from %q for heat %s: %v", activityText, heatNum, err)
				return
			}

			cache.UpdateCachedRace(store, heatNum, func(r *models.Race) {
				result := r.Results[dID]
				result.KartID = kartID
				r.Results[dID] = result
			})
		})
	})
}

func parseProskill(text string) (int, error) {
	proskill, err := strconv.Atoi(text)
	if err != nil {
		return 0, fmt.Errorf("failed to parse proskill from [%s]: %w", text, err)
	}
	return proskill, nil
}
