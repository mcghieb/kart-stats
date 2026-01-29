package models

import "fmt"

type RaceResult struct {
	ID             string // HeatNo //FIXME: this is missing
	KartID         int    // kart number
	DriverID       string // CustID
	Position       int
	Penalties      int // FIXME: also missing
	BestLaptime    float64
	AvgLaptime     float64
	NumLaps        int
	GapFromLeader  float64
	ProskillRating int // FIXME: add to db schema
	// FIXME: missing: timestamp, track name
}

func (rr RaceResult) String() string {
	return fmt.Sprintf(
		"RaceResult{\n"+
			"\tID: %s\n"+
			"\tKartID: %d\n"+
			"\tDriverID: %s\n"+
			"\tPosition: %d\n"+
			"\tPenalties: %d\n"+
			"\tBestLaptime: %.3f\n"+
			"\tAvgLaptime: %.3f\n"+
			"\tNumLaps: %d\n"+
			"\tGapFromLeader: %.3f\n"+
			"\tProskillRating: %d\n"+
			"}",
		rr.ID,
		rr.KartID,
		rr.DriverID,
		rr.Position,
		rr.Penalties,
		rr.BestLaptime,
		rr.AvgLaptime,
		rr.NumLaps,
		rr.GapFromLeader,
		rr.ProskillRating,
	)
}
