package user

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

	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}
	
	loginUserID := auth.GetUserIDInt(c)
	loginUserRole := auth.GetUserRole(c)

	if loginUserID != userID && loginUserRole != domain.OWNER {
		c.String(http.StatusUnauthorized, "unauthorized")
		return
	}

	// ownerが自分を削除する場合はownerが0人になる場合がある
	// ownerは一人以上必要
	if loginUserID == userID && loginUserRole == domain.OWNER {
		var cnt int
		err = db.Get(&cnt, "SELECT count(*) FROM users WHERE role = 'owner'")
		if cnt <= 1 {
			c.String(http.StatusBadRequest, "owner needs at least 1")
			return
		}
	}

	_, err = db.Exec("DELETE FROM users WHERE user_ID = $1", userID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/user/"+userIDStr)
	return
}
