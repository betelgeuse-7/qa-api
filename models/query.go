package models

import "github.com/doug-martin/goqu"

/*
type Table uint

const (
	USERS Table = iota
	QUESTIONS
	ANSWERS
)

func (t Table) String() string {
	switch t {
	case USERS:
		return "users"
	case QUESTIONS:
		return "questions"
	default:
		return "answers"
	}
}
*/

func (u *User) GetById(id uint) /*(*User, error)*/ string {
	query, _, err := goqu.From("users").Prepared(true).Where(goqu.Ex{
		"user_id": id,
	}).ToSql()

	if err != nil {
		return "panic(err)"
	}

	return query
}
