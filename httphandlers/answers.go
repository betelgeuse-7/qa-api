package httphandlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/betelgeuse-7/qa/storage/models"
	"github.com/betelgeuse-7/qa/storage/postgres"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func (h *Handler) NewAnswer(c *gin.Context) {
	questionId, err := getInt64IdParam(c)
	if err != nil {
		return
	}
	newAnswerPayload := models.NewAnswerPayload{}
	newAnswerPayload.ToQuestion = questionId
	answerBy := c.GetInt64(ContextUserIdKey)
	newAnswerPayload.AnswerBy = answerBy
	if err = c.BindJSON(&newAnswerPayload); err != nil {
		c.Status(http.StatusInternalServerError)
		h.logger.Error("*Handler.NewAnswer: bind json: %s\n", err.Error())
		return
	}
	if len(newAnswerPayload.Text) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing text"})
		return
	}
	nar, err := h.answerRepo.NewAnswer(newAnswerPayload)
	if err != nil {
		// no question with provided question id
		if err := err.(*pq.Error); err.Code == postgres.ERROR_FOREIGN_KEY_VIOLATION {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no such question"})
			return
		}
		c.Status(http.StatusInternalServerError)
		h.logger.Error("*Handler.NewAnswer: *Handler.answerRepo.NewAnswer: %s\n", err.Error())
		return
	}
	msg := gin.H{"message": "answered successfully", "record": gin.H{
		"text": nar.Text, "answered_at": nar.CreatedAt,
	}}
	c.JSON(http.StatusCreated, msg)
}

func (h *Handler) DeleteAnswer(c *gin.Context) {
	answerId, err := getInt64IdParam(c)
	if err != nil {
		return
	}
	userId := c.GetInt64(ContextUserIdKey)
	as, err := h.answerRepo.GetAnswerStatus(answerId)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no such answer"})
			return
		}
		c.Status(http.StatusInternalServerError)
		h.logger.Error("*Handler.DeleteAnswer: GetAnswerStatus: %s\n", err.Error())
		return
	}
	if as.UserId != userId {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
		return
	}
	if as.DeletedAt != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no such answer"})
		return
	}
	err = h.answerRepo.DeleteAnswer(answerId)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no such answer"})
			return
		}
		c.Status(http.StatusInternalServerError)
		h.logger.Error("*Handler.DeleteAnswer: delete answer, err: %s\n", err.Error())
		return
	}
	msg := gin.H{"message": fmt.Sprintf("deleted answer with id '%d'", answerId)}
	c.JSON(http.StatusOK, msg)
}

func (h *Handler) UpdateAnswer(c *gin.Context) {
	answerId, err := getInt64IdParam(c)
	if err != nil {
		return
	}
	userId := c.GetInt64(ContextUserIdKey)
	ok, err := h.answerRepo.AnswerBelongsToUser(answerId, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no such answer"})
			return
		}
		c.Status(http.StatusInternalServerError)
		h.logger.Error("*Handler.UpdateAnswer: answer belongs to user, err: %s\n", err.Error())
		return
	}
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
		return
	}
	var uap models.UpdateAnswerPayload
	if err := c.BindJSON(&uap); err != nil {
		c.Status(http.StatusInternalServerError)
		h.logger.Error("*Handler.UpdateAnswer: %s\n", err.Error())
		return
	}
	if len(uap.Text) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing text"})
		return
	}
	uar, err := h.answerRepo.UpdateAnswer(uap, answerId)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no such answer"})
			return
		}
		c.Status(http.StatusInternalServerError)
		h.logger.Error("*Handler.UpdateAnswer: update answer, err: %s\n", err.Error())
		return
	}
	msg := gin.H{"message": "updated answer", "record": gin.H{"text": uar.Text}}
	c.JSON(http.StatusCreated, msg)
}
