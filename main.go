package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hiroyaonoe/le4-db-go/config"
	"github.com/hiroyaonoe/le4-db-go/pkg/index"
)

func main() {
	e := gin.Default()
	e.LoadHTMLGlob("view/*.html")

	e.GET("/", index.Get)

	e.Run(config.Port())
}
