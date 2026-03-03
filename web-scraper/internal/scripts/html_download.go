package scripts

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/mcghieb/kart-stats/web-scraper/internal/collector"
)

const (
	startHeatno = 1
	endHeatno   = 50_000
	batchSize   = 100
	workers     = 10
)

func numDigits(n int) int {
	if n <= 0 {
		return 1
	}

	return int(math.Log10(float64(n))) + 1
}

func DownloadHTMLpages() {
	fmt.Println("BEGIN HTML DOWNLOAD")
	startTime := time.Now()

	// Counter for completed heats
	var completed atomic.Int64

	// Progress printer goroutine
	var progressWg sync.WaitGroup
	progressWg.Add(1)
	go progressPrinter(&completed, &progressWg, startTime)

	// Create work channel for heat numbers
	heatNums := make(chan int, batchSize)

	// WaitGroup to track worker completion
	var wg sync.WaitGroup

	// Start worker goroutines
	heatDigits := numDigits(endHeatno)
	folderDigits := numDigits(endHeatno / batchSize)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(i, heatNums, &completed, &wg, heatDigits, folderDigits)
	}

	// Send heat numbers to workers
	for i := startHeatno; i <= endHeatno; i++ {
		heatNums <- i
	}
	close(heatNums)

	// Wait for all workers to finish
	wg.Wait()

	// Signal progress printer to stop
	completed.Store(-1)
	progressWg.Wait()

	elapsed := time.Since(startTime)
	fmt.Println("\n=== DOWNLOAD COMPLETE ===")
	fmt.Printf("Downloaded heats from %d to %d\n", startHeatno, endHeatno)
	fmt.Printf("Total time: %v\n", elapsed)
	fmt.Printf("Average rate: %.1f heats/sec\n", float64(endHeatno-startHeatno+1)/elapsed.Seconds())
	fmt.Println("\nEND")
}

func worker(id int, heatNums <-chan int, completed *atomic.Int64, wg *sync.WaitGroup, heatDigits, folderDigits int) {
	defer wg.Done()

	// Each worker gets its own collector
	c := collector.NewCollector()
	attachHandler(c, heatDigits, folderDigits)

	for heatNo := range heatNums {
		url := fmt.Sprintf("https://rrorem.clubspeed.com/sp_center/HeatDetails.aspx?HeatNo=%d", heatNo)
		if err := c.Visit(url); err != nil {
			// Silently skip errors to reduce noise
		}
		completed.Add(1)
	}
}

func attachHandler(c *colly.Collector, heatDigits, folderDigits int) {
	folderFmt := fmt.Sprintf("heatdata-files/batch-%%0%dd", folderDigits)
	fileFmt := fmt.Sprintf("%s/%%0%dd.html", folderFmt, heatDigits)

	c.OnHTML("body", func(e *colly.HTMLElement) {
		if strings.Contains(e.Request.URL.Path, "/sp_center/ServerError.html") {
			// Silently skip errors to reduce noise
			return
		}

		// Skip races with no winner (winner is "-")
		winner := e.ChildText("span#lblWinner")
		if strings.TrimSpace(winner) == "-" {
			return
		}

		heatno := e.Request.URL.Query().Get("HeatNo")
		heat, err := strconv.Atoi(heatno)
		if err != nil {
			log.Printf("failed to convert heatno [%s] to int: %v", heatno, err)
		}

		folder := heat / batchSize
		path := fmt.Sprintf(folderFmt, folder)
		if !folderExists(path) {
			err := os.MkdirAll(path, 0755)
			if err != nil {
				panic(err)
			}
		}

		filename := fmt.Sprintf(fileFmt, folder, heat)
		if err := e.Response.Save(filename); err != nil {
			log.Printf("failed to save %s: %v", filename, err)
			return
		}
	})
}

func folderExists(path string) bool {
	info, err := os.Stat(path)
	return !os.IsNotExist(err) && info.IsDir()
}

func progressPrinter(completed *atomic.Int64, wg *sync.WaitGroup, startTime time.Time) {
	defer wg.Done()
	lastReported := int64(0)

	for {
		time.Sleep(500 * time.Millisecond)
		current := completed.Load()

		// Check if we should exit
		if current == -1 {
			return
		}

		// Report every 100 heats
		if current >= lastReported+100 {
			newMilestone := (current / 100) * 100
			elapsed := time.Since(startTime)
			rate := float64(current) / elapsed.Seconds()
			remaining := endHeatno - int(current)
			eta := time.Duration(float64(remaining)/rate) * time.Second

			fmt.Printf("[Progress] %d heats completed | %.1f heats/sec | Elapsed: %v | ETA: %v\n",
				newMilestone, rate, elapsed.Round(time.Second), eta.Round(time.Second))
			lastReported = newMilestone
		}
	}
}
