package index

import (
	"net/http"

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

	threads := []domain.Thread{}
	query := "SELECT thread_id, title, users.user_id, users.name AS user_name, created_at, categories.category_id, categories.name AS category_name " +
		"FROM threads " +
		"NATURAL JOIN post_threads " +
		"NATURAL JOIN users " +
		"NATURAL JOIN link_categories " +
		"JOIN categories ON categories.category_id = link_categories.category_id"
	err = db.Select(&threads, query)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	userID, ok := c.Get("UserID")
	userName, _ := c.Get("UserName")

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
	err = db.Select(&AddTags, "SELECT tag_id, name, thread_id FROM tags NATURAL JOIN add_tags")
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

	c.HTML(http.StatusOK, "index.html", gin.H{
		"threads":     threads,
		"user_id":     userID,
		"user_name":   userName,
		"user_exists": ok,
		"categories":  categories,
		"tags":        tags,
	})
	return
}
