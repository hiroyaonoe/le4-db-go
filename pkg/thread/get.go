package thread

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hiroyaonoe/le4-db-go/db"
	"github.com/hiroyaonoe/le4-db-go/domain"
	"github.com/hiroyaonoe/le4-db-go/pkg/auth"
)

func Get(c *gin.Context) {
	db, err := db.NewDB()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	threadID, err := strconv.Atoi(c.Param("thread_id"))
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	thread := domain.Thread{}
	query := "SELECT thread_id, title, created_at, user_id, users.name AS user_name, categories.category_id, categories.name AS category_name " +
		"FROM threads " +
		"NATURAL JOIN post_threads " +
		"NATURAL JOIN users " +
		"NATURAL JOIN link_categories " +
		"JOIN categories ON categories.category_id = link_categories.category_id " +
		"WHERE thread_id = $1"
	err = db.Get(&thread, query, threadID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.String(http.StatusNotFound, fmt.Sprintf("thread %d not found", threadID))
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	comments := []domain.Comment{}
	err = db.Select(&comments, "SELECT comment_id, thread_id, content, created_at, user_id, users.name AS user_name FROM comments NATURAL JOIN post_comments NATURAL JOIN users WHERE thread_id = $1 ORDER BY created_at ASC", threadID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	tags := []domain.Tag{}
	err = db.Select(&tags, "SELECT tag_id, name FROM tags NATURAL JOIN add_tags WHERE thread_id = $1", threadID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	userID := auth.GetUserIDInt(c)
	userRole := auth.GetUserRole(c)

	c.HTML(http.StatusOK, "thread.html", gin.H{
		"thread":   thread,
		"comments": comments,
		"userID":   userID,
		"userRole": userRole,
		"tags": tags,
	})
	return
}
