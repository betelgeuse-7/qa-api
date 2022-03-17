package models

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var DB *sqlx.DB

func RegisterPostgresDB(db *sqlx.DB) {
	DB = db
}

func (u *User) Insert(hashedPwd string) error {
	q, args, err := sq.Insert("users").Columns("username", "email", "password", "handle").
		Values(u.Username, u.Email, hashedPwd, u.Handle).ToSql()

	if err != nil {
		return err_QUERY_BUILDING_FAIL()
	}

	// turn ? bindvars to driver specific bindvars
	// $n for postgres (where n is > 0)
	q = DB.Rebind(q)

	_, err = DB.Exec(q, args...)
	if err, ok := err.(*pq.Error); ok {
		fmt.Println(err.Constraint)
		return err_DB_EXEC_FAIL(err)
	}

	return nil
}
