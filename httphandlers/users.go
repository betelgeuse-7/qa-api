package httphandlers

import (
	"net/http"

	"github.com/betelgeuse-7/qa/service/logger"
	"github.com/gin-gonic/gin"
)

func (h *Handler) NewUser(c *gin.Context) {
	_, err := c.Writer.Write([]byte("NEW USER"))
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		logger.Error(err.Error())
	}
	//h.userRepo.Register()
}
