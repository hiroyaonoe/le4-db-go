package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hiroyaonoe/le4-db-go/db"
	"github.com/hiroyaonoe/le4-db-go/domain"
	"github.com/hiroyaonoe/le4-db-go/pkg/session"
)

func Create(c *gin.Context) {
	db, err := db.NewDB()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	user := domain.User{}
	user.Name = c.PostForm("user_name")
	user.Password, err = domain.NewPassword(c.PostForm("password"))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	user.Role = domain.MEMBER

	_, err = db.NamedExec("INSERT INTO users (name, password, role) VALUES (:name, :password, :role)", user)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	ids := []int{}
	err = db.Select(&ids, "SELECT user_id FROM users WHERE name = $1", user.Name)
	user.UserID = ids[0]

	session.SetSession(c, user.UserID)

	id := strconv.Itoa(user.UserID)
	c.Redirect(http.StatusMovedPermanently, "/user/"+id)
	return
}
