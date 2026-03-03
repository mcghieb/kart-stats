package scripts

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/mcghieb/kart-stats/web-scraper/internal/cache"
	"github.com/mcghieb/kart-stats/web-scraper/internal/collector"
	"github.com/mcghieb/kart-stats/web-scraper/internal/scrapers/driverhistory"
)

const (
	BASEURL       = "https://rrorem.clubspeed.com/sp_center/RacerHistory.aspx?CustID"
	driverWorkers = 10
)

type driverVisit struct {
	url      string
	driverID string
}

func ParseDriverPages(appCache *cache.Cache) {
	startTime := time.Now()

	// Build list of URLs from all drivers in cache
	var visits []driverVisit
	appCache.Driver.Range(func(id string) bool {
		visits = append(visits, driverVisit{
			url:      createURL(id),
			driverID: id,
		})
		return true
	})

	total := len(visits)

	// Counter for completed pages
	var completed atomic.Int64

	// Progress printer goroutine
	var progressWg sync.WaitGroup
	progressWg.Add(1)
	go progressPrinter(&completed, &progressWg, startTime, total)

	// Create work channel
	work := make(chan driverVisit, batchSize)

	// Start worker goroutines
	var wg sync.WaitGroup
	for range driverWorkers {
		wg.Add(1)
		go driverWorker(work, &completed, &wg, appCache)
	}

	// Send visits to workers
	for _, v := range visits {
		work <- v
	}
	close(work)

	// Wait for all workers to finish
	wg.Wait()

	// Signal progress printer to stop
	completed.Store(-1)
	progressWg.Wait()

	elapsed := time.Since(startTime)
	fmt.Println("\n=== DRIVER HISTORY PARSING COMPLETE ===")
	fmt.Printf("Parsed %d driver history pages\n", total)
	fmt.Printf("Total time: %v\n", elapsed)
	if elapsed.Seconds() > 0 {
		fmt.Printf("Average rate: %.1f pages/sec\n", float64(total)/elapsed.Seconds())
	}
}

func driverWorker(work <-chan driverVisit, completed *atomic.Int64, wg *sync.WaitGroup, appCache *cache.Cache) {
	defer wg.Done()

	// Each worker gets its own collector
	c := collector.NewCollector()
	driverhistory.AttachHandlers(c, appCache)

	for v := range work {
		ctx := colly.NewContext()
		ctx.Put("driverID", v.driverID)

		if err := c.Request("GET", v.url, nil, ctx, nil); err != nil {
			log.Printf("failed to visit driver page for %s: %v", v.driverID, err)
		}
		completed.Add(1)
	}
}

func createURL(driverID string) string {
	return fmt.Sprintf("%s=%s==", BASEURL, driverID)
}
