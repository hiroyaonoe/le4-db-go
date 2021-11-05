package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/hiroyaonoe/le4-db-go/db"
	"github.com/hiroyaonoe/le4-db-go/domain"
)

func Authenticate(c *gin.Context) {
	authenticate(c)
	c.Next()
}

func AuthenticateWithRedirect(c *gin.Context) {
	_, ok := c.Get("UserID")
	if !ok {
		c.Redirect(http.StatusMovedPermanently, "/login")
		c.Abort()
	}
	c.Next()
}

func authenticate(c *gin.Context) error {
	session := sessions.Default(c)
	userID := session.Get("UserID")
	if userID == nil {
		return fmt.Errorf("UserID is not setted")
	}

	db, err := db.NewDB()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	user := domain.User{}
	err = db.Get(&user, "SELECT * FROM users WHERE user_id = $1", userID)
	if err != nil {
		return err
	}

	c.Set("UserID", user.UserID)
	c.Set("UserName", user.Name)
	c.Set("UserRole", user.Role)
	return nil
}
