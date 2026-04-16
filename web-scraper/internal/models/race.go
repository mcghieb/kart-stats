package models

import (
	"fmt"
	"strings"
	"time"
)

type Lap struct {
	LapNumber int
	LapTime   float64
}

type Result struct {
	KartIDs                []int  // kart numbers (may be multiple in endurance/swap races)
	DriverID               string // CustID
	Position               int
	Penalties              int
	BestLaptime            float64
	AvgLaptime             float64
	NumLaps                int
	GapFromLeader          float64
	SnapshotProskillRating int // this is the proskill rating after the heat (proskill on driver table is current proskill)     // FIXME:
	Laps                   []Lap
}

// Race represents a race heat
type Race struct {
	ID        string // HeatNo
	Track     string
	WinBy     string
	Timestamp time.Time
	Results   map[string]Result
}

func (r Race) String() string {
	timestampStr := "(not set)"
	if !r.Timestamp.IsZero() {
		timestampStr = r.Timestamp.Format("2006-01-02 3:04 PM")
	}
	
	return fmt.Sprintf(
		"Race Heat Result {\n"+
			"  ID: %s\n"+
			"  Track: %s\n"+
			"  WinBy: %s\n"+
			"  Timestamp: %s\n"+
			"  Results:\n%s"+
			"}",
		r.ID,
		r.Track,
		r.WinBy,
		timestampStr,
		r.makeResultsString(),
	)
}

func (r Race) makeResultsString() string {
	if len(r.Results) == 0 {
		return "    (no results)\n"
	}
	
	var results []string
	for _, v := range r.Results {
		results = append(results, v.String())
	}

	return strings.Join(results, "\n")
}

func (r Result) String() string {
	lapsStr := ""
	if len(r.Laps) > 0 {
		lapsStr = fmt.Sprintf("      Laps: %d recorded\n", len(r.Laps))
	}
	
	return fmt.Sprintf(
		"    Result {\n"+
			"      DriverID: %s\n"+
			"      Position: %d\n"+
			"      KartIDs: %v\n"+
			"      Penalties: %d\n"+
			"      BestLaptime: %.3f\n"+
			"      AvgLaptime: %.3f\n"+
			"      NumLaps: %d\n"+
			"      GapFromLeader: %.3f\n"+
			"      SnapshotProskillRating: %d\n"+
			"%s"+
			"    }",
		r.DriverID,
		r.Position,
		r.KartIDs,
		r.Penalties,
		r.BestLaptime,
		r.AvgLaptime,
		r.NumLaps,
		r.GapFromLeader,
		r.SnapshotProskillRating,
		lapsStr,
	)
}
