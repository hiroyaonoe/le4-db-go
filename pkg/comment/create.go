package comment

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hiroyaonoe/le4-db-go/db"
	"github.com/hiroyaonoe/le4-db-go/domain"
	"github.com/hiroyaonoe/le4-db-go/pkg/auth"
)

func Create(c *gin.Context) {
	db, err := db.NewDB()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	comment := domain.Comment{}
	comment.Content = c.PostForm("comment_content")
	if comment.Content == "" {
		c.String(http.StatusBadRequest, "comment's content cannot be null")
		return
	}
	comment.ThreadID, err = strconv.Atoi(c.Param("thread_id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	comment.UserID = auth.GetUserID(c)
	comment.CreatedAt = domain.NewDateTime(time.Now())

	tx, err := db.Beginx()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	defer tx.Rollback()

	err = tx.Get(&comment.CommentID, "INSERT INTO comments (content, thread_id) VALUES ($1, $2) RETURNING comment_id", comment.Content, comment.ThreadID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	_, err = tx.NamedExec("INSERT INTO post_comments (comment_id, thread_id, user_id, created_at) VALUES (:comment_id, :thread_id, :user_id, :created_at)", comment)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	tx.Commit()

	id := strconv.Itoa(comment.ThreadID)
	c.Redirect(http.StatusSeeOther, "/thread/"+id)
	return
}
