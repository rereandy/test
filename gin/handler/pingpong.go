package handler

import "github.com/gin-gonic/gin"

func PingHandler() gin.HandlerFunc {
	counter := 0
	return func(c *gin.Context) {
		counter++
		c.JSON(200, gin.H{
			"message": "pong",
			"counter": counter,
		})
	}
}
