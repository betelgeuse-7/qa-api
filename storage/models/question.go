package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/betelgeuse-7/qa/service/sqlbuild"
	"github.com/jmoiron/sqlx"
)

type QuestionRepository interface {
	NewQuestion(*NewQuestionPayload) (NewQuestionResponse, error)
	GetQuestion(int64) (ViewQuestionResponse, error)
	UpdateQuestion(int64, *UpdateQuestionPayload) (UpdateQuestionResponse, error)
	DeleteQuestion(int64) error
	GetQuestionStatus(int64) (QuestionStatus, error)
	UpvoteQuestion(int64, int64) error
	DownvoteQuestion(int64, int64) error
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
	BasicUserResponse `json:"author"`
	QuestionId        int64                 `json:"question_id" db:"question_id"`
	Title             string                `json:"title" db:"title"`
	Text              string                `json:"text" db:"text"`
	CreatedAt         *time.Time            `json:"created_at" db:"created_at"`
	UpvoteCount       uint64                `json:"upvotes"`
	DownvoteCount     uint64                `json:"downvotes"`
	Answers           []BasicAnswerResponse `json:"answers"`
	Tags              []string              `json:"tags"`
}

func (qr *QuestionRepo) GetQuestion(questionId int64) (ViewQuestionResponse, error) {
	res := ViewQuestionResponse{}
	qs, err := qr.GetQuestionStatus(questionId)
	if err != nil {
		return res, err
	}
	if qs.DeletedAt != nil {
		return res, sql.ErrNoRows
	}
	tags, err := qr.getTagsForQuestion(questionId)
	if err != nil {
		return res, err
	}
	res.Tags = tags
	answers, err := qr.getAnswersForQuestion(questionId)
	if err != nil {
		return res, err
	}
	res.Answers = answers
	upvotes, downvotes, err := qr.getUpvoteAndDownvotesForQuestion(questionId)
	if err != nil {
		return res, err
	}
	res.UpvoteCount = upvotes
	res.DownvoteCount = downvotes
	q, args, err := qr.sqlbuilder.Select("q.question_id", "q.title", "q.text", "q.created_at",
		"u.username", "u.handle", "u.created_at").
		From("questions q").
		InnerJoin("users u ON u.user_id = q.question_by").
		Where(squirrel.Eq{"q.question_id": questionId}).
		ToSql()
	if err != nil {
		return res, err
	}
	row := qr.db.QueryRowx(q, args...)
	err = row.StructScan(&res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (qr *QuestionRepo) getTagsForQuestion(questionId int64) ([]string, error) {
	res := []string{}
	tagsQuery, tagsQueryArgs, err := qr.sqlbuilder.Select("DISTINCT t.tag").From("tags t").
		InnerJoin("question_tags qt ON qt.question_id = $1", questionId).
		ToSql()
	if err != nil {
		return res, err
	}
	rows, err := qr.db.Queryx(tagsQuery, tagsQueryArgs...)
	if err != nil {
		return res, err
	}
	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		if err != nil {
			return res, err
		}
		res = append(res, tag)
	}
	return res, nil
}

func (qr *QuestionRepo) getAnswersForQuestion(questionId int64) ([]BasicAnswerResponse, error) {
	res := []BasicAnswerResponse{}
	q, args, err := qr.sqlbuilder.Select("a.answer_id", "u.username", "u.handle", "u.created_at",
		"a.text", "a.created_at").
		From("answers a").InnerJoin("users u ON a.answer_by = u.user_id").
		Where(squirrel.Eq{"a.to_question": questionId}).
		ToSql()
	if err != nil {
		return res, err
	}
	rows, err := qr.db.Queryx(q, args...)
	if err != nil {
		return res, err
	}
	for rows.Next() {
		var answer BasicAnswerResponse
		err = rows.StructScan(&answer)
		if err != nil {
			return res, err
		}
		res = append(res, answer)
	}
	return res, nil
}

// return: 	upvotes, downvotes, error
func (qr *QuestionRepo) getUpvoteAndDownvotesForQuestion(questionId int64) (uint64, uint64, error) {
	q, args, err := qr.sqlbuilder.Select("COUNT(qu.question_id)").From("question_upvotes qu").
		Where(squirrel.Eq{"qu.question_id": questionId}).ToSql()
	if err != nil {
		return 0, 0, err
	}
	q2, _, err := qr.sqlbuilder.Select("COUNT(qd.question_id)").From("question_downvotes qd").
		Where(squirrel.Eq{"qd.question_id": questionId}).Limit(1).ToSql()
	if err != nil {
		return 0, 0, err
	}
	// i couldn't figure out how to create nested select statements using squirrel
	// squirrel.Expr() maybe?
	query := "SELECT " + "(" + q + ") AS upvotes, (" + q2 + ") AS downvotes;"
	var up uint64
	var down uint64
	row := qr.db.QueryRowx(query, args...)
	err = row.Scan(&up, &down)
	if err != nil {
		return 0, 0, fmt.Errorf("getUpvoteAndDownvotesForQuestion Scanning err: %s", err.Error())
	}
	return up, down, nil
}

type UpdateQuestionPayload struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

func (u *UpdateQuestionPayload) Validate() (res []string) {
	res = []string{}
	if len(u.Title) == 0 && len(u.Text) == 0 {
		res = append(res, "empty payload")
	}
	return
}

type UpdateQuestionResponse struct {
	QuestionId int64 `json:"question_id" db:"question_id"`
	UpdateQuestionPayload
}

func (qr *QuestionRepo) UpdateQuestion(questionId int64, uqp *UpdateQuestionPayload) (UpdateQuestionResponse, error) {
	res := UpdateQuestionResponse{}
	whichFields := []string{}
	if len(uqp.Text) > 0 {
		whichFields = append(whichFields, "text")
	}
	if len(uqp.Title) > 0 {
		whichFields = append(whichFields, "title")
	}
	updateBuilder := qr.sqlbuilder.Update("questions")
	for _, v := range whichFields {
		switch v {
		case "text":
			updateBuilder = updateBuilder.Set("text", uqp.Text)
		default:
			// title
			updateBuilder = updateBuilder.Set("title", uqp.Title)
		}
	}
	q, args, err := updateBuilder.
		Where(squirrel.Eq{
			"question_id": questionId,
			"deleted_at":  nil,
		}).
		Suffix("RETURNING question_id, title, text").ToSql()
	if err != nil {
		return res, err
	}
	row := qr.db.QueryRowx(q, args...)
	err = row.StructScan(&res)
	return res, err
}

type QuestionStatus struct {
	AuthorId  int64      `db:"question_by"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (qr *QuestionRepo) GetQuestionStatus(questionId int64) (QuestionStatus, error) {
	var qs QuestionStatus
	q, args, err := qr.sqlbuilder.Select("question_by", "deleted_at").From("questions").
		Where(squirrel.Eq{"question_id": questionId}).Limit(1).ToSql()
	if err != nil {
		return qs, err
	}
	row := qr.db.QueryRowx(q, args...)
	err = row.StructScan(&qs)
	return qs, err
}

func (qr *QuestionRepo) DeleteQuestion(questionId int64) error {
	q, args, err := qr.sqlbuilder.Update("questions").
		Set("deleted_at", time.Now()).
		Where(squirrel.Eq{
			"question_id": questionId,
		}).ToSql()
	if err != nil {
		return err
	}
	_, err = qr.db.Exec(q, args...)
	return err
}

const (
	ERROR_UPVOTE_OWN_QUESTION   = "cannot upvote own question"
	ERROR_DOWNVOTE_OWN_QUESTION = "cannot downvote own question"
)

func (qr *QuestionRepo) UpvoteQuestion(questionId, upvoteBy int64) error {
	// can't upvote own question
	qs, err := qr.GetQuestionStatus(questionId)
	if err != nil {
		return err
	}
	if qs.AuthorId == upvoteBy {
		return fmt.Errorf(ERROR_UPVOTE_OWN_QUESTION)
	}
	return voteQuestion(qr, "upvote", questionId, upvoteBy)
}

func (qr *QuestionRepo) DownvoteQuestion(questionId, downvoteBy int64) error {
	qs, err := qr.GetQuestionStatus(questionId)
	if err != nil {
		return err
	}
	if qs.AuthorId == downvoteBy {
		return fmt.Errorf(ERROR_DOWNVOTE_OWN_QUESTION)
	}
	return voteQuestion(qr, "downvote", questionId, downvoteBy)
}
