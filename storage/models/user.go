package models

import (
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/betelgeuse-7/qa/service/hashpwd"
	"github.com/betelgeuse-7/qa/service/sqlbuild"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Register(*UserRegisterPayload) (int64, error)
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
	if len(u.Password) == 0 {
		errs = append(errs, "missing password")
	} else if len(u.Password) < 6 {
		errs = append(errs, "password not long enough")
	}
	if len(u.Handle) == 0 {
		errs = append(errs, "missing handle")
	}
	if ok := strings.HasPrefix(u.Handle, "@"); ok {
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

func (u *UserRepo) Register(payload *UserRegisterPayload) (int64, error) {
	hasher := hashpwd.New(payload.Password)
	hasher.HashPwd()
	if err := hasher.Error(); err != nil {
		return -1, err
	}
	payload.Password = hasher.Hashed()
	q, args, _ := u.sqlbuilder.Insert("users").
		Columns("username", "email", "handle", "password").
		Values(payload.Username, payload.Email, payload.Handle, payload.Password).
		Suffix("RETURNING user_id").
		ToSql()
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
