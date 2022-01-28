package storage

import (
	"github.com/jameycribbs/hare"
	"github.com/jameycribbs/hare/datastores/disk"
)

var Database *hare.Database

func init() {
	ds, err := disk.New("./", ".json")
	if err != nil {
		panic(err)
	}

	db, err := hare.New(ds)
	if err != nil {
		panic(err)
	}

	Database = db
}
