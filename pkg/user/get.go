package user

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hiroyaonoe/le4-db-go/db"
	"github.com/hiroyaonoe/le4-db-go/domain"
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

	users := []domain.User{}
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

	threads := []domain.Thread{}
	err = db.Select(&threads, "SELECT thread_id, title, created_at FROM post_threads NATURAL JOIN threads WHERE user_id = $1", userID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	threadsC := []domain.Thread{}
	query := "SELECT DISTINCT threads.thread_id, threads.title, post_threads.created_at, users.user_id, users.name AS user_name " +
		"FROM post_comments " +
		"JOIN comments ON post_comments.thread_id = comments.thread_id AND post_comments.comment_id = comments.comment_id " +
		"JOIN threads ON comments.thread_id = threads.thread_id " +
		"JOIN post_threads ON threads.thread_id = post_threads.thread_id " +
		"JOIN users ON post_threads.user_id = users.user_id " +
		"WHERE post_comments.user_id = $1"
	err = db.Select(&threadsC, query, userID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.HTML(http.StatusOK, "user.html", gin.H{
		"user":     user,
		"threads":  threads,
		"threadsC": threadsC,
	})
	return
}
