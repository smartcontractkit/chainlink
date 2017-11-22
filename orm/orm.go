package orm

import (
	"github.com/asdine/storm"
	homedir "github.com/mitchellh/go-homedir"
	"log"
	"os"
	"path"
)

var db *storm.DB

func Init() {
	dir, err := homedir.Expand("~/.chainlink")
	if err != nil {
		log.Fatal(err)
	}

	os.MkdirAll(dir, os.FileMode(0700))
	db, err = storm.Open(path.Join(dir, "db.bolt"))
	if err != nil {
		log.Fatal(err)
	}

	migrate()
}

func GetDB() *storm.DB {
	return db
}

func Close() {
	db.Close()
}
