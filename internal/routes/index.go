package routes

import (
	"log"

	"github.com/devhsoj/fizzy/internal/lib"
	"github.com/gofiber/fiber/v2"
)

func IndexPage(c *fiber.Ctx) error {
	uploads, err := lib.ParseIndexEntries()

	if err != nil {
		log.Println(err)
		return c.SendStatus(500)
	}

	return c.Render("index", fiber.Map{
		"uploads":     uploads,
		"uploadCount": len(uploads),
	})
}
