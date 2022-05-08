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

type removeOpts struct {
	Name  string `json:"name"`
	Index bool   `json:"index"`
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

		if err = h.CreateFile(f, args.Name, args.Password, args.TimeTillDeath, args.Index, args.Override); err != nil {
			return err
		}

		return util.Status(c, http.StatusOK)
	}
}

func (h *Handler) RemoveFile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		args := removeOpts{}

		if err := c.BodyParser(&args); err != nil {
			return util.WrapFiberError(http.StatusBadRequest, err)
		}

		if !h.Exists(args.Name, args.Index) {
			return util.WrapFiberErrorText(http.StatusNotFound, "file not found")
		}

		f, err := h.Get(args.Name, args.Index)
		if err != nil {
			return err
		}

		if err = f.Delete(); err != nil {
			return err
		}

		return util.Status(c, http.StatusOK)
	}
}

func (h *Handler) GetFiles() fiber.Handler {
	return func(c *fiber.Ctx) error {
		all, err := h.GetAll()
		if err != nil {
			return err
		}

		return util.JsonWithStatus(c, http.StatusOK, all)
	}
}
