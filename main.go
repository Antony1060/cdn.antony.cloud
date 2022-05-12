package main

import (
	"cdn/db"
	"cdn/env"
	"cdn/filesystem"
	"cdn/routes"
	"cdn/routes/middleware"
	"cdn/util"
	"encoding/base64"
	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"net/http"
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
		relativePath := c.Path()
		if !strings.HasPrefix(relativePath, "/secret") {
			relativePath = "/index" + relativePath
		}

		fullPath, err := filepath.Abs("./files" + relativePath) // ikik, hard-coded values, will maybe fix later :/
		if err != nil {
			return true
		}

		target, ok := db.Get().FilePasswords[fullPath]
		if !ok || target == "" {
			return false
		}

		val, ok := c.GetReqHeaders()["Authorization"]
		if !ok || !strings.HasPrefix(val, "Basic ") {
			c.Append("WWW-Authenticate", "Basic realm=Forbidden")
			_ = util.Status(c, http.StatusUnauthorized)
			return true
		}

		basicAuth := strings.TrimPrefix(val, "Basic ")

		decoded, err := base64.StdEncoding.DecodeString(basicAuth)
		if err != nil {
			return true
		}

		s := strings.SplitAfterN(string(decoded), ":", 2)
		if len(s) < 2 {
			c.Append("WWW-Authenticate", "Basic realm=Forbidden")
			_ = util.Status(c, http.StatusUnauthorized)
			return true
		}

		pass := s[1]
		if bcrypt.CompareHashAndPassword([]byte(target), []byte(pass)) != nil {
			c.Append("WWW-Authenticate", "Basic realm=Forbidden")
			_ = util.Status(c, http.StatusUnauthorized)
			return true
		}

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

	fs, err := filesystem.NewFileSystem("./files/index", "./files/secret")
	if err != nil {
		panic("can't load filesystem") // eh :/
	}

	h := routes.Handler{FileSystem: fs}

	api := app.Group("/api")
	{
		api.Post("/add", middleware.VerifyToken(env.Token), h.AddFile())
		api.Post("/delete", middleware.VerifyToken(env.Token), h.RemoveFile())
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

			if c.Context().Response.StatusCode() == http.StatusUnauthorized {
				code = http.StatusUnauthorized
			} else if e, ok := err.(*fiber.Error); ok {
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
