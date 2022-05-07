package routes

import (
	"cdn/util"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type addOpts struct {
	Index bool `json:"index"`
}

func AddFile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		f, err := c.FormFile("file")
		if err != nil {
			return util.WrapFiberError(http.StatusBadRequest, err)
		}

		args := addOpts{}

		if err = c.BodyParser(&args); err != nil {
			return util.WrapFiberError(http.StatusBadRequest, err)
		}

		if args.Index {
			err := c.SaveFile(f, "./files/index/"+f.Filename)
			if err != nil {
				return util.WrapFiberError(http.StatusInternalServerError, err)
			}
			return util.Status(c, http.StatusCreated)
		}

		return util.Status(c, http.StatusOK)
	}
}

func RemoveFile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return util.Status(c, http.StatusOK)
	}
}

func GetFiles() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return util.Status(c, http.StatusOK)
	}
}
