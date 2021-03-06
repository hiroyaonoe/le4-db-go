package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hiroyaonoe/le4-db-go/db"
	"github.com/hiroyaonoe/le4-db-go/domain"
	"github.com/hiroyaonoe/le4-db-go/lib/session"
)

func Create(c *gin.Context) {
	db := db.GetDB()

	user := domain.User{}
	user.Name = c.PostForm("user_name")
	if user.Name == "" {
		c.String(http.StatusBadRequest, "user's name cannot be null")
		return
	}
	var err error
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
	err = db.Get(&user.UserID, "SELECT user_id FROM users WHERE name = $1", user.Name)

	session.SetUserID(c, user.UserID)

	id := strconv.Itoa(user.UserID)
	c.Redirect(http.StatusSeeOther, "/user/"+id)
	return
}
