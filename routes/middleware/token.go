package middleware

import (
	"cdn/util"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func VerifyToken(token string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if val, ok := c.GetReqHeaders()["Authorization"]; !ok || val != "Bearer "+token {
			return util.Status(c, http.StatusForbidden)
		}

		return c.Next()
	}
}
