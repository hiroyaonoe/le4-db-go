package user

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hiroyaonoe/le4-db-go/db"
	"github.com/hiroyaonoe/le4-db-go/domain"
	"github.com/hiroyaonoe/le4-db-go/lib/session"
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

	user := domain.User{}
	err = db.Get(&user, "SELECT * FROM users WHERE user_id = $1", userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.String(http.StatusNotFound, fmt.Sprintf("user %d not found", userID))
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	loginUserID := auth.GetUserID(c)
	loginUserRole := auth.GetUserRole(c)

	if loginUserID != userID && loginUserRole != domain.OWNER {
		c.String(http.StatusUnauthorized, "unauthorized")
		return
	}

	tx, err := db.Beginx()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	defer tx.Rollback()

	// ownerが自分を削除する場合はownerが0人になる場合がある
	// ownerは一人以上必要
	if loginUserID == userID && loginUserRole == domain.OWNER {
		var cnt int
		err = tx.Get(&cnt, "SELECT count(*) FROM users WHERE role = 'owner'")
		if cnt <= 1 {
			c.String(http.StatusBadRequest, "owner needs at least 1")
			return
		}
	}

	query := "DELETE FROM threads " +
		"WHERE thread_id IN " +
		"(" +
		"	SELECT thread_id " +
		"	FROM post_threads " +
		"	WHERE user_id = $1 " +
		")" // スレッドの削除(このスレッドに紐づくコメントはON DELETE CASCADEで削除される)
	_, err = tx.Exec(query, userID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	query = "DELETE FROM comments " +
		"WHERE comment_id IN " +
		"(" +
		"	SELECT comment_id " +
		"	FROM post_comments " +
		"	WHERE user_id = $1 " +
		")" // コメントの削除
	_, err = tx.Exec(query, userID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	_, err = tx.Exec("DELETE FROM users WHERE user_ID = $1", userID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	tx.Commit()

	if loginUserID == userID {
		session.Clear(c)
	}

	c.HTML(http.StatusOK, "message.html", gin.H{
		"message": fmt.Sprintf("ユーザー%sを削除しました", user.Name),
	})
	return
}
