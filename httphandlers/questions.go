package httphandlers

import (
	"database/sql"
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
	questionIdStr := c.Param("id")
	questionId, err := strconv.ParseInt(questionIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter. need an integer"})
		return
	}
	// check to see if the question belongs to the user
	userId := c.GetInt64(ContextUserIdKey)
	if userId <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
		return
	}
	authorId, err := h.questionRepo.GetAuthorIdOfQuestion(questionId)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "no such question"})
			return
		}
		c.Status(http.StatusInternalServerError)
		h.logger.Error("*Handler.UpdateQuestion: GetAuthorIdOfQuestion: %s\n", err.Error())
		return
	}
	if authorId != userId {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
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

}
