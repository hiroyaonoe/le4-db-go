package thread

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hiroyaonoe/le4-db-go/db"
	"github.com/hiroyaonoe/le4-db-go/domain"
	"github.com/hiroyaonoe/le4-db-go/pkg/auth"
	"github.com/jmoiron/sqlx"
)

func Create(c *gin.Context) {
	db, err := db.NewDB()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	thread := domain.Thread{}
	thread.Title = c.PostForm("thread_title")
	thread.UserID = auth.GetUserIDInt(c)
	thread.CategoryID, err = strconv.Atoi(c.PostForm("category_id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	thread.CreatedAt = domain.NewDateTime(time.Now())

	tagNames := strings.Fields(c.PostForm("thread_tags"))
	tags := make([]domain.Tag, len(tagNames))
	for i, v := range tagNames {
		tags[i].Name = v
	}

	tx, err := db.Beginx()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	defer tx.Rollback()

	err = tx.Get(&thread.ThreadID, "INSERT INTO threads (title) VALUES ($1) RETURNING thread_id", thread.Title)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	_, err = tx.NamedExec("INSERT INTO post_threads (thread_id, user_id, created_at) VALUES (:thread_id, :user_id, :created_at)", thread)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	_, err = tx.NamedExec("INSERT INTO link_categories (thread_id, category_id) VALUES (:thread_id, :category_id)", thread)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if len(tagNames) > 0 {
		query := "INSERT INTO tags (name) VALUES (:name)" +
			"ON CONFLICT (name) DO NOTHING" // Bulk Upsert
		_, err = tx.NamedExec(query, tags)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		query = "SELECT tag_id, name FROM tags WHERE name IN (?)"
		query, args, err := sqlx.In(query, tagNames)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		query = tx.Rebind(query)
		tags = []domain.Tag{}
		err = tx.Select(&tags, query, args...)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		for i := 0; i < len(tags); i++ {
			tags[i].ThreadID = thread.ThreadID
		}
		_, err = tx.NamedExec("INSERT INTO add_tags (thread_id, tag_id) VALUES (:thread_id, :tag_id)", tags)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}

	tx.Commit()

	id := strconv.Itoa(thread.ThreadID)
	c.Redirect(http.StatusMovedPermanently, "/thread/"+id)
	return
}
