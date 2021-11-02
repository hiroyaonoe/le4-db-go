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
	}

	threads := []thread.Thread{}
	err = db.Select(&threads, "SELECT thread_id, title, created_at FROM threads NATURAL JOIN post_threads")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	userID, ok := c.Get("UserID")
	userName, _ := c.Get("UserName")

	categories := []category.Category{}
	err = db.Select(&categories, "SELECT category_id, name FROM categories")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"threads":     threads,
		"user_id":     userID,
		"user_name":   userName,
		"user_exists": ok,
		"categories": categories,
	})
}
