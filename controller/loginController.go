package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var store = session.New()

func ShowLogin(ctx *fiber.Ctx) error {
	return ctx.Render("login", fiber.Map{})
}

func HandleLogin(ctx *fiber.Ctx) error {
	token := ctx.FormValue("token")
	url := ctx.FormValue("url")

	sess, err := store.Get(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Session error")
	}

	sess.Set("token", token)
	sess.Set("url", url)
	sess.Save()

	return ctx.Redirect("/")
}
