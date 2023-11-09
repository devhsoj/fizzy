package main

import (
	"log"
	"os"

	"github.com/devhsoj/fizzy/index"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/handlebars/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print("Failed to load variables from .env!")
	}

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

	var listenAddress = os.Getenv("LISTEN_ADDRESS")

	if len(listenAddress) == 0 {
		listenAddress = "localhost:3000"
	}

	log.Printf("starting @ %s", listenAddress)

	app.Listen(listenAddress)
}
