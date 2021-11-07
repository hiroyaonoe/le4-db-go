package session

import (
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func NewStore(key, name string) gin.HandlerFunc {
	store := cookie.NewStore([]byte(key))
	return sessions.Sessions(name, store)
}

func SetUserID(c *gin.Context, userID int) {
	session := sessions.Default(c)
	session.Set("UserID", userID)
	session.Save()
}

func GetUserID(c *gin.Context) (interface{}, error) {
	session := sessions.Default(c)
	userID := session.Get("UserID")
	if userID == nil {
		return nil, fmt.Errorf("UserID is not setted")
	}
	return userID, nil
}

func Clear(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
}
