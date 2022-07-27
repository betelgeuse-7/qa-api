package httphandlers

import (
	"net/http"

	"database/sql"

	"github.com/betelgeuse-7/qa/service/hashpwd"
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
		h.logger.Error("*Handler.NewUser: %s\n", err.Error())
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
	cookieHttpOnly := true
	c.SetCookie(h.atCookieName, at, int(jwtauth.AT_EXPIRY.Seconds()), "/", h.domain, h.useHTTPS, cookieHttpOnly)
	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

// set access-token cookie after a successfull log in
func (h *Handler) Login(c *gin.Context) {
	ulp := &models.UserLoginPayload{}
	if err := c.BindJSON(ulp); err != nil {
		c.Status(http.StatusInternalServerError)
		h.logger.Error("*Handler.Login: %s\n", err.Error())
		return
	}
	validationErrs := ulp.Validate()
	if len(validationErrs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrs})
		return
	}
	ulr, err := h.userRepo.GetUserLoginResults(ulp.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "no such user"})
			return
		}
		c.Status(http.StatusInternalServerError)
		h.logger.Error("*Handler.Login: %s\n", err.Error())
		return
	}
	if err := hashpwd.CompareHashAndPwd(ulr.Pwd, ulp.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong password"})
		return
	}
	t, err := h.jwtRepo.NewToken(ulr.UserId, jwtauth.NewAccessToken)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		h.logger.Error("*Handler.Error: new token: %s\n", err.Error())
		return
	}
	cookieMaxAge := int(jwtauth.AT_EXPIRY.Seconds())
	cookiePath := "/"
	cookieHttpOnly := true
	c.SetCookie(h.atCookieName, t, cookieMaxAge, cookiePath, h.domain, h.useHTTPS, cookieHttpOnly)
	c.Set(ContextUserIdKey, ulr.UserId)
	c.JSON(http.StatusOK, gin.H{"message": "login successful (no redirect)"})
}
