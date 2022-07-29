package models

import "fmt"

// type_ is either "downvote", or "upvote"
func voteQuestion(qr *QuestionRepo, type_ string, questionId, voteBy int64) error {
	table := ""
	columns := []string{}
	switch type_ {
	case "downvote":
		table = "question_downvotes"
		columns = append(columns, "question_id", "downvote_by")
	case "upvote":
		table = "question_upvotes"
		columns = append(columns, "question_id", "upvote_by")
	default:
		return fmt.Errorf("models.voteQuestion: invalid vote type '%s'", type_)
	}
	q, args, err := qr.sqlbuilder.Insert(table).Columns(columns...).Values(questionId, voteBy).ToSql()
	if err != nil {
		return err
	}
	_, err = qr.db.Exec(q, args...)
	return err
}
