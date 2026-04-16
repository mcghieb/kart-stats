package providers

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mcghieb/kart-stats/web-scraper/internal/repository"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(connStr string) (*PostgresStore, error) {
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate database connection pool: %w", err)
	}

	return &PostgresStore{pool: pool}, nil
}

// Close closes the connection pool.
func (r *PostgresStore) Close() {
	r.pool.Close()
}

// Migrate runs the schema setup script. Safe to call multiple times.
func (r *PostgresStore) Migrate() error {
	if _, err := r.pool.Exec(context.Background(), repository.Schema); err != nil {
		return fmt.Errorf("failed to migrate sql schema: %w", err)
	}
	return nil
}

// BatchRaceUpload inserts race data into the races, karts, race_results, race_result_karts, and laps tables.
//
// Each row is a denormalized entry representing one lap:
//
//	[race_id(int), race_time(time.Time), driver_id(string),
//	 position(int), penalties(int), best_laptime(float64), avg_laptime(float64),
//	 num_laps(int), gap_from_leader(float64), lap_time(float64)]
//
// Each kartRow maps a driver to a kart used in a race:
//
//	[race_id(int), driver_id(string), kart_id(int)]
func (r *PostgresStore) BatchRaceUpload(rows [][]any, kartRows [][]any) error {
	if len(rows) == 0 {
		return nil
	}

	ctx := context.Background()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	races, results, laps := splitRaceRows(rows)
	resultKarts, karts := deduplicateKartRows(kartRows)

	// Insert races and karts via batch exec (supports ON CONFLICT).
	batch := &pgx.Batch{}
	for _, race := range races {
		batch.Queue(
			"INSERT INTO races (id, time) VALUES ($1, $2) ON CONFLICT DO NOTHING",
			race[0], race[1],
		)
	}
	for _, kartID := range karts {
		batch.Queue(
			"INSERT INTO karts (id) VALUES ($1) ON CONFLICT DO NOTHING",
			kartID,
		)
	}
	if err := tx.SendBatch(ctx, batch).Close(); err != nil {
		return fmt.Errorf("failed to insert races/karts: %w", err)
	}

	// COPY race_results.
	_, err = tx.CopyFrom(ctx,
		pgx.Identifier{"race_results"},
		[]string{"race_id", "driver_id", "position", "penalties",
			"best_laptime", "avg_laptime", "num_laps", "gap_from_leader"},
		pgx.CopyFromRows(results),
	)
	if err != nil {
		return fmt.Errorf("failed to copy race results: %w", err)
	}

	// COPY race_result_karts.
	_, err = tx.CopyFrom(ctx,
		pgx.Identifier{"race_result_karts"},
		[]string{"race_id", "driver_id", "kart_id"},
		pgx.CopyFromRows(resultKarts),
	)
	if err != nil {
		return fmt.Errorf("failed to copy race result karts: %w", err)
	}

	// COPY laps.
	_, err = tx.CopyFrom(ctx,
		pgx.Identifier{"laps"},
		[]string{"race_id", "driver_id", "lap_time"},
		pgx.CopyFromRows(laps),
	)
	if err != nil {
		return fmt.Errorf("failed to copy laps: %w", err)
	}

	return tx.Commit(ctx)
}

// BatchDriverUpload inserts drivers into the drivers table.
//
// Each row: [id(string), alias(string), proskill_rating(int)]
func (r *PostgresStore) BatchDriverUpload(rows [][]any) error {
	if len(rows) == 0 {
		return nil
	}

	_, err := r.pool.CopyFrom(
		context.Background(),
		pgx.Identifier{"drivers"},
		[]string{"id", "alias", "proskill_rating"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return fmt.Errorf("failed to copy drivers: %w", err)
	}
	return nil
}

// splitRaceRows deduplicates denormalized race rows into separate table data.
func splitRaceRows(rows [][]any) (races [][]any, results [][]any, laps [][]any) {
	seenRaces := make(map[any]bool)

	type resultKey struct{ raceID, driverID any }
	seenResults := make(map[resultKey]bool)

	for _, row := range rows {
		raceID, raceTime := row[0], row[1]
		driverID := row[2]

		if !seenRaces[raceID] {
			races = append(races, []any{raceID, raceTime})
			seenRaces[raceID] = true
		}

		key := resultKey{raceID, driverID}
		if !seenResults[key] {
			results = append(results, []any{
				raceID, driverID,
				row[3], row[4], row[5], row[6], row[7], row[8],
			})
			seenResults[key] = true
		}

		laps = append(laps, []any{raceID, driverID, row[9]})
	}

	return
}

// deduplicateKartRows deduplicates kart junction rows and collects unique kart IDs.
func deduplicateKartRows(kartRows [][]any) (deduped [][]any, uniqueKarts []any) {
	type kartKey struct{ raceID, driverID, kartID any }
	seen := make(map[kartKey]bool)
	seenKarts := make(map[any]bool)

	for _, row := range kartRows {
		kartID := row[2]

		if !seenKarts[kartID] {
			uniqueKarts = append(uniqueKarts, kartID)
			seenKarts[kartID] = true
		}

		key := kartKey{row[0], row[1], kartID}
		if !seen[key] {
			deduped = append(deduped, row)
			seen[key] = true
		}
	}

	return
}
