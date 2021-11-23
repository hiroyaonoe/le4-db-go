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
	db := db.GetDB()

	threadID, err := strconv.Atoi(c.Param("thread_id"))
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	thread := domain.Thread{}
	query := "SELECT thread_id, title, created_at, user_id, user_name, category_id, category_name " +
		"FROM threads_with_user_category " +
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
	query = "SELECT comment_id, thread_id, content, created_at, user_id, user_name " +
		"FROM comments_with_user " +
		"WHERE thread_id = $1 " +
		"ORDER BY created_at ASC"
	err = db.Select(&comments, query, threadID)
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

	loginUserID := auth.GetUserID(c)
	loginUserName := auth.GetUserName(c)
	loginUserRole := auth.GetUserRole(c)

	c.HTML(http.StatusOK, "thread.html", gin.H{
		"thread":          thread,
		"comments":        comments,
		"login_user_id":   loginUserID,
		"login_user_name": loginUserName,
		"login_user_role": loginUserRole,
		"tags":            tags,
	})
	return
}
