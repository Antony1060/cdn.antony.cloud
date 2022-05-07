package middleware

import (
	"cdn/util"
	"fmt"
	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"
	"time"
)

func DefaultLogger(c *fiber.Ctx) error {
	start := time.Now()
	err := c.Next()

	duration := time.Since(start)
	latencyFormatted := ""
	if duration.Microseconds() > 1000 {
		latencyFormatted = fmt.Sprintf("%dms", duration.Milliseconds())
	} else {
		latencyFormatted = fmt.Sprintf("%dµs", duration.Microseconds())
	}

	f := log.Fields{
		"client":  c.IP(),
		"status":  c.Context().Response.StatusCode(),
		"latency": latencyFormatted,
	}

	if err != nil {
		f["error"] = err
	}

	log.WithFields(f).Debugf("%s %s", util.MethodColor(c.Method())+" "+c.Method()+" "+util.ResetColor(), c.Path())

	return err
}
