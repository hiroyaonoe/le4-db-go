package thread

import (
	"net/http"
	"strconv"
	"time"

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
	thread.UserID = c.GetInt("UserID") // AuthenticateWithRedirectでユーザーの存在確認は済
	thread.CreatedAt = time.Now()
	ids := []int{}

	tx, err := db.Beginx()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	defer tx.Rollback()

	err = tx.Select(&ids, "INSERT INTO threads (title) VALUES ($1) RETURNING thread_id", thread.Title)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	thread.ThreadID = ids[0]
	_, err = tx.NamedExec("INSERT INTO post_threads (thread_id, user_id, created_at) VALUES (:thread_id, :user_id, :created_at)", thread)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	tx.Commit()

	id := strconv.Itoa(thread.ThreadID)
	c.Redirect(http.StatusMovedPermanently, "thread/"+id)
}
