package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"repo.polytron.co.id/rd-special-project/auth-service/backend/configs"
)

func main() {
	gin.SetMode(configs.GIN_MODE)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello Auth Service",
		})
	})
	fmt.Println("BE PORT: ", configs.BE_PORT)
	if err := r.Run(fmt.Sprintf(":%s", configs.BE_PORT)); err != nil {
		panic(err)
	}
}
