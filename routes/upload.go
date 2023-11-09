package routes

import "github.com/gofiber/fiber/v2"

func UploadRoute(c *fiber.Ctx) error {
	return c.SendStatus(200)
}
