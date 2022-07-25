package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/betelgeuse-7/qa/service/sqlbuild"
	"github.com/jmoiron/sqlx"
)

type QuestionRepository interface {
	NewQuestion(*NewQuestionPayload) (NewQuestionResponse, error)
	GetQuestion(questionId int64) (ViewQuestionResponse, error)
	//DeleteQuestion()
	//UpdateQuestion()
}

type QuestionRepo struct {
	db         *sqlx.DB
	sqlbuilder squirrel.StatementBuilderType
}

func NewQuestionRepo(db *sqlx.DB, sqlbuilder *sqlbuild.Builder) *QuestionRepo {
	return &QuestionRepo{
		db:         db,
		sqlbuilder: sqlbuilder.B,
	}
}

type NewQuestionPayload struct {
	UserId int64  // set this from request context's user id
	Title  string `json:"title"`
	Text   string `json:"text"`
}

type NewQuestionResponse struct {
	QuestionId int64      `db:"question_id" json:"question_id"`
	Title      string     `db:"title" json:"title"`
	Text       string     `db:"text" json:"text"`
	CreatedAt  *time.Time `db:"created_at" json:"created_at"`
}

func (nqp *NewQuestionPayload) Validate() []string {
	errs := []string{}
	if len(nqp.Title) == 0 {
		errs = append(errs, "missing title")
	}
	if len(nqp.Text) == 0 {
		errs = append(errs, "missing text")
	}
	return errs
}

func (qr *QuestionRepo) NewQuestion(payload *NewQuestionPayload) (NewQuestionResponse, error) {
	res := NewQuestionResponse{}
	title, text, questionBy := payload.Title, payload.Text, payload.UserId
	q, args, err := qr.sqlbuilder.Insert("questions").
		Columns("title", "text", "question_by").
		Values(title, text, questionBy).
		Suffix("RETURNING question_id, title, text, created_at").
		ToSql()
	if err != nil {
		return res, err
	}
	tx, err := qr.db.Beginx()
	if err != nil {
		return res, errors.New("could not begin a new transaction")
	}
	row := tx.QueryRowx(q, args...)
	err = row.StructScan(&res)
	if err != nil {
		tx.Rollback()
		return res, err
	}
	tx.Commit()
	return res, nil
}

type ViewQuestionResponse struct {
	QuestionId    int64                 `json:"question_id" db:"question_id"`
	Title         string                `json:"title" db:"title"`
	Text          string                `json:"text" db:"text"`
	CreatedAt     *time.Time            `json:"created_at" db:"created_at"`
	Author        BasicUserResponse     `json:"author"`
	Tags          []string              `json:"tags"`
	UpvoteCount   uint64                `json:"upvotes"`
	DownvoteCount uint64                `json:"downvotes"`
	Answers       []BasicAnswerResponse `json:"answers"`
}

func (qr *QuestionRepo) GetQuestion(questionId int64) (ViewQuestionResponse, error) {
	res := ViewQuestionResponse{}
	// ! HELP
	// i am bad at sql
	q, args, err := qr.sqlbuilder.Select("q.question_id", "q.title", "q.text", "q.created_at",
		"u.username", "u.handle", "u.created_at",
		"t.tag").From("questions q").InnerJoin("users u ON u.user_id = q.question_by").
		InnerJoin("question_tags qt ON qt.question_id = q.question_id").
		InnerJoin("tags t ON t.tag_id = qt.tag_id").
		ToSql()
	if err != nil {
		return res, err
	}
	fmt.Println("query: ", q)
	fmt.Println("args: ", args)

	return res, nil
}
