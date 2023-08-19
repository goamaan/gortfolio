package main

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/gofiber/fiber/v2"
)

type AuthBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	IsAuthed bool `json:"isAuthed"`
}

type CookieValue struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

func Login(c *fiber.Ctx) error {
	var authBody AuthBody

	if err := c.BodyParser(&authBody); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	user := os.Getenv("ADMIN_USERNAME")
	password := os.Getenv("ADMIN_PASSWORD")

	if authBody.Username != user || authBody.Password != password {
		return c.Status(401).JSON(errors.New("Invalid username/password"))
	}

	cookieValue, err := json.Marshal(&CookieValue{Username: user, Role: "admin"})
	if err != nil {
		return c.Status(500).JSON(errors.New(""))
	}

	c.Cookie(&fiber.Cookie{
		Name:     "goamaan_session",
		Value:    string(cookieValue),
		HTTPOnly: true,
		SameSite: "Lax"})

	return c.Redirect("/")
}

func Logout(c *fiber.Ctx) error {
	c.ClearCookie("goamaan_session")
	return c.Redirect("/")
}
