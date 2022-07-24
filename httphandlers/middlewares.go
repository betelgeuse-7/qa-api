package httphandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const ContextUserIdKey = "user"

func (h *Handler) AuthTokenMiddleware(c *gin.Context) {
	at, err := c.Cookie(h.atCookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing access token cookie"})
			return
		}
		errStr := err.Error()
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errStr})
		return
	}
	atTok, atClaims, err := h.jwtRepo.ParseToken(at)
	if err != nil {
		errStr := err.Error()
		h.logger.Info(errStr)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid access token"})
		return
	}
	atTokValid := atTok.Valid
	atClaimsUserId := atClaims.UserId
	if !(atTokValid) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid access token"})
		return
	}
	c.Set(ContextUserIdKey, atClaimsUserId)
	c.Next()
}
