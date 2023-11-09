package routes

import (
	"github.com/devhsoj/fizzy/index"
	"github.com/gofiber/fiber/v2"
)

func IndexRoute(c *fiber.Ctx) error {
	uploads, err := index.ParseIndexEntries()

	if err != nil {
		return c.SendStatus(500)
	}

	return c.Render("index", fiber.Map{
		"uploads":     uploads,
		"uploadCount": len(uploads),
	})
}
