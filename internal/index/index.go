package index

import (
	"os"
)

/*
this indexing implementation only exists because i didn't want to deal with setting up cgo/mingw/gcc/go-sqlite3 on windows lol.
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
