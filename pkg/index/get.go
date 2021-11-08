package index

import (
	"net/http"

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

	threads := []domain.Thread{}
	err = db.Select(&threads, "SELECT thread_id, title, created_at, user_id, user_name, category_id, category_name FROM threads_with_user_category LIMIT 10")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	categories := []domain.Category{}
	err = db.Select(&categories, "SELECT category_id, name FROM categories")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	tags := []domain.Tag{}
	err = db.Select(&tags, "SELECT tag_id, name FROM tags")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	AddTags := []domain.Tag{}
	query := "SELECT tag_id, name, thread_id FROM tag_with_thread_id WHERE thread_id IN (:thread_id)"
	query, args, err := sqlx.Named(query, threads)
	query, args, err = sqlx.In(query, args)
	query = db.Rebind(query)
	err = db.Select(&AddTags, query, args...)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	numComments := []domain.Thread{}
	query = "SELECT thread_id, num_comment " +
		"FROM num_comments " +
		"WHERE thread_id IN (:thread_id)"
	query, args, err = sqlx.Named(query, threads)
	query, args, err = sqlx.In(query, args)
	query = db.Rebind(query)
	err = db.Select(&numComments, query, args...)
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

	loginUserID := auth.GetUserID(c)
	loginUserName := auth.GetUserName(c)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"threads":         threads,
		"login_user_id":   loginUserID,
		"login_user_name": loginUserName,
		"categories":      categories,
		"tags":            tags,
	})
	return
}
