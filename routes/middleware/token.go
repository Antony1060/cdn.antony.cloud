package middleware

import (
	"cdn/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func VerifyToken(token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") != "Bearer " + token {
			util.JsonWithStatus(c, http.StatusForbidden, &gin.H{})
		}
	}
}