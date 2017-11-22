package orm

import (
	"github.com/asdine/storm"
	homedir "github.com/mitchellh/go-homedir"
	"log"
	"os"
	"path"
)

type chainlinkDB interface {
	GetDB() *storm.DB
	Close()
}

type ephemeralDB struct {
	*storm.DB
}

type persistentDB struct {
	*storm.DB
}

var db chainlinkDB

func Init() {
	db = persistentDB{initializeDatabase("production")}
	migrate()
}

func InitTest() {
	os.Remove(dbpath("test"))
	db = ephemeralDB{initializeDatabase("test")}
	migrate()
}

func GetDB() *storm.DB {
	return db.GetDB()
}

func Close() {
	db.Close()
}

func (d ephemeralDB) GetDB() *storm.DB {
	return d.DB
}

func (d persistentDB) GetDB() *storm.DB {
	return d.DB
}

func (d ephemeralDB) Close() {
	d.GetDB().Close()
}

func (d persistentDB) Close() {
	d.GetDB().Close()
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
