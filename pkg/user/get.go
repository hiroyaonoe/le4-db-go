package user

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
	"github.com/jmoiron/sqlx"
)

func Get(c *gin.Context) {
	db := db.GetDB()

	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	user := domain.User{}
	err = db.Get(&user, "SELECT * FROM users WHERE user_id = $1", userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.String(http.StatusNotFound, fmt.Sprintf("user %d not found", userID))
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	threads := []domain.Thread{}
	query := "SELECT thread_id, title, created_at, category_id, category_name " +
		"FROM threads_with_user_category " +
		"WHERE user_id = $1"
	err = db.Select(&threads, query, userID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	AddTags := []domain.Tag{}
	query = "SELECT tag_id, name, thread_id FROM tag_with_thread_id WHERE thread_id IN (:thread_id)"
	query, argsT, err := sqlx.Named(query, threads)
	query, argsT, err = sqlx.In(query, argsT)
	query = db.Rebind(query)
	err = db.Select(&AddTags, query, argsT...)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	numComments := []domain.Thread{}
	query = "SELECT thread_id, num_comment " +
		"FROM num_comments " +
		"WHERE thread_id IN (:thread_id)"
	query, argsT, err = sqlx.Named(query, threads)
	query, argsT, err = sqlx.In(query, argsT)
	query = db.Rebind(query)
	err = db.Select(&numComments, query, argsT...)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	indexThread := map[int]*domain.Thread{}
	for i := 0; i < len(threads); i++ {
		indexThread[threads[i].ThreadID] = &threads[i]
	}
	for _, v := range AddTags {
		t := indexThread[v.ThreadID]
		t.Tags = append(t.Tags, v)
	}
	for _, v := range numComments {
		t := indexThread[v.ThreadID]
		t.NumComment = v.NumComment
	}

	comments := []domain.Comment{}
	query = "SELECT content, comment_id, thread_id, thread_title, created_at " +
		"FROM comments_with_user_thread " +
		"WHERE user_id = $1"
	err = db.Select(&comments, query, userID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	loginUserID := auth.GetUserID(c)
	loginUserName := auth.GetUserName(c)
	loginUserRole := auth.GetUserRole(c)
	c.HTML(http.StatusOK, "user.html", gin.H{
		"user":            user,
		"threads":         threads,
		"comments":        comments,
		"login_user_id":   loginUserID,
		"login_user_name": loginUserName,
		"login_user_role": loginUserRole,
		"is_user_page":    true,
	})
	return
}
