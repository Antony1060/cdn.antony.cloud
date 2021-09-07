package main

import (
	"cdn/env"
	"cdn/routes"
	"cdn/routes/middleware"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
)

func main() {
	startServer(env.New())
}

func loadRoutes(router *gin.Engine, env *env.EnvConfig) {
	// logging
	router.Use(gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		log.WithFields(log.Fields{
			"client": params.ClientIP,
			"status": params.StatusCode,
			"latency": params.Latency,
		}).Debugf("%s %s", params.MethodColor() + params.Method + params.ResetColor(), params.Path)

		return ""
	}))

	router.POST("/add", middleware.VerifyToken(env.Token), routes.AddFile())
	router.POST("/delete", middleware.VerifyToken(env.Token), routes.RemoveFile())
	router.GET("/get", middleware.VerifyToken(env.Token), routes.GetFiles())
}

func startServer(env *env.EnvConfig) {
	g := gin.New()
	loadRoutes(g, env)
	log.WithError(g.Run(":" + env.Port)).Fatal("HTTP server ended")
}
