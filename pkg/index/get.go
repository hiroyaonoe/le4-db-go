package index

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hiroyaonoe/le4-db-go/db"
	"github.com/hiroyaonoe/le4-db-go/pkg/thread"
)

func Get(c *gin.Context) {
	db, err := db.NewDB()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	threads := []thread.Thread{}
	err = db.Select(&threads, "SELECT thread_id, title FROM threads")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	userName, ok := c.Get("UserName")

	c.HTML(http.StatusOK, "index.html", gin.H{
		"threads": threads,
		"user_name": userName,
		"user_exists": ok,
	})
}
