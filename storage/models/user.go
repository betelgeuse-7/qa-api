package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/betelgeuse-7/okay"
	"github.com/betelgeuse-7/qa/service/hashpwd"
	"github.com/betelgeuse-7/qa/service/sqlbuild"
	"github.com/betelgeuse-7/qa/storage/postgres"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Register(*UserRegisterPayload) (int64, error)
	GetUserLoginResults(string) (UserLoginResults, error)
	DeleteUser(int64) error
	IsUserDeleted(int64) (bool, error)
	GetUserProfile(int64, ServerInfo) (UserProfileResponse, error)
}

type UserRepo struct {
	db         *sqlx.DB
	sqlbuilder squirrel.StatementBuilderType
}

func NewUserRepo(db *sqlx.DB, builder *sqlbuild.Builder) *UserRepo {
	return &UserRepo{db: db, sqlbuilder: builder.B}
}

type UserRegisterPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Handle   string `json:"handle"`
}

func (u *UserRegisterPayload) Okay() (okay.ValidationErrors, error) {
	o := okay.New()
	o.Text(u.Username, "username").Required().IsAlphanumeric()
	o.Text(u.Email, "email").Required().IsEmail()
	o.Text(u.Password, "password").Required().MinLength(6)
	o.Text(u.Handle, "handle").Required().IsAlphanumeric().DoesNotStartWith("@")
	return o.Errors()
}

func (u *UserRegisterPayload) Validate() ([]string, error) {
	return okay.Validate(u)
}

type UserLoginPayload struct {
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}

func (u *UserLoginPayload) Okay() (okay.ValidationErrors, error) {
	o := okay.New()
	o.Text(u.Email, "email").Required().IsEmail()
	o.Text(u.Password, "password").Required().MinLength(6)
	return o.Errors()
}

// the information necessary for the controller, for authorization purposes
type UserLoginResults struct {
	Pwd    string `db:"password"`
	UserId int64  `db:"user_id"`
}

func (u *UserLoginPayload) Validate() ([]string, error) {
	return okay.Validate(u)
}

type BasicUserResponse struct {
	Username string `json:"username" db:"username"`
	Handle   string `json:"handle" db:"handle"`
	// ? why isn't this populated when scanning results from db ? vvv
	CreatedAt *time.Time `json:"registered_at" db:"created_at"`
}

func (u *UserRepo) Register(payload *UserRegisterPayload) (int64, error) {
	hasher := hashpwd.New(payload.Password)
	hasher.HashPwd()
	if err := hasher.Error(); err != nil {
		return -1, err
	}
	payload.Password = hasher.Hashed()
	q, args, err := u.sqlbuilder.Insert("users").
		Columns("username", "email", "handle", "password").
		Values(payload.Username, payload.Email, payload.Handle, payload.Password).
		Suffix("RETURNING user_id").
		ToSql()
	if err != nil {
		return -1, err
	}
	// begin a transaction here, to make sure we don't insert a new row, in case we can't get the
	// last inserted id. getting last inserted id is important, because we need it to build access,
	// and refresh tokens upon registration.
	tx, err := u.db.Begin()
	if err != nil {
		return -1, err
	}
	var userId int64
	row := tx.QueryRow(q, args...)
	err = row.Scan(&userId)
	if err != nil {
		tx.Rollback()
		return -1, err
	}
	tx.Commit()
	return userId, nil
}

// returns bcrypt-hashed password, and an error
func (u *UserRepo) GetUserLoginResults(email string) (UserLoginResults, error) {
	ulr := UserLoginResults{}
	q, args, err := u.sqlbuilder.Select("user_id", "password").From("users").Where(squirrel.Eq{
		"email":      email,
		"deleted_at": nil,
	}).Limit(1).ToSql()
	if err != nil {
		return ulr, err
	}
	if err := u.db.Get(&ulr, q, args...); err != nil {
		return ulr, err
	}
	return ulr, nil
}

func (u *UserRepo) IsUserDeleted(userId int64) (bool, error) {
	var deletedAt *time.Time
	q, args, err := u.sqlbuilder.Select("deleted_at").From("users").Where(squirrel.Eq{
		"user_id": userId,
	}).ToSql()
	if err != nil {
		return true, fmt.Errorf("error construction query")
	}
	row := u.db.QueryRowx(q, args...)
	row.Scan(&deletedAt)
	return deletedAt != nil, nil
}

func (u *UserRepo) DeleteUser(userId int64) error {
	q, args, err := u.sqlbuilder.Update("users").Set("deleted_at", time.Now()).Where(squirrel.Eq{
		"user_id": userId,
	}).ToSql()
	if err != nil {
		return err
	}
	_, err = u.db.Exec(q, args...)
	if err != nil {
		return err
	}
	return nil
}

type UserLastQuestionResponse struct {
	Id        int64      `db:"question_id" json:"-"`
	Title     string     `db:"title" json:"title"`
	Text      string     `db:"text" json:"question_text"`
	CreatedAt *time.Time `db:"created_at" json:"asked_at"`
	Link      string     `json:"question_link"`
}

type UserLastAnswerResponse struct {
	Id        int64      `db:"answer_id" json:"-"`
	Text      string     `db:"text" json:"answer_text"`
	CreatedAt *time.Time `db:"created_at" json:"answered_at"`
	Link      string     `json:"answer_link"`
}

type UserProfileResponse struct {
	Username       string                     `db:"username" json:"username"`
	Handle         string                     `db:"handle" json:"handle"`
	CreatedAt      *time.Time                 `db:"created_at" json:"registered_at"`
	TotalUpvotes   int64                      `json:"total_upvotes"`
	TotalDownvotes int64                      `json:"total_downvotes"`
	LastQuestions  []UserLastQuestionResponse `json:"last_questions"`
	LastAnswers    []UserLastAnswerResponse   `json:"last_answers"`
}

