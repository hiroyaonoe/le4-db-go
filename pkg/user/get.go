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
	query := "SELECT thread_id, title, created_at, categories.category_id, categories.name AS category_name " +
		"FROM threads " +
		"NATURAL JOIN post_threads " +
		"NATURAL JOIN link_categories " +
		"JOIN categories ON categories.category_id = link_categories.category_id " +
		"WHERE user_id = $1"
	err = db.Select(&threads, query, userID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	AddTags := []domain.Tag{}
	query = "SELECT tag_id, name, thread_id FROM tags NATURAL JOIN add_tags WHERE thread_id IN (:thread_id)"
	query, argsT, err := sqlx.Named(query, threads)
	query, argsT, err = sqlx.In(query, argsT)
	query = db.Rebind(query)
	err = db.Select(&AddTags, query, argsT...)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	numComments := []domain.Thread{}
	query = "SELECT thread_id, COUNT(comment_id) AS num_comment " +
		"FROM comments " +
		"WHERE thread_id IN (:thread_id) " +
		"GROUP BY thread_id"
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

	threadsC := []domain.Thread{}
	query = "SELECT DISTINCT threads.thread_id, threads.title, post_threads.created_at, users.user_id, users.name AS user_name " +
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

	loginUserID := auth.GetUserIDInt(c)
	loginUserRole := auth.GetUserRole(c)
	c.HTML(http.StatusOK, "user.html", gin.H{
		"user":     user,
		"threads":  threads,
		"threadsC": threadsC,
		"userID":   loginUserID,
		"userRole": loginUserRole,
	})
	return
}
