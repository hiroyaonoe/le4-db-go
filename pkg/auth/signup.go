package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Signup(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", gin.H{})
}
