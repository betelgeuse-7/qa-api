package httphandlers

import (
	"net/http"

	"github.com/betelgeuse-7/qa/service/logger"
	"github.com/betelgeuse-7/qa/storage/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) NewUser(c *gin.Context) {
	urp := &models.UserRegisterPayload{}
	err := c.BindJSON(urp)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		logger.Error("NewUser: %s\n", err.Error())
		return
	}
	errs := urp.Validate()
	if len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}
	if err := h.userRepo.Register(urp); err != nil {
		c.Status(http.StatusInternalServerError)
		logger.Error("NewUser: %s\n", err.Error())
		return
	}
	// TODO also return an access, and a refresh token
	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}
