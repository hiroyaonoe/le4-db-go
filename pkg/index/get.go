package index

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hiroyaonoe/le4-db-go/db"
	"github.com/hiroyaonoe/le4-db-go/pkg/category"
	"github.com/hiroyaonoe/le4-db-go/pkg/thread"
)

func Get(c *gin.Context) {
	db, err := db.NewDB()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	threads := []thread.Thread{}
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

	categories := []category.Category{}
	err = db.Select(&categories, "SELECT category_id, name FROM categories")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"threads":     threads,
		"user_id":     userID,
		"user_name":   userName,
		"user_exists": ok,
		"categories":  categories,
	})
	return
}
