package logout

import (
	"net/http"

	"github.com/hiroyaonoe/le4-db-go/lib/session"
	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context) {
	session.Clear(c)
	c.HTML(http.StatusOK, "logout.html", gin.H{})
	return
}
