package comment

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hiroyaonoe/le4-db-go/db"
	"github.com/hiroyaonoe/le4-db-go/domain"
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
		"FROM comments " +
		"NATURAL JOIN post_comments " +
		"WHERE comment_id = $1 " +
		"	AND thread_id = $2"
	err = db.Get(&comment, query, comment.CommentID, comment.ThreadID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	userID := c.GetInt("UserID")
	userRole, _ := c.Get("UserRole")
	if userID != comment.UserID && userRole != "admin" && userRole != "owner" {
		log.Printf("userID: %d, comment.UserID: %d, userRole: %s\n", userID, comment.UserID, userRole)
		log.Printf("%#v\n", comment)
		c.String(http.StatusUnauthorized, "unauthorized")
		return
	}

	_, err = db.Exec("DELETE FROM comments WHERE comment_id = $1 AND thread_id = $2", comment.CommentID, comment.ThreadID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	id := strconv.Itoa(comment.ThreadID)
	c.Redirect(http.StatusMovedPermanently, "/thread/"+id)
	return
}
