package logout

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hiroyaonoe/le4-db-go/lib/session"
)

func Get(c *gin.Context) {
	session.Clear(c)
	c.HTML(http.StatusOK, "message.html", gin.H{
		"message": "ログアウトしました",
	})
	return
}
