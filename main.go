package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hiroyaonoe/le4-db-go/config"
	"github.com/hiroyaonoe/le4-db-go/pkg/index"
	"github.com/hiroyaonoe/le4-db-go/pkg/thread"
	"github.com/hiroyaonoe/le4-db-go/pkg/auth"
)

func main() {
	e := gin.Default()
	e.LoadHTMLGlob("view/*.html")

	e.GET("", index.Get)

	e.GET("/signup", auth.Signup)

	th := e.Group("/thread")
	th.POST("", thread.Create)
	th.GET("/:thread_id", thread.Get)

	e.Run(config.Port())
}
