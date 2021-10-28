package thread

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hiroyaonoe/le4-db-go/db"
)

func Get(c *gin.Context) {
	db, err := db.NewDB()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	threadID, err := strconv.Atoi(c.Param("thread_id"))
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
	}

	threads := []Thread{}
	err = db.Select(&threads, "SELECT thread_id, title FROM threads WHERE thread_id = $1", threadID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	thread := threads[0]

	c.HTML(http.StatusOK, "thread_get.html", gin.H{
		"thread": thread,
	})
}
