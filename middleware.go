package main

import (
	"encoding/json"
	"os"

	"github.com/gofiber/fiber/v2"
)

const LocalsAdminkey = "admin"

func ParseCookieMiddleware(c *fiber.Ctx) error {
	sessionCookie := c.Cookies("goamaan_session")
	if sessionCookie == "" {
		return c.Next()
	}
	var parsed CookieValue

	err := json.Unmarshal([]byte(sessionCookie), &parsed)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if parsed.Username != os.Getenv("ADMIN_USERNAME") || parsed.Role != "admin" {
		c.ClearCookie("goamaan_session")
	} else {
		c.Locals(LocalsAdminkey, &parsed)
	}

	return c.Next()
}
