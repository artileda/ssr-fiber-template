package configs

import (
	fiber "github.com/gofiber/fiber/v2"
)

func PublicRoutes(app *fiber.App) *fiber.App {

	app.Get("/", func(c *fiber.Ctx) error {
		// Render index
		return c.Render("index", fiber.Map{
			"Title": "Hello, World!",
		})
	})

	return app
}

func PrivateRoute() {

}
