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

// TODO this requires authorization
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

type updatableUserField uint

const (
	_USERNAME updatableUserField = iota
	_EMAIL
	_PASSWORD
	_HANDLE
)

func (u *User) update(field updatableUserField, newValue string) error {
	updateBuilder := sq.Update("users")

	switch field {
	case _USERNAME:
		q, args, err := updateBuilder.Set("username", newValue).Where(sq.Eq{
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

	panic("not implemented")
}

func (u *User) UpdateUsername(newUsername string) error {
	return u.update(_USERNAME, newUsername)
}

func (u *User) UpdateEmail(newEmail string) error {
	return u.update(_EMAIL, newEmail)
}

// don't forget @
func (u *User) UpdateHandle(newHandle string) error {
	return u.update(_HANDLE, newHandle)
}

// this is a bit more complex
func (u *User) UpdatePassword(newPassword string) error {
	return u.update(_PASSWORD, newPassword)
}
