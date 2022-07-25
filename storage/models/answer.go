package models

import "time"

type BasicAnswerResponse struct {
	AnswerId  int64             `json:"answer_id" db:"answer_id"`
	AnswerBy  BasicUserResponse `json:"answer_author"`
	Text      string            `json:"text" db:"text"`
	CreatedAt *time.Time        `json:"created_at" db:"created_at"`
}
