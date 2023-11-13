package lib

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/cespare/xxhash"
	"github.com/gofiber/fiber/v2"
)

type Upload struct {
	Id           uint64
	Filename     string
	Size         uint64
	DateUploaded time.Time
}

func (upload *Upload) Serialize() []byte {
	/*
		serializaton format:
		\r{Filename}{Id}{Size}{DateUploaded}

		example entry:
		\rtest.txtppppppppiiiiiiiiyyyyyyyy
	*/

	data := []byte("\r" + upload.Filename)

	leBytes := make([]byte, 8)

	binary.LittleEndian.PutUint64(leBytes, upload.Id)
	data = append(data, leBytes...)

	binary.LittleEndian.PutUint64(leBytes, upload.Size)
	data = append(data, leBytes...)

	binary.LittleEndian.PutUint64(leBytes, uint64(upload.DateUploaded.UnixMilli()))
	data = append(data, leBytes...)

	return data
}

func StoreUploadsFromRequest(c *fiber.Ctx) error {
	form, err := c.MultipartForm()

	if err != nil {
		return err
	}

	files := form.File["files"]

	var uploads []Upload

	for _, file := range files {
		now := time.Now()
		id := xxhash.Sum64String(fmt.Sprintf("%s-%d-%d", file.Filename, file.Size, now.UnixMilli()))

		uploads = append(uploads, Upload{
			Id:           id,
			Filename:     file.Filename,
			Size:         uint64(file.Size),
			DateUploaded: now,
		})

		if err := c.SaveFile(file, fmt.Sprintf("./static/%d", id)); err != nil {
			return err
		}
	}

	if _, err := WriteUploadsToIndex(&uploads); err != nil {
		return err
	}

	return nil
}
