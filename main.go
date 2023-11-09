package main

import (
	"log"
	"os"

	"github.com/devhsoj/fizzy/index"
	"github.com/devhsoj/fizzy/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/handlebars/v2"
	"github.com/joho/godotenv"
)

func main() {
	// try to load environment variables from .env
	if err := godotenv.Load(); err != nil {
		log.Print("Failed to load variables from .env!")
	}

	// create the index file or load the file pointer into memory
	if err := index.SetupIndexFile(); err != nil {
		panic(err)
	}

	// setup views
	viewEngine := handlebars.New("./views/", ".hbs")

	// setup app
	app := fiber.New(fiber.Config{
		Views:                 viewEngine,
		DisableStartupMessage: true,
		BodyLimit:             10 * 1_024 * 1_024 * 1_024, // 10 Gib
	})

	// setup static dir
	app.Static("/static", "./static/")

	// setup routes
	app.Get("/", routes.IndexRoute)
	app.Post("/upload", routes.UploadRoute)

	// determine a listening address for the app
	var listenAddress = os.Getenv("LISTEN_ADDRESS")

	if len(listenAddress) == 0 {
		listenAddress = "localhost:3000"
	}

	log.Printf("starting @ %s", listenAddress)

	// start app
	app.Listen(listenAddress)
}
