package util

import "github.com/gofiber/fiber/v2"

func JsonWithStatus(c *fiber.Ctx, status int, data any) error {
	return c.Status(status).JSON(fiber.Map{
		"status": status,
		"data":   data,
	})
}

func ErrorWithStatus(c *fiber.Ctx, status int, error error) error {
	return c.Status(status).JSON(fiber.Map{
		"status": status,
		"error":  error.Error(),
	})
}

func Status(c *fiber.Ctx, status int) error {
	return c.Status(status).JSON(fiber.Map{
		"status": status,
	})
}
