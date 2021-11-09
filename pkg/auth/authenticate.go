package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hiroyaonoe/le4-db-go/db"
	"github.com/hiroyaonoe/le4-db-go/domain"
	"github.com/hiroyaonoe/le4-db-go/lib/session"
)

const (
	userIDKey   = "UserID"
	userNameKey = "UserName"
	userRoleKey = "UserRole"
)

func Authenticate(c *gin.Context) {
	authenticate(c)
	c.Next()
}

func AuthenticateWithRedirect(c *gin.Context) {
	_, ok := c.Get(userIDKey)
	if !ok {
		c.Redirect(http.StatusSeeOther, "/login")
		c.Abort()
	}
	c.Next()
}

func authenticate(c *gin.Context) error {
	userID, err := session.GetUserID(c)
	if err != nil {
		return err
	}

	db := db.GetDB()
	user := domain.User{}
	err = db.Get(&user, "SELECT * FROM users WHERE user_id = $1", userID)
	if err != nil {
		return err
	}

	c.Set(userIDKey, user.UserID)
	c.Set(userNameKey, user.Name)
	c.Set(userRoleKey, user.Role)
	return nil
}

func GetUserID(c *gin.Context) int {
	_, exists := c.Get(userIDKey)
	if exists {
		return c.GetInt(userIDKey)
	}
	return -1 // 未ログインなら-1
}

func GetUserName(c *gin.Context) interface{} {
	userName, _ := c.Get(userNameKey)
	return userName
}

func GetUserRole(c *gin.Context) interface{} {
	userRole, _ := c.Get(userRoleKey)
	return userRole
}
