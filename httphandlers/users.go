package httphandlers

import (
	"net/http"

	"github.com/betelgeuse-7/qa/service/jwtauth"
	"github.com/betelgeuse-7/qa/storage/models"
	"github.com/betelgeuse-7/qa/storage/postgres"
	"github.com/gin-gonic/gin"
	pq "github.com/lib/pq"
)

func (h *Handler) NewUser(c *gin.Context) {
	urp := &models.UserRegisterPayload{}
	err := c.BindJSON(urp)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		h.logger.Error("NewUser: %s\n", err.Error())
		return
	}
	errs := urp.Validate()
	if len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}
	userId, err := h.userRepo.Register(urp)
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok && pqError.Code == postgres.ERROR_UNIQUE_VIOLATION {
			c.String(http.StatusBadRequest, "this user already exists")
			h.logger.Info("NewUser: tried to insert duplicate entry\n")
			return
		}
		c.Status(http.StatusInternalServerError)
		h.logger.Error("NewUser: %s\n", err.Error())
		return
	}
	at, err := h.jwtRepo.NewToken(userId, jwtauth.NewAccessToken)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		h.logger.Error("NewUser: %s\n", err.Error())
		return
	}
	cookieSecure := false
	cookieHttpOnly := true
	c.SetCookie(h.atCookieName, at, int(jwtauth.AT_EXPIRY.Seconds()), "/", h.domain, cookieSecure, cookieHttpOnly)
	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}
