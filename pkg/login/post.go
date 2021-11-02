package login

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hiroyaonoe/le4-db-go/db"
	"github.com/hiroyaonoe/le4-db-go/pkg/session"
	"github.com/hiroyaonoe/le4-db-go/pkg/user"
)

func Post(c *gin.Context) {
	db, err := db.NewDB()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	userName := c.PostForm("user_name")
	password := c.PostForm("password")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	users := []user.User{}
	err = db.Select(&users, "SELECT * FROM users WHERE name = $1", userName)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if len(users) == 0 {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"message": "ユーザー名かパスワードが間違っています",
		})
		return
	}
	u := users[0]
	ok := u.Password.Authenticate(password)
	if !ok {
		c.String(http.StatusUnauthorized, "This user is unauthorized")
		return
	}

	session.SetSession(c, u.UserID)

	id := strconv.Itoa(u.UserID)
	c.Redirect(http.StatusMovedPermanently, "/user/"+id)
	return
}
