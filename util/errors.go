package util

import "github.com/gofiber/fiber/v2"

func WrapFiberError(status int, error error) error {
	return fiber.NewError(status, error.Error())
}

func WrapFiberErrorText(status int, error string) error {
	return fiber.NewError(status, error)
}
