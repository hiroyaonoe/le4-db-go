package comment

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hiroyaonoe/le4-db-go/db"
	"github.com/hiroyaonoe/le4-db-go/pkg/datetime"
)

func Create(c *gin.Context) {
	db, err := db.NewDB()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	comment := Comment{}
	comment.Content = c.PostForm("comment_content")
	comment.ThreadID, err = strconv.Atoi(c.Param("thread_id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	comment.UserID = c.GetInt("UserID") // AuthenticateWithRedirectでユーザーの存在確認は済
	comment.CreatedAt = datetime.NewDateTime(time.Now())
	ids := []int{}

	tx, err := db.Beginx()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	defer tx.Rollback()

	err = tx.Select(&ids, "INSERT INTO comments (content, thread_id) VALUES ($1, $2) RETURNING comment_id", comment.Content, comment.ThreadID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	comment.CommentID = ids[0]
	_, err = tx.NamedExec("INSERT INTO post_comments (comment_id, thread_id, user_id, created_at) VALUES (:comment_id, :thread_id, :user_id, :created_at)", comment)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	tx.Commit()

	id := strconv.Itoa(comment.ThreadID)
	c.Redirect(http.StatusMovedPermanently, "/thread/"+id)
}
