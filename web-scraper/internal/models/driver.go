package models

import "fmt"

type Driver struct {
	ID             string
	Alias          string
	ProskillRating int
	// The field "Name" in the db is not parsed in this scraping script
}

func (d Driver) String() string {
	return fmt.Sprintf("Driver { ID: %s, Alias: %s, ProskillRating: %d }",
		d.ID, d.Alias, d.ProskillRating)
}
