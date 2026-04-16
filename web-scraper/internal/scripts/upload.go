package scripts

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/mcghieb/kart-stats/web-scraper/internal/cache"
	"github.com/mcghieb/kart-stats/web-scraper/internal/repository"
)

func UploadToDatabase(c *cache.Cache, repo repository.Repository) error {
	startTime := time.Now()

	// Build driver rows: [id, alias, proskill_rating]
	var driverRows [][]any
	c.Driver.Range(func(id string) bool {
		d, ok := c.Driver.Get(id)
		if ok {
			driverRows = append(driverRows, []any{d.ID, d.Alias, d.ProskillRating})
		}
		return true
	})

	fmt.Printf("Uploading %d drivers...\n", len(driverRows))
	if err := repo.BatchDriverUpload(driverRows); err != nil {
		return fmt.Errorf("driver upload: %w", err)
	}

	// Build denormalized race rows (one row per lap):
	// [race_id, race_time, driver_id, position, penalties,
	//  best_laptime, avg_laptime, num_laps, gap_from_leader, lap_time]
	// Plus kart junction rows: [race_id, driver_id, kart_id]
	var raceRows [][]any
	var kartRows [][]any
	var skipped int
	c.Race.Range(func(raceID string) bool {
		race := c.Race.Get(raceID)

		id, err := strconv.Atoi(race.ID)
		if err != nil {
			log.Printf("skipping race with non-numeric ID %q: %v", race.ID, err)
			skipped++
			return true
		}

		for _, result := range race.Results {
			if len(result.Laps) == 0 {
				continue
			}

			for _, kartID := range result.KartIDs {
				kartRows = append(kartRows, []any{id, result.DriverID, kartID})
			}

			for _, lap := range result.Laps {
				raceRows = append(raceRows, []any{
					id, race.Timestamp,
					result.DriverID,
					result.Position, result.Penalties,
					result.BestLaptime, result.AvgLaptime,
					result.NumLaps, result.GapFromLeader,
					lap.LapTime,
				})
			}
		}
		return true
	})

	fmt.Printf("Uploading %d lap rows, %d kart mappings (skipped %d races)...\n", len(raceRows), len(kartRows), skipped)
	if err := repo.BatchRaceUpload(raceRows, kartRows); err != nil {
		return fmt.Errorf("race upload: %w", err)
	}

	fmt.Printf("Upload complete in %v\n", time.Since(startTime).Round(time.Millisecond))
	return nil
}
