package user

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hiroyaonoe/le4-db-go/db"
	"github.com/hiroyaonoe/le4-db-go/pkg/thread"
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

	threads := []thread.Thread{}
	err = db.Select(&threads, "SELECT thread_id, title, created_at FROM post_threads NATURAL JOIN threads WHERE user_id = $1", userID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.HTML(http.StatusOK, "user.html", gin.H{
		"user": user,
		"threads": threads,
	})
}
