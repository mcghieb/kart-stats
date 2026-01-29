package models

type Driver struct {
	ID             string
	Alias          string
	ProskillRating int
	// The field "Name" in the db is not parsed in this scraping script
}
