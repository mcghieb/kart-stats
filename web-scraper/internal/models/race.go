package models

import (
	"fmt"
	"time"
)

// Race represents a race heat result for a driver
type Race struct {
	ID                     string // HeatNo
	KartID                 int    // kart number
	DriverID               string // CustID
	Position               int
	Penalties              int
	BestLaptime            float64
	AvgLaptime             float64
	NumLaps                int
	GapFromLeader          float64
	SnapshotProskillRating int       // this is the proskill rating after the heat (proskill on driver table is current proskill)     // FIXME:
	Track                  string    // FIXME:
	Timestamp              time.Time // FIXME:
}

func (r Race) String() string {
	return fmt.Sprintf(
		"Race Heat Result{\n"+
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
		r.ID,
		r.KartID,
		r.DriverID,
		r.Position,
		r.Penalties,
		r.BestLaptime,
		r.AvgLaptime,
		r.NumLaps,
		r.GapFromLeader,
		r.SnapshotProskillRating,
	)
}
