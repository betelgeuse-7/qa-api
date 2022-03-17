package models

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func RegisterPostgresDB(db *sqlx.DB) {
	DB = db
}

func (u *User) Insert(hashedPwd string) (int64, error) {
	q, args, err := sq.Insert("users").Columns("username", "email", "password", "handle").
		Values(u.Username, u.Email, hashedPwd, u.Handle).ToSql()

	if err != nil {
		return 0, err_QUERY_BUILDING_FAIL()
	}

	res, err := DB.Exec(q, args...)
	if err != nil {
		return 0, err_DB_EXEC_FAIL(err)
	}

	lastInsertedId, err := res.LastInsertId()
	if err != nil {
		return -1, nil
	}

	return lastInsertedId, nil
}
