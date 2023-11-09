package lib

import (
	"encoding/binary"
	"time"
)

type Upload struct {
	Id           uint64
	Name         string
	Size         uint64
	DateUploaded time.Time
}

func (upload *Upload) Serialize() []byte {
	data := []byte("\r" + upload.Name)

	leBytes := make([]byte, 8)

	binary.LittleEndian.PutUint64(leBytes, upload.Id)
	data = append(data, leBytes...)

	binary.LittleEndian.PutUint64(leBytes, upload.Size)
	data = append(data, leBytes...)

	binary.LittleEndian.PutUint64(leBytes, uint64(upload.DateUploaded.UnixMilli()))
	data = append(data, leBytes...)

	return data
}

func DeserializeUploadEntry(data []byte) Upload {
	var upload Upload

	entryLength := len(data)

	if entryLength < 8+8+8+1 {
		return upload
	}

	dateUploadedLeBytes := data[entryLength-8:]
	sizeLeBytes := data[entryLength-16 : entryLength-8]
	idLeBytes := data[entryLength-24 : entryLength-16]
	nameBytes := data[:entryLength-24]

	upload.Id = binary.LittleEndian.Uint64(idLeBytes)
	upload.Name = string(nameBytes)
	upload.Size = binary.LittleEndian.Uint64(sizeLeBytes)
	upload.DateUploaded = time.UnixMilli(int64(binary.LittleEndian.Uint64(dateUploadedLeBytes)))

	return upload
}
