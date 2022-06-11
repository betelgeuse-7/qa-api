package httphandlers

import (
	"fmt"
	"net/http"

	"github.com/betelgeuse-7/qa/storage/inmem"
	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) {
	var u *inmem.User = &inmem.User{}
	if err := c.ShouldBindUri(u); err != nil {
		c.JSON(400, gin.H{"err": err})
		return
	}
	u = inmem.GetUser(u.Username)
	if u == nil {
		c.JSON(404, gin.H{"err": "user not found"})
	}
	c.JSON(200, u)
}

func NewUser(c *gin.Context) {
	var u *inmem.User = &inmem.User{}
	if err := c.BindJSON(u); err != nil {
		c.JSON(400, gin.H{"err": err})
		return
	}
	// validate
	// ...
	if err := u.Save(); err != nil {
		c.JSON(400, err)
		return
	}
	fmt.Println(u)
	c.JSON(http.StatusCreated, gin.H{"msg": "user successfully registered"})
}
