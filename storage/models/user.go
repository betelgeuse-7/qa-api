package models

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/betelgeuse-7/qa/service/hashpwd"
	"github.com/betelgeuse-7/qa/service/sqlbuild"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Register(*UserRegisterPayload) error
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

func (u *UserRegisterPayload) Validate() []string {
	errs := []string{}
	if len(u.Username) == 0 {
		errs = append(errs, "missing username")
	}
	if len(u.Email) == 0 {
		errs = append(errs, "missing email")
	}
	if len(u.Password) < 6 {
		errs = append(errs, "password not long enough")
	}
	if len(u.Handle) == 0 {
		errs = append(errs, "missing handle")
	}
	if u.Handle[0] == '@' {
		errs = append(errs, "handle cannot start with '@'")
	}
	return errs
}

type UserLoginPayload struct {
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}

type UserProfileResponse struct {
	Username      string            `db:"username" json:"username"`
	Email         string            `db:"email" json:"email"`
	Handle        string            `db:"handle" json:"handle"`
	LastOnline    *time.Time        `db:"last_online" json:"last_online"`
	CreatedAt     *time.Time        `db:"created_at" json:"created_at"`
	LastQuestions []*SingleQuestion `json:"last_questions"`
	LastAnswers   []*SingleAnswer   `json:"last_answers"`
}

func (u *UserRepo) Register(payload *UserRegisterPayload) error {
	hasher := hashpwd.New(payload.Password)
	hasher.HashPwd()
	if err := hasher.Error(); err != nil {
		return err
	}
	payload.Password = hasher.Hashed()
	q, args, _ := u.sqlbuilder.Insert("users").
		Columns("username", "email", "handle", "password").
		Values(payload.Username, payload.Email, payload.Handle, payload.Password).
		ToSql()
	_, err := u.db.Exec(q, args...)
	if err != nil {
		return err
	}
	return nil
}
