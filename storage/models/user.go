package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/betelgeuse-7/qa/service/hashpwd"
	"github.com/betelgeuse-7/qa/service/sqlbuild"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Register(*UserRegisterPayload) (int64, error)
	GetUserLoginResults(string) (UserLoginResults, error)
	DeleteUser(int64) error
	IsUserDeleted(int64) (bool, error)
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

// the information necessary for the controller, for authorization purposes
type UserLoginResults struct {
	Pwd    string `db:"password"`
	UserId int64  `db:"user_id"`
}

func (u *UserLoginPayload) Validate() []string {
	errs := []string{}
	if len(u.Email) == 0 {
		errs = append(errs, "missing email")
	}
	if len(u.Email) > 0 {
		if !(strings.Contains(u.Email, "@")) {
			errs = append(errs, "invalid email")
		}
	}
	if len(u.Password) == 0 {
		errs = append(errs, "missing password")
	} else if len(u.Password) < 6 {
		errs = append(errs, "password not long enough")
	}
	return errs
}

type UserProfileResponse struct {
	Username   string     `db:"username" json:"username"`
	Email      string     `db:"email" json:"email"`
	Handle     string     `db:"handle" json:"handle"`
	LastOnline *time.Time `db:"last_online" json:"last_online"`
	CreatedAt  *time.Time `db:"created_at" json:"created_at"`
	//LastQuestions []*SingleQuestion `json:"last_questions"`
	LastAnswers []*BasicAnswerResponse `json:"last_answers"`
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
