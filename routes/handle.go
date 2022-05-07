package routes

import (
	"cdn/filesystem"
	"cdn/util"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type addOpts struct {
	Index         bool   `json:"index" default:"false"`
	TimeTillDeath int    `json:"timeTillDeath" default:"-1"`
	Name          string `json:"name" default:""`
	Password      string `json:"password" default:""`
	Override      bool   `json:"override" default:"false"`
}

type Handler struct {
	*filesystem.FileSystem
}

func (h *Handler) AddFile() fiber.Handler {
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

func (h *Handler) RemoveFile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return util.Status(c, http.StatusOK)
	}
}

func (h *Handler) GetFiles() fiber.Handler {
	return func(c *fiber.Ctx) error {
		lookup, err := h.GetAll()
		if err != nil {
			return err
		}

		return util.JsonWithStatus(c, http.StatusOK, lookup)
	}
}
