package main

import "github.com/gin-gonic/gin"

func main() {
	// Создаем новый роутер Gin
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Привет я Gin!",
		})
	})

	router.Run(":8101")
}
