package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/hiroyaonoe/le4-db-go/db"
	"github.com/hiroyaonoe/le4-db-go/pkg/user"
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
	users := []user.User{}
	err = db.Select(&users, "SELECT * FROM users WHERE user_id = $1", userID)
	if err != nil {
		return err
	}
	if len(users) == 0 {
		return fmt.Errorf("User not found")
	}
	u := users[0]

	c.Set("UserID", u.UserID)
	c.Set("UserName", u.Name)
	c.Set("UserRole", u.Role)
	return nil
}
