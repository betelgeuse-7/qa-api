package models

import (
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/betelgeuse-7/qa/service/sqlbuild"
	"github.com/jmoiron/sqlx"
)

type AnswerRepository interface {
	NewAnswer(NewAnswerPayload) (NewAnswerResponse, error)
	UpdateAnswer(UpdateAnswerPayload, int64) (UpdateAnswerResponse, error)
	DeleteAnswer(int64) error
	// answerId, userId -> answer.answer_by == userId, err
	AnswerBelongsToUser(int64, int64) (bool, error)
	GetAnswerStatus(int64) (AnswerStatus, error)
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

type NewAnswerPayload struct {
	Text       string     `json:"text" db:"text"`
	ToQuestion int64      `json:"question_id" db:"to_question"`
	AnswerBy   int64      `json:"answer_by" db:"answer_by"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type NewAnswerResponse struct {
	NewAnswerPayload
	AnswerId  int64      `json:"answer_id" db:"answer_id"`
	CreatedAt *time.Time `json:"answered_at" db:"created_at"`
}

func (a *AnswerRepo) NewAnswer(nap NewAnswerPayload) (NewAnswerResponse, error) {
	nar := NewAnswerResponse{}
	q, args, err := a.sqlbuilder.Insert("answers").Columns("text", "to_question", "answer_by").
		Values(nap.Text, nap.ToQuestion, nap.AnswerBy).Suffix("RETURNING *").ToSql()
	if err != nil {
		return nar, fmt.Errorf("error building NewAnswer query: %w", err)
	}
	row := a.db.QueryRowx(q, args...)
	if err := row.StructScan(&nar); err != nil {
		return nar, err
	}
	return nar, err
}

type UpdateAnswerPayload struct {
	Text string `json:"text"`
}

type UpdateAnswerResponse struct {
	UpdateAnswerPayload
}

func (a *AnswerRepo) UpdateAnswer(uap UpdateAnswerPayload, answerId int64) (UpdateAnswerResponse, error) {
	uar := UpdateAnswerResponse{}
	q, args, err := a.sqlbuilder.Update("answers").Set("text", uap.Text).Where(squirrel.Eq{
		"answer_id": answerId, "deleted_at": nil,
	}).Suffix("RETURNING \"text\"").ToSql()
	if err != nil {
		return uar, fmt.Errorf("error while building query for UpdateAnswer: %w", err)
	}
	row := a.db.QueryRowx(q, args...)
	err = row.StructScan(&uar)
	return uar, err
}

func (a *AnswerRepo) AnswerBelongsToUser(answerId, userId int64) (bool, error) {
	q, args, err := a.sqlbuilder.Select("answer_by").From("answers").Where(squirrel.Eq{
		"answer_id": answerId,
	}).ToSql()
	if err != nil {
		return false, err
	}
	row := a.db.QueryRowx(q, args...)
	var answerBy int64
	err = row.Scan(&answerBy)
	return answerBy == userId, err
}

func (a *AnswerRepo) DeleteAnswer(answerId int64) error {
	q, args, err := a.sqlbuilder.Update("answers").Set("deleted_at", time.Now()).Where(squirrel.Eq{
		"deleted_at": nil,
		"answer_id":  answerId,
	}).ToSql()
	if err != nil {
		return fmt.Errorf("error while building query for DeleteAnswer: %w", err)
	}
	_, err = a.db.Exec(q, args...)
	return err
}

type AnswerStatus struct {
	UserId    int64
	DeletedAt *time.Time
}

func (a *AnswerRepo) GetAnswerStatus(answerId int64) (AnswerStatus, error) {
	as := AnswerStatus{}
	q, args, err := a.sqlbuilder.Select("answer_by", "deleted_at").From("answers").
		Where(squirrel.Eq{"answer_id": answerId}).ToSql()
	if err != nil {
		return as, fmt.Errorf("error while building query for GetAnswerStatus: %w", err)
	}
	row := a.db.QueryRowx(q, args...)
	err = row.Scan(&as.UserId, &as.DeletedAt)
	return as, err
}
