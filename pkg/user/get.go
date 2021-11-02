package user

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hiroyaonoe/le4-db-go/db"
)

func Get(c *gin.Context) {
	db, err := db.NewDB()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	users := []User{}
	err = db.Select(&users, "SELECT * FROM users WHERE user_id = $1", userID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if len(users) == 0 {
		c.String(http.StatusNotFound, fmt.Sprintf("user %d not found", userID))
		return
	}
	user := users[0]

	c.HTML(http.StatusOK, "user.html", gin.H{
		"user": user,
	})
}
