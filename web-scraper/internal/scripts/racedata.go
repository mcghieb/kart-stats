package scripts

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/mcghieb/kart-stats/web-scraper/internal/cache"
	"github.com/mcghieb/kart-stats/web-scraper/internal/collector"
	"github.com/mcghieb/kart-stats/web-scraper/internal/scrapers/racedata"
)

func ParseRaceData(cache *cache.Cache) {
	fmt.Println("BEGIN")

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

	fmt.Printf("Found %d HTML files to parse\n", len(htmlFiles))

	for _, file := range htmlFiles {
		absPath, err := filepath.Abs(file)
		if err != nil {
			log.Printf("failed to get absolute path for %s: %v", file, err)
			continue
		}
		fileURL := "file://" + absPath
		if err := c.Visit(fileURL); err != nil {
			log.Printf("failed to parse %s: %v", file, err)
		}
	}

	// Print the cached race data for heat 42004
	race := cache.Race.Get("42004")
	fmt.Println("\n=== RACE 42004 ===")
	fmt.Println(race)

	// Print the last 10 drivers from cache
	var driverIDs []string
	cache.Driver.Range(func(id string) bool {
		driverIDs = append(driverIDs, id)
		return true
	})

	fmt.Println("\n=== LAST 10 DRIVERS ===")
	if len(driverIDs) > 0 {
		startIdx := len(driverIDs) - 10
		if startIdx < 0 {
			startIdx = 0
		}
		for i := startIdx; i < len(driverIDs); i++ {
			driver, _ := cache.Driver.Get(driverIDs[i])
			fmt.Println(driver.String())
		}
		fmt.Printf("\nTotal drivers in cache: %d\n", len(driverIDs))
	} else {
		fmt.Println("No drivers in cache")
	}

	fmt.Println("\nEND")

	// TODO: AT THE END OF THE SCRIPT:
	// go through all of the drivers in the cache
	// and update their proskills in the db. This
	// avoids us needing to make a bunch of updates
	// mid process when there could be a lot of
	// separate races of the same drivers
}
