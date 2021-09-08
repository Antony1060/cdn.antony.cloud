package routes

import (
	"cdn/util"
	"github.com/gin-gonic/gin"
)

func AddFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		f, err := c.FormFile("file")
		if err != nil {
			util.JsonWithStatus(c, 400, &gin.H{
				"error": "File not found",
			})
			return
		}

		if c.PostForm("index") == "true" {
			err := c.SaveUploadedFile(f, "./files/index/" + f.Filename)
			if err != nil {
				util.JsonWithStatus(c, 500, &gin.H{
					"error": "Couldn't save file: " + err.Error(),
				})
				return
			}
			util.Status(c, 201)
			return
		}

		util.JsonWithStatus(c, 200, &gin.H{})
	}
}

func RemoveFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		util.JsonWithStatus(c, 200, &gin.H{})
	}
}

func GetFiles() gin.HandlerFunc {
	return func(c *gin.Context) {
		util.JsonWithStatus(c, 200, &gin.H{})
	}
}