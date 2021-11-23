package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hiroyaonoe/le4-db-go/db"
	"github.com/hiroyaonoe/le4-db-go/domain"
	"github.com/hiroyaonoe/le4-db-go/pkg/auth"
)

func UpdateRole(c *gin.Context) {
	db := db.GetDB()

	if role := auth.GetUserRole(c); role != domain.OWNER {
		c.String(http.StatusUnauthorized, "you are not owner")
		return
	}

	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	newRole := c.PostForm("new_role")
	if newRole != domain.MEMBER && newRole != domain.ADMIN && newRole != domain.OWNER {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	tx, err := db.Beginx()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	defer tx.Rollback()

	if userID == auth.GetUserID(c) { // ownerは一人以上必要
		var cnt int
		err = tx.Get(&cnt, "SELECT count(*) FROM users WHERE role = 'owner'")
		if cnt <= 1 {
			c.String(http.StatusBadRequest, "owner needs at least 1")
			return
		}
	}

	_, err = tx.Exec("UPDATE users SET role = $1 WHERE user_ID = $2", newRole, userID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	tx.Commit()

	c.Redirect(http.StatusSeeOther, "/user/"+userIDStr)
	return
}
