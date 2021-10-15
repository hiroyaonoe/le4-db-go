package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hiroyaonoe/le4db-go/pkg/user"
	"github.com/hiroyaonoe/le4db-go/config"
)

func main() {
	e := gin.Default()
	e.LoadHTMLGlob("view/*.html")

	e.GET("/user", user.Get)

	e.Run(config.Port())
}
