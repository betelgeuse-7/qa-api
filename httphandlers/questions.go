package httphandlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/betelgeuse-7/qa/storage/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) AskQuestion(c *gin.Context) {
	nqp := &models.NewQuestionPayload{}
	err := c.BindJSON(nqp)
	if err != nil {
		if err.Error() == "EOF" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no json body"})
			return
		}
		c.Status(http.StatusInternalServerError)
		h.logger.Error("*Handler.AskQuestion: bind json: %s\n", err.Error())
		return
	}
	if validationErrs := nqp.Validate(); len(validationErrs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrs})
		return
	}
	userId := c.GetInt64(ContextUserIdKey)
	if userId <= 0 {
		c.Status(http.StatusUnauthorized)
		return
	}
	nqp.UserId = userId
	response, err := h.questionRepo.NewQuestion(nqp)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		h.logger.Error("*Handler.AskQuestion: new question: %s\n", err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "question added", "question": response})
}

func (h *Handler) ViewQuestion(c *gin.Context) {
	questionIdStr := c.Param("id")
	questionId, err := strconv.ParseInt(questionIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid path parameter. need an integer"})
		return
	}
	q, err := h.questionRepo.GetQuestion(questionId)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "no such question"})
			return
		}
		c.Status(http.StatusInternalServerError)
		h.logger.Error("*Handler.ViewQuestion: get question: %s\n", err.Error())
		return
	}
	c.JSON(http.StatusOK, q)
}

// can only update the text or the title
func (h *Handler) UpdateQuestion(c *gin.Context) {
	var payload *models.UpdateQuestionPayload = &models.UpdateQuestionPayload{}
	questionId, err := getQuestionIdParam(c)
	if err != nil {
		return
	}
	if err := checkUserIsTheAuthorOfQuestion(h, c, questionId); err != nil {
		return
	}
	err = c.BindJSON(payload)
	if err != nil {
		if err.Error() == "EOF" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no json body"})
			return
		}
		c.Status(http.StatusInternalServerError)
		h.logger.Error("*Handler.UpdateQuestion: bind json: %s\n", err.Error())
		return
	}
	if validationErrs := payload.Validate(); len(validationErrs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrs})
		return
	}
	res, err := h.questionRepo.UpdateQuestion(questionId, payload)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "no such question"})
			return
		}
		c.Status(http.StatusInternalServerError)
		h.logger.Error("*Handler.UpdateQuestion: UpdateQuestion: %s\n", err.Error())
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (h *Handler) DeleteQuestion(c *gin.Context) {
	questionId, err := getQuestionIdParam(c)
	if err != nil {
		return
	}
	if err := checkUserIsTheAuthorOfQuestion(h, c, questionId); err != nil {
		return
	}
	err = h.questionRepo.DeleteQuestion(questionId)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "no such question"})
			return
		}
		c.Status(http.StatusInternalServerError)
		h.logger.Error("*Handler.UpdateQuestion: UpdateQuestion: %s\n", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted question"})
}

func checkUserIsTheAuthorOfQuestion(h *Handler, c *gin.Context, questionId int64) error {
	userId := c.GetInt64(ContextUserIdKey)
	if userId <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
		return fmt.Errorf("err")
	}
	qs, err := h.questionRepo.GetQuestionStatus(questionId)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "no such question"})
			return err
		}
		c.Status(http.StatusInternalServerError)
		h.logger.Error("*Handler.UpdateQuestion: GetAuthorIdOfQuestion: %s\n", err.Error())
		return err
	}
	if userId != qs.AuthorId {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return fmt.Errorf("err")
	}
	// question's deleted_at column is set, which means, this question is deleted.
	if qs.DeletedAt != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no such question"})
		return fmt.Errorf("err")
	}
	return nil
}

func getQuestionIdParam(c *gin.Context) (int64, error) {
	questionIdStr := c.Param("id")
	questionId, err := strconv.ParseInt(questionIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter. need an integer"})
		return -1, fmt.Errorf("err")
	}
	return questionId, nil
}
