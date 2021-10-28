package thread

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hiroyaonoe/le4-db-go/db"
)

func Create(c *gin.Context) {
	db, err := db.NewDB()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	thread := Thread{}
	thread.Title = c.PostForm("thread_title")
	ids := []int{}
	err = db.Select(&ids, "INSERT INTO threads (title) VALUES ($1) RETURNING thread_id", thread.Title)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	thread.ThreadID = ids[0]

	id := strconv.Itoa(thread.ThreadID)
	c.Redirect(http.StatusMovedPermanently, "thread/" + id)
}
