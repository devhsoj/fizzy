package index

import (
	"bytes"
	"io"
	"os"

	"github.com/devhsoj/fizzy/lib"
)

/*
this indexing implementation only exists because i didn't want to deal with setting up cgo/mingw/gcc/go-sqlite3 on windows lol.
i'm on windows, yes i know, don't make fun of me...
i just play too many video games to not have my personal pc using windows.
anyways, this implementation is pretty fast, probably not sqlite fast, but it gets the job done while keeping dependencies super light.
*/

const INDEX_FILENAME = "fizzy.idx"

var IndexFile *os.File

func SetupIndexFile() error {
	var err error

	IndexFile, err = os.OpenFile(INDEX_FILENAME, os.O_RDWR|os.O_APPEND, 0700)

	if err != nil {
		if os.IsNotExist(err) {
			if IndexFile, err = os.Create(INDEX_FILENAME); err != nil {
				return err
			}
		}

		return err
	}

	return nil
}

func WriteUploadToIndex(upload *lib.Upload) (int, error) {
	data := upload.Serialize()

	return IndexFile.Write(data)
}

func WriteUploadsToIndex(uploads *[]lib.Upload) (int, error) {
	var data []byte

	for i := 0; i < len(*uploads); i++ {
		data = append(data, (*uploads)[i].Serialize()...)
	}

	return IndexFile.Write(data)
}

func ParseIndexEntries() ([]lib.Upload, error) {
	var uploads []lib.Upload

	var offset int64 = 0
	var buf []byte = make([]byte, 16_384)
	var data []byte

	for {
		n, err := IndexFile.ReadAt(buf, offset)

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
			uploads = append(uploads, lib.DeserializeUploadEntry(entry))
		}
	}

	return uploads, nil
}
