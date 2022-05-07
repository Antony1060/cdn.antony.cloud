package main

import (
	"cdn/db"
	"cdn/env"
	"cdn/filesystem"
	"cdn/routes"
	"cdn/routes/middleware"
	"cdn/util"
	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"
	"os"
	"path/filepath"
	"strings"
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

func loadRoutes(app *fiber.App, env *env.EnvConfig) {
	app.Use(middleware.DefaultLogger)

	nextFunc := func(c *fiber.Ctx) bool {
		// TODO: handle password protected
		return false
	}

	// ./files/index files are shown in the public browsable index
	app.Static("/", "./files/index", fiber.Static{
		Browse: true,
		Next:   nextFunc,
	})

	// files here are not shown in the browsable index
	// this is intended as a default location for file if `index` is not set
	//   Note for CLI: if not parameters are set, file should be uploaded to here with a random hash
	app.Static("/secret", "./files/secret", fiber.Static{
		Next: nextFunc,
	})

	h := routes.Handler{
		FileSystem: &filesystem.FileSystem{
			IndexDir:  "./files/index",
			SecretDir: "./files/secret",
		},
	}

	api := app.Group("/api")
	{
		// TODO: error if already exists, else write if override is present
		api.Post("/add", middleware.VerifyToken(env.Token), h.AddFile())
		// TODO: delete file
		api.Post("/delete", middleware.VerifyToken(env.Token), h.RemoveFile())
		// TODO: list all files
		api.Get("/get", middleware.VerifyToken(env.Token), h.GetFiles())
	}
}

func startServer(env *env.EnvConfig) {
	if strings.ToLower(env.Mode) == "debug" {
		log.SetLevel(log.DebugLevel)
	}

	app := fiber.New(fiber.Config{
		ServerHeader: "go-fiber",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			return util.ErrorWithStatus(c, code, err)
		},
	})
	loadRoutes(app, env)
	log.WithError(app.Listen(":" + env.Port)).Fatal("HTTP server ended")
}

func generateDirs() error {
	err := os.MkdirAll(filepath.Join("files", "index"), os.ModePerm)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Join("files", "secret"), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
