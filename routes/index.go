package routes

import "github.com/gofiber/fiber/v2"

func IndexRoute(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{})
}
