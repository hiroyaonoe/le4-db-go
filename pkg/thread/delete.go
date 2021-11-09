package thread

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

	thread := domain.Thread{}
	thread.ThreadID, err = strconv.Atoi(c.Param("thread_id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	query := "SELECT thread_id, user_id " +
		"FROM post_threads " +
		"WHERE thread_id = $1"
	err = db.Get(&thread, query, thread.ThreadID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	userID := auth.GetUserID(c)
	userRole := auth.GetUserRole(c)
	if userID != thread.UserID && userRole != domain.ADMIN && userRole != domain.OWNER {
		c.String(http.StatusUnauthorized, "unauthorized")
		return
	}

	_, err = db.Exec("DELETE FROM threads WHERE thread_id = $1", thread.ThreadID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/")
	return
}
