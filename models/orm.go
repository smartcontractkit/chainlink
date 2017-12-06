package models

import (
	"fmt"
	"log"
	"os"
	"path"
	"reflect"
	"time"

	"github.com/asdine/storm"
	homedir "github.com/mitchellh/go-homedir"
)

type ORM struct {
	*storm.DB
}

func InitORM(env string) ORM {
	orm := ORM{initializeDatabase(env)}
	orm.migrate()
	return orm
}

func initializeDatabase(env string) *storm.DB {
	db, err := storm.Open(DBPath(env))
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func DBPath(env string) string {
	dir, err := homedir.Expand("~/.chainlink")
	if err != nil {
		log.Fatal(err)
	}

	if env == "test" {
		dir = path.Join(dir, "tmp")
		os.MkdirAll(dir, os.FileMode(0700))
		return path.Join(dir, "db."+fmt.Sprintf("%d", time.Now().UnixNano())+".bolt")
	}

	os.MkdirAll(dir, os.FileMode(0700))
	return path.Join(dir, "db."+env+".bolt")
}

func (self ORM) Where(field string, value interface{}, instance interface{}) error {
	err := self.Find(field, value, instance)
	if err == storm.ErrNotFound {
		emptySlice(instance)
		return nil
	}
	return err
}

func (self ORM) InitBucket(model interface{}) error {
	return self.Init(model)
}

func (self ORM) JobsWithCron() ([]Job, error) {
	jobs := []Job{}
	err := self.AllByIndex("Cron", &jobs)
	return jobs, err
}

func emptySlice(to interface{}) {
	ref := reflect.ValueOf(to)
	results := reflect.MakeSlice(reflect.Indirect(ref).Type(), 0, 0)
	reflect.Indirect(ref).Set(results)
}
