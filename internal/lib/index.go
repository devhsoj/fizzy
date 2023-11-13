package lib

import (
	"bytes"
	"encoding/binary"
	"io"
	"time"

	"github.com/devhsoj/fizzy/internal/index"
)

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

func WriteUploadToIndex(upload *Upload) (int, error) {
	data := upload.Serialize()

	return index.IndexFile.Write(data)
}

func WriteUploadsToIndex(uploads *[]Upload) (int, error) {
	var data []byte

	for i := 0; i < len(*uploads); i++ {
		data = append(data, (*uploads)[i].Serialize()...)
	}

	return index.IndexFile.Write(data)
}

func ParseIndexUploadEntries() ([]Upload, error) {
	var uploads []Upload

	var offset int64 = 0
	var buf []byte = make([]byte, 16_384)
	var data []byte

	for {
		n, err := index.IndexFile.ReadAt(buf, offset)

		if err != nil && err != io.EOF {
			if err == io.EOF {
				data = append(data, buf[:n]...)
				offset += int64(n)
				break
			}

			return uploads, nil
		}

		if n == 0 {
			break
		}

		data = append(data, buf[:n]...)
		offset += int64(n)
	}

	uploadEntries := bytes.Split(data, []byte("\r"))

	for _, entry := range uploadEntries {
		if len(entry) > 0 {
			uploads = append(uploads, DeserializeUploadEntry(entry))
		}
	}

	return uploads, nil
}
