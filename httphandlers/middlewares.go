package httphandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AuthTokenMiddleware(c *gin.Context) {
	at, err := c.Cookie(h.atCookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing access token cookie"})
			h.logger.Info("Missing %s cookie header\n", h.atCookieName)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		h.logger.Error("*Handler.TokenAuthMiddleware: err: %s\n", err.Error())
		return
	}
	rt, err := c.Cookie(h.rtCookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing refresh token cookie"})
			h.logger.Info("Missing %s cookie header\n", h.rtCookieName)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		h.logger.Error("*Handler.TokenAuthMiddleware: err: %s\n", err.Error())
		return
	}
	atTok, atClaims, err := h.jwtRepo.ParseToken(at)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		h.logger.Error("*Handler.TokenAuthMiddleware: err: %s\n", err.Error())
		return
	}
	rtTok, rtClaims, err := h.jwtRepo.ParseToken(rt)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		h.logger.Error("*Handler.TokenAuthMiddleware: err: %s\n", err.Error())
		return
	}
	atTokValid := atTok.Valid
	rtTokValid := rtTok.Valid
	atClaimsUserId := atClaims.UserId
	rtClaimsUserId := rtClaims.UserId

	if !(atTokValid) && !(rtTokValid) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "both the access, and the refresh tokens are invalid"})
		return
	}
	if !(atTokValid) {
		if rtTokValid {
			// TODO
		}
	}

	// if access token is not valid, abort

	// if refresh token is not valid, abort

	// if both of them are valid, c.Next()

	// set user_id in request context to token's user id

	c.Next()
}
