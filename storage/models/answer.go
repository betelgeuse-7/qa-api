package models

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/betelgeuse-7/qa/service/sqlbuild"
	"github.com/jmoiron/sqlx"
)

type AnswerRepository interface {
	NewAnswer(int64)
}

type AnswerRepo struct {
	db         *sqlx.DB
	sqlbuilder squirrel.StatementBuilderType
}

func NewAnswerRepo(db *sqlx.DB, sqlbuilder *sqlbuild.Builder) *AnswerRepo {
	return &AnswerRepo{db: db, sqlbuilder: sqlbuilder.B}
}

type BasicAnswerResponse struct {
	AnswerId          int64      `json:"answer_id" db:"answer_id"`
	Text              string     `json:"text" db:"text"`
	CreatedAt         *time.Time `json:"created_at" db:"created_at"`
	BasicUserResponse `json:"answer_author"`
}

func (a *AnswerRepo) NewAnswer(toQuestion int64) {}
