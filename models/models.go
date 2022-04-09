package models

import (
	"qa/models/pgmodels"

	"github.com/jmoiron/sqlx"
)

func RegisterPostgresDB(db *sqlx.DB) {
	pgmodels.RegisterPostgresDB(db)
}
