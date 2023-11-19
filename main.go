package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/devhsoj/fizzy/internal/index"
	"github.com/devhsoj/fizzy/internal/routes"
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
	viewEngine := handlebars.New("./web/views/", ".hbs")

	viewEngine.AddFunc("formatSize", func(size uint64) string {
		// modified from clever solution @ https://yourbasic.org/golang/formatting-byte-size-to-human-readable-format/
		const unit = 1000

		if size < unit {
			return fmt.Sprintf("%d b", size)
		}

		div, exp := int64(unit), 0

		for n := size / unit; n >= unit; n /= unit {
			div *= unit
			exp++
		}

		return fmt.Sprintf("%.2f %cb", float64(size)/float64(div), "kmgtpe"[exp])
	})

	viewEngine.AddFunc("formatDate", func(date time.Time) string {
		return date.Format(time.RFC822)
	})

	// setup app
	app := fiber.New(fiber.Config{
		Views:                 viewEngine,
		DisableStartupMessage: true,
		BodyLimit:             10 * 1_024 * 1_024 * 1_024, // 10 GiB
	})

	// setup static dir
	app.Static("/static", "./web/static/")

	// setup routes
	app.Get("/", routes.IndexPage)
	app.Post("/upload", routes.Upload)

	// determine a listening address for the app
	var listenAddress = os.Getenv("LISTEN_ADDRESS")

	if len(listenAddress) == 0 {
		listenAddress = "localhost:3000"
	}

	log.Printf("starting @ %s", listenAddress)

	// start app
	app.Listen(listenAddress)
}
