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
	db = initializeDatabase("production")
	migrate()
}

func InitTest() {
	os.Remove(dbpath("test"))
	db = initializeDatabase("test")
	migrate()
}

func GetDB() *storm.DB {
	return db
}

func Close() {
	db.Close()
}

func initializeDatabase(env string) *storm.DB {
	db, err := storm.Open(dbpath(env))
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func dbpath(env string) string {
	dir, err := homedir.Expand("~/.chainlink")
	if err != nil {
		log.Fatal(err)
	}

	os.MkdirAll(dir, os.FileMode(0700))
	return path.Join(dir, "db."+env+".bolt")
}
