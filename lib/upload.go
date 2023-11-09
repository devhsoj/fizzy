package lib

import (
	"encoding/binary"
	"time"
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

func DeserializeUploadEntry(data []byte) Upload {
	var upload Upload

	entryLength := len(data)

	/*
		filename 		string X>=1 b
		id 		 		uint64 8    b
		size 	 		uint64 8    b
		dateUploaded	uint64 8    b
		=
		(X>=1)+8+8+8 = min length 25
	*/

	if entryLength < 25 {
		return upload
	}

	dateUploadedLeBytes := data[entryLength-8:]
	sizeLeBytes := data[entryLength-16 : entryLength-8]
	idLeBytes := data[entryLength-24 : entryLength-16]
	filenameBytes := data[:entryLength-24]

	upload.Id = binary.LittleEndian.Uint64(idLeBytes)
	upload.Filename = string(filenameBytes)
	upload.Size = binary.LittleEndian.Uint64(sizeLeBytes)
	upload.DateUploaded = time.UnixMilli(int64(binary.LittleEndian.Uint64(dateUploadedLeBytes)))

	return upload
}
