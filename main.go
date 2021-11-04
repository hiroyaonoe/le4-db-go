package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/hiroyaonoe/le4-db-go/config"
	"github.com/hiroyaonoe/le4-db-go/pkg/auth"
	"github.com/hiroyaonoe/le4-db-go/pkg/comment"
	"github.com/hiroyaonoe/le4-db-go/pkg/index"
	"github.com/hiroyaonoe/le4-db-go/pkg/login"
	"github.com/hiroyaonoe/le4-db-go/pkg/logout"
	"github.com/hiroyaonoe/le4-db-go/pkg/signup"
	"github.com/hiroyaonoe/le4-db-go/pkg/thread"
	"github.com/hiroyaonoe/le4-db-go/pkg/user"
)

func main() {
	e := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	e.Use(sessions.Sessions("mysession", store))

	e.LoadHTMLGlob("view/*.html")

	e.Use(auth.Authenticate)

	e.GET("", index.Get)

	e.GET("/login", login.Get)
	e.POST("/login", login.Post)
	e.GET("/signup", signup.Get)
	e.GET("/logout", logout.Get)

	u := e.Group("/user")
	u.POST("", user.Create)
	u.GET("/:user_id", user.Get)

	th := e.Group("/thread")
	th.GET("/:thread_id", thread.Get)
	thAuth := e.Group("/thread", auth.AuthenticateWithRedirect)
	thAuth.POST("", thread.Create)
	
	coAuth := e.Group("/thread/:thread_id/comment", auth.AuthenticateWithRedirect)
	coAuth.POST("", comment.Create)

	e.Run(config.Port())
}
