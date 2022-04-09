package pgmodels

import "github.com/jmoiron/sqlx"

var DB *sqlx.DB

func RegisterPostgresDB(db *sqlx.DB) {
	DB = db
}
