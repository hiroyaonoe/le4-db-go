package thread

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hiroyaonoe/le4-db-go/db"
	"github.com/hiroyaonoe/le4-db-go/domain"
)

func Create(c *gin.Context) {
	db, err := db.NewDB()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	thread := domain.Thread{}
	thread.Title = c.PostForm("thread_title")
	thread.UserID = c.GetInt("UserID") // AuthenticateWithRedirectでユーザーの存在確認は済
	thread.CategoryID, err = strconv.Atoi(c.PostForm("category_id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	thread.CreatedAt = domain.NewDateTime(time.Now())

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
	tx.Commit()

	id := strconv.Itoa(thread.ThreadID)
	c.Redirect(http.StatusMovedPermanently, "/thread/"+id)
	return
}
