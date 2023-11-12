package routes

import (
	"log"

	"github.com/devhsoj/fizzy/internal/lib"
	"github.com/gofiber/fiber/v2"
)

func Upload(c *fiber.Ctx) error {
	if err := lib.StoreUploadsFromRequest(c); err != nil {
		log.Println(err)
		return c.SendStatus(500)
	}

	return c.Redirect("/")
}
