package routes

import (
	"gitchecker/controller"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/template/html/v2"
	"html/template"
)

var templateFuncs = template.FuncMap{
	"add": func(a, b int) int {
		return a + b
	},
	"sub": func(a, b int) int {
		return a - b
	},
	"until": func(n int) []int {
		var res []int
		for i := 0; i < n; i++ {
			res = append(res, i)
		}
		return res
	},
}

func PanelHandler() *fiber.App {

	//engine := html.New("/Users/vefgh/Desktop/server/gitchecker/templates/views", ".html")
	engine := html.New("./templates/views/", ".html")
	engine.AddFunc("add", templateFuncs["add"])
	engine.AddFunc("sub", templateFuncs["sub"])
	engine.AddFunc("until", templateFuncs["until"])
	panel := fiber.New(fiber.Config{
		Views: engine,
	})
	panel.Static("/statics", "./templates/views/statics")
	panel.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))
	panel.Get("/", controller.Home)
	panel.Get("/login", controller.ShowLogin)
	panel.Post("/login", controller.HandleLogin)
	panel.Post("/delete", controller.DeleteBranches)

	panel.Get("/cpu", monitor.New())

	return panel
}
