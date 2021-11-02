package session

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SetSession(c *gin.Context, userID int) {
	session := sessions.Default(c)
	session.Set("UserID", userID)
	session.Save()
}
