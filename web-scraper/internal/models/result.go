package models

type RaceResult struct {
	ID            string // HeatNo
	KartID        int    // kart number
	DriverID      string // CustID
	Position      int
	Penalties     int
	BestLaptime   float32
	AvgLaptime    float32
	NumLaps       int
	GapFromLeader float32
}
