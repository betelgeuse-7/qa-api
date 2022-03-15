package models

import "time"

type Model interface {
	GetById(id uint) (Model, error)
}

type BaseModel struct {
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	DeletedAt time.Time `json:"deleted_at" db:"deleted_at"`
}

type User struct {
	BaseModel

	UserId     uint      `json:"user_id" db:"user_id"`
	Username   string    `json:"username" db:"username"`
	Email      string    `json:"email" db:"email"`
	Handle     string    `json:"handle" db:"handle"`
	LastOnline time.Time `json:"last_online" db:"last_online"`
}

type Question struct {
	BaseModel

	QuestionId uint     `json:"question_id" db:"question_id"`
	Title      string   `json:"title" db:"title"`
	Text       string   `json:"text" db:"text"`
	Tags       []string `json:"tags"`
	QuestionBy User     `json:"question_by"`
	Upvotes    uint     `json:"upvotes"`
	Downvotes  uint     `json:"downvotes"`
}

type Answer struct {
	BaseModel

	AnswerId  uint   `json:"answer_id" db:"answer_id"`
	AnswerBy  User   `json:"answer_by"`
	Text      string `json:"text" db:"text"`
	Upvotes   uint   `json:"upvotes"`
	Downvotes uint   `json:"downvotes"`
}

type Comment struct {
	BaseModel

	CommentId uint   `json:"comment_id" db:"comment_id"`
	Text      string `json:"text" db:"text"`
	CommentBy User   `json:"comment_by"`
}
