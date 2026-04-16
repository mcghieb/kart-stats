package repository

import _ "embed"

//go:embed schema.sql
var Schema string

type Repository interface {
	Migrate() error
	Close()
	BatchRaceUpload(rows [][]any, kartRows [][]any) error
	BatchDriverUpload(rows [][]any) error
}
