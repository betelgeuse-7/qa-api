package pgmodels

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

func (u *User) Insert(hashedPwd string) error {
	q, args, err := sq.Insert("users").Columns("username", "email", "password", "handle").
		Values(u.Username, u.Email, hashedPwd, u.Handle).ToSql()

	if err != nil {
		return err_QUERY_BUILDING_FAIL()
	}

	// turn ? bindvars to driver specific bindvars. $n for postgres
	q = DB.Rebind(q)

	_, err = DB.Exec(q, args...)
	if err, ok := err.(*pq.Error); ok {
		fmt.Println(err.Constraint)
		return err_DB_EXEC_FAIL(err)
	}

	return nil
}

func (u *User) GetUser() error {
	if len(u.Handle) < 1 {
		return errors.New("no user handle. need one")
	}
	/*
		return username, handle
	*/
	q, args, err := sq.Select("username, handle").From("users").Where(sq.Eq{
		"handle": u.Handle, "deleted_at": nil,
	}).ToSql()
	if err != nil {
		return err_QUERY_BUILDING_FAIL()
	}
	q = DB.Rebind(q)
	err = sqlx.Get(DB, u, q, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		return err_DB_EXEC_FAIL(err)
	}
	return nil
}

func (u *User) Delete() error {
	if len(u.Handle) < 1 {
		return errors.New("no user handle. need one")
	}
	q, args, err := sq.Update("users").Set("deleted_at", time.Now()).Where(sq.Eq{
		"handle": u.Handle,
	}).ToSql()
	if err != nil {
		return err_QUERY_BUILDING_FAIL()
	}
	q = DB.Rebind(q)
	_, err = DB.Exec(q, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		return err_DB_EXEC_FAIL(err)
	}
	return nil
}

// populate UserId and Password fields of u
func (u *User) GetPassword() error {
	email := u.Email
	username := u.Username
	which := ""    // which field
	whichVal := "" // that field's value

	if len(email) > 0 {
		which = "email"
		whichVal = email
	}
	if len(username) > 0 {
		which = "username"
		whichVal = username
	}
	q, args, err := sq.Select("user_id, password").From("users").Where(sq.Eq{
		which: whichVal, "deleted_at": nil,
	}).ToSql()
	if err != nil {
		return err_QUERY_BUILDING_FAIL()
	}
	q = DB.Rebind(q)
	err = sqlx.Get(DB, u, q, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		return err_DB_EXEC_FAIL(err)
	}
	return nil
}
