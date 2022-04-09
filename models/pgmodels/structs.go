package pgmodels

import (
	"strings"
	"time"
)

type ValidationErrors []string

type Model interface {
	Validate() (ValidationErrors, bool)
}

type BaseModel struct {
	CreatedAt *time.Time `json:"created_at,omitempty" db:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type User struct {
	BaseModel

	UserId     uint       `json:"user_id,omitempty" db:"user_id"`
	Username   string     `json:"username,omitempty" db:"username"`
	Email      string     `json:"email,omitempty" db:"email"`
	Password   string     `json:"password,omitempty"`
	Handle     string     `json:"handle,omitempty" db:"handle"`
	LastOnline *time.Time `json:"last_online,omitempty" db:"last_online"`
}

func (u *User) Validate() (ValidationErrors, bool) {
	ve := ValidationErrors{}

	if len(u.Username) < 1 {
		ve = append(ve, "a username must be present")
	}
	if len(u.Email) < 1 {
		ve = append(ve, "an email must be present")
	}
	if !(strings.Contains(u.Email, "@")) {
		ve = append(ve, "invalid email")
	}
	if len(u.Password) < 6 {
		ve = append(ve, "the password must be at least 6 characters long")
	}
	if len(u.Handle) < 1 {
		ve = append(ve, "the handle must be at least 1 character long")
	}

	if len(ve) > 0 {
		return ve, false
	}
	return ve, true
}

type Question struct {
	BaseModel

	QuestionId uint     `json:"question_id,omitempty" db:"question_id"`
	Title      string   `json:"title,omitempty" db:"title"`
	Text       string   `json:"text,omitempty" db:"text"`
	Tags       []string `json:"tags,omitempty"`
	QuestionBy User     `json:"question_by,omitempty"`
	Upvotes    uint     `json:"upvotes,omitempty"`
	Downvotes  uint     `json:"downvotes,omitempty"`
}

type Answer struct {
	BaseModel

	AnswerId  uint   `json:"answer_id,omitempty" db:"answer_id"`
	AnswerBy  User   `json:"answer_by,omitempty"`
	Text      string `json:"text,omitempty" db:"text"`
	Upvotes   uint   `json:"upvotes,omitempty"`
	Downvotes uint   `json:"downvotes,omitempty"`
}

type Comment struct {
	BaseModel

	CommentId uint   `json:"comment_id,omitempty" db:"comment_id"`
	Text      string `json:"text,omitempty" db:"text"`
	CommentBy User   `json:"comment_by,omitempty"`
}
