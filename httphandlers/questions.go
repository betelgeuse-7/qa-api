package httphandlers

import (
	"net/http"

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
	_, _ = h.questionRepo.GetQuestion(1)
	c.String(200, "he")
}

func (h *Handler) UpdateQuestion(c *gin.Context) {

}

func (h *Handler) DeleteQuestion(c *gin.Context) {

}
