package main

import (
	"log"

	"github.com/devhsoj/fizzy/index"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/handlebars/v2"
)

func main() {
	if err := index.SetupIndexFile(); err != nil {
		panic(err)
	}

	viewEngine := handlebars.New("./views/", ".hbs")

	app := fiber.New(fiber.Config{
		Views:                 viewEngine,
		DisableStartupMessage: true,
	})

	app.Static("/static", "./static/")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})

	app.Post("/upload", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	log.Println("starting @ http://localhost:3000")

	app.Listen("localhost:3000")
}
