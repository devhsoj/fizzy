package routes

import (
	"fmt"
	"log"
	"time"

	"github.com/cespare/xxhash/v2"
	"github.com/devhsoj/fizzy/index"
	"github.com/devhsoj/fizzy/lib"
	"github.com/gofiber/fiber/v2"
)

func UploadRoute(c *fiber.Ctx) error {
	form, err := c.MultipartForm()

	if err != nil {
		log.Println(err)
		return c.SendStatus(500)
	}

	files := form.File["files"]

	// i really don't like having logic like this in a route, but since this is one of only two total routes, i can live with it :)

	for _, file := range files {
		now := time.Now()
		id := xxhash.Sum64String(fmt.Sprintf("%s-%d-%d", file.Filename, file.Size, now.UnixMilli()))

		upload := lib.Upload{
			Id:           id,
			Filename:     file.Filename,
			Size:         uint64(file.Size),
			DateUploaded: now,
		}

		if err := c.SaveFile(file, fmt.Sprintf("./static/%d", id)); err != nil {
			log.Println(err)
			return c.SendStatus(500)
		}

		if _, err = index.WriteUploadToIndex(&upload); err != nil {
			log.Println(err)
			return c.SendStatus(500)
		}
	}

	return c.Redirect("/")
}
