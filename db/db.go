package db

import (
	"github.com/jinzhu/gorm"
	homedir "github.com/mitchellh/go-homedir"
	"log"
	"os"
	"path"
)

var db *gorm.DB

func Init() {
	dir, err := homedir.Expand("~/.chainlink")
	if err != nil {
		log.Fatal(err)
	}

	os.MkdirAll(dir, os.FileMode(0700))
	db, err = gorm.Open("sqlite3", path.Join(dir, "db.sqlite3"))
	if err != nil {
		log.Fatal(err)
	}
	if err = db.DB().Ping(); err != nil {
		log.Fatal(err)
	}

	migrate()
}

func GetDB() *gorm.DB {
	return db
}

func Close() {
	GetDB().Close()
}
