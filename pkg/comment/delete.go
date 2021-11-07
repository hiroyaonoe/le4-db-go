package comment

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hiroyaonoe/le4-db-go/db"
	"github.com/hiroyaonoe/le4-db-go/domain"
	"github.com/hiroyaonoe/le4-db-go/pkg/auth"
)

func Delete(c *gin.Context) {
	db, err := db.NewDB()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	comment := domain.Comment{}
	comment.CommentID, err = strconv.Atoi(c.Param("comment_id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	comment.ThreadID, err = strconv.Atoi(c.Param("thread_id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	query := "SELECT comment_id, thread_id, user_id " +
		"FROM post_comments " +
		"WHERE comment_id = $1 " +
		"	AND thread_id = $2"
	err = db.Get(&comment, query, comment.CommentID, comment.ThreadID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	userID := auth.GetUserID(c)
	userRole := auth.GetUserRole(c)
	if userID != comment.UserID && userRole != domain.ADMIN && userRole != domain.OWNER {
		c.String(http.StatusUnauthorized, "unauthorized")
		return
	}

	_, err = db.Exec("DELETE FROM comments WHERE comment_id = $1 AND thread_id = $2", comment.CommentID, comment.ThreadID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	id := strconv.Itoa(comment.ThreadID)
	c.Redirect(http.StatusSeeOther, "/thread/"+id)
	return
}
