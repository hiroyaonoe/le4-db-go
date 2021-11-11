package login

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hiroyaonoe/le4-db-go/db"
	"github.com/hiroyaonoe/le4-db-go/domain"
	"github.com/hiroyaonoe/le4-db-go/lib/session"
)

func Post(c *gin.Context) {
	db := db.GetDB()

	userName := c.PostForm("user_name")
	password := c.PostForm("password")

	user := domain.User{}
	err := db.Get(&user, "SELECT * FROM users WHERE name = $1", userName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.HTML(http.StatusOK, "login.html", gin.H{
				"message": "ユーザー名かパスワードが間違っています",
			})
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	ok := user.Password.Authenticate(password)
	if !ok {
		c.String(http.StatusUnauthorized, "This user is unauthorized")
		return
	}

	session.SetUserID(c, user.UserID)

	c.Redirect(http.StatusSeeOther, "/")
	return
}
