package scripts

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"time"

	"github.com/mcghieb/kart-stats/web-scraper/internal/cache"
	"github.com/mcghieb/kart-stats/web-scraper/internal/collector"
	"github.com/mcghieb/kart-stats/web-scraper/internal/scrapers/racedata"
)

func ParseRacePages(cache *cache.Cache) {
	startTime := time.Now()

	c := collector.NewLocalCollector()
	racedata.AttachHandlers(c, cache)

	// TODO: put path stuff in environment variables?
	// i don't like that the path is hardcoded here
	var htmlFiles []string
	err := filepath.WalkDir("./heatdata-files", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && filepath.Ext(d.Name()) == ".html" {
			htmlFiles = append(htmlFiles, path)
		}

		return nil
	})
	if err != nil {
		panic(err) // not being able to parse html file names is a fatal error
	}

	total := len(htmlFiles)
	fmt.Printf("Found %d HTML files to parse\n", total)

	for i, file := range htmlFiles {
		absPath, err := filepath.Abs(file)
		if err != nil {
			log.Printf("failed to get absolute path for %s: %v", file, err)
			continue
		}
		fileURL := "file://" + absPath
		if err := c.Visit(fileURL); err != nil {
			log.Printf("failed to parse %s: %v", file, err)
		}

		completed := i + 1
		if completed%100 == 0 {
			elapsed := time.Since(startTime)
			rate := float64(completed) / elapsed.Seconds()
			remaining := total - completed
			eta := time.Duration(float64(remaining)/rate) * time.Second

			fmt.Printf("[Progress] %d/%d heats parsed | %.1f heats/sec | Elapsed: %v | ETA: %v\n",
				completed, total, rate, elapsed.Round(time.Second), eta.Round(time.Second))
		}
	}

	elapsed := time.Since(startTime)
	fmt.Println("\n=== PARSING COMPLETE ===")
	fmt.Printf("Parsed %d heats\n", total)
	fmt.Printf("Total time: %v\n", elapsed)
	if elapsed.Seconds() > 0 {
		fmt.Printf("Average rate: %.1f heats/sec\n", float64(total)/elapsed.Seconds())
	}
}
