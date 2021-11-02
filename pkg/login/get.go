package login

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
	return
}
