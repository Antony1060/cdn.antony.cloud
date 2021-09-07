package util

import "github.com/gin-gonic/gin"

func JsonWithStatus(c *gin.Context, status int, json *gin.H) {
	c.JSON(status, &gin.H{
		"status": status,
		"data": json,
	})
}