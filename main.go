package main

import (
	"cdn/db"
	"cdn/env"
	"cdn/routes"
	"cdn/routes/middleware"
	"github.com/apex/log"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
)

func main() {
	err := generateDirs()
	if err != nil {
		log.WithError(err).Fatal("Failed to created needed directories")
		return
	}
	err = db.Init()
	if err != nil {
		log.WithError(err).Fatal("Failed to create database")
		return
	}
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

	router.Use(static.Serve("/", static.LocalFile("./files/index", true)))

	api := router.Group("/api")
	{
		api.POST("/add", middleware.VerifyToken(env.Token), routes.AddFile())
		api.POST("/delete", middleware.VerifyToken(env.Token), routes.RemoveFile())
		api.GET("/get", middleware.VerifyToken(env.Token), routes.GetFiles())
	}

	log.Infof("Loaded %d routes", len(router.Routes()))
}

func startServer(env *env.EnvConfig) {
	if env.Mode == "DEBUG" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	g := gin.New()
	loadRoutes(g, env)
	log.WithError(g.Run(":" + env.Port)).Fatal("HTTP server ended")
}

func generateDirs() error {
	err := os.MkdirAll(filepath.Join("files", "index"), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}