package models

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Register(UserRegisterPayload) error
}

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

type UserRegisterPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Handle   string `json:"handle"`
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

func (u *UserRepo) Register(payload UserRegisterPayload) error {
	return nil
}