type ServerInfo struct {
	Domain string
	Ssl    bool
}

const POSTGRES_INVALID_ID = "max integer value exceeded"

func (u *UserRepo) GetUserProfile(userId int64, serverInfo ServerInfo) (UserProfileResponse, error) {
	res := UserProfileResponse{}
	if userId > int64(postgres.MAX_INT_VAL) {
		return res, errors.New(POSTGRES_INVALID_ID)
	}
	q, args, err := u.sqlbuilder.Select("username", "handle", "created_at").From("users").
		Where(squirrel.Eq{"deleted_at": nil, "user_id": userId}).ToSql()
	if err != nil {
		return res, err
	}
	row := u.db.QueryRowx(q, args...)
	err = row.StructScan(&res)
	if err != nil {
		return res, err
	}
	limit := uint64(10)
	lastAnswers, err := u.getAnswersForUser(userId, limit, serverInfo)
	if err != nil {
		return res, err
	}
	lastQuestions, err := u.getQuestionsForUser(userId, limit, serverInfo)
	if err != nil {
		return res, err
	}
	res.LastAnswers = lastAnswers
	res.LastQuestions = lastQuestions
	downs, ups, err := u.getTotalUpvoteDownvotes(userId)
	if err != nil {
		return res, err
	}
	res.TotalDownvotes = downs
	res.TotalUpvotes = ups
	return res, nil
}

// upvotes, downvotes, error
func (u *UserRepo) getTotalUpvoteDownvotes(userId int64) (int64, int64, error) {
	q, args, err := u.sqlbuilder.Select("COUNT(qu.question_id)").From("question_upvotes qu").
		Where(squirrel.Eq{
			"qu.upvote_by": userId,
		}).ToSql()
	if err != nil {
		return -1, -1, err
	}
	q2, _, err := u.sqlbuilder.Select("COUNT(qd.question_id)").From("question_downvotes qd").
		Where(squirrel.Eq{
			"qd.downvote_by": userId,
		}).ToSql()
	if err != nil {
		return -1, -1, err
	}
	query := "SELECT " + "(" + q + ") AS upvotes, (" + q2 + ") AS downvotes;"
	var up, down int64
	row := u.db.QueryRowx(query, args...)
	err = row.Scan(&up, &down)
	return up, down, err
}

func (u *UserRepo) getAnswersForUser(userId int64, limit uint64, serverInfo ServerInfo) ([]UserLastAnswerResponse, error) {
	_, ular, err := __getForUser(u, userId, limit, "answers", serverInfo)
	if ular == nil {
		return []UserLastAnswerResponse{}, err
	}
	return *ular, err
}

func (u *UserRepo) getQuestionsForUser(userId int64, limit uint64, serverInfo ServerInfo) ([]UserLastQuestionResponse, error) {
	ulqr, _, err := __getForUser(u, userId, limit, "questions", serverInfo)
	if ulqr == nil {
		return []UserLastQuestionResponse{}, err
	}
	return *ulqr, err
}

func __getForUser(u *UserRepo, userId int64, limit uint64, table string, serverInfo ServerInfo) (*[]UserLastQuestionResponse, *[]UserLastAnswerResponse, error) {
	var query string
	var arguments []interface{}

	switch table {
	case "questions":
		q, args, err := u.sqlbuilder.Select("question_id", "title", "text", "created_at").From("questions").
			Where(squirrel.Eq{"deleted_at": nil, "question_by": userId}).Limit(limit).ToSql()
		if err != nil {
			return nil, nil, err
		}
		query = q
		arguments = args
	case "answers":
		q, args, err := u.sqlbuilder.Select("answer_id", "text", "created_at").From("answers").
			Where(squirrel.Eq{"deleted_at": nil, "answer_by": userId}).Limit(limit).ToSql()
		if err != nil {
			return nil, nil, err
		}
		query = q
		arguments = args
	default:
		return nil, nil, errors.New("__getForUser: unknown table: '" + table + "'")
	}
	var lastAnswers []UserLastAnswerResponse
	var lastQuestions []UserLastQuestionResponse

	rows, err := u.db.Queryx(query, arguments...)
	if err != nil {
		return nil, nil, err
	}
	if table == "questions" {
		var ulqr UserLastQuestionResponse
		for rows.Next() {
			err = rows.StructScan(&ulqr)
			if err != nil {
				return nil, nil, err
			}
			ulqr.Link = generateLink(serverInfo.Domain, "questions", ulqr.Id, serverInfo.Ssl)
			lastQuestions = append(lastQuestions, ulqr)
		}
		err := rows.Err()
		return &lastQuestions, nil, err
	}
	var ular UserLastAnswerResponse
	for rows.Next() {
		err = rows.StructScan(&ular)
		if err != nil {
			return nil, nil, err
		}
		ular.Link = generateLink(serverInfo.Domain, "answers", ular.Id, serverInfo.Ssl)
		lastAnswers = append(lastAnswers, ular)
	}
	return nil, &lastAnswers, rows.Err()
}

func generateLink(domain, resource string, id int64, ssl bool) string {
	scheme := "http"
	if ssl {
		scheme += "s"
	}
	scheme += "://"
	res := ""
	domain = strings.TrimSuffix(domain, "/")
	resource = strings.TrimSuffix(resource, "/")
	res += fmt.Sprintf("%s%s/api/v1/%s/%d/", scheme, domain, resource, id)
	return res
}
