package models

import (
	"log"
	"path"
	"reflect"

	"github.com/asdine/storm"
)

type ORM struct {
	*storm.DB
}

func NewORM(dir string) *ORM {
	path := path.Join(dir, "db.bolt")
	orm := &ORM{initializeDatabase(path)}
	orm.migrate()
	return orm
}

func initializeDatabase(path string) *storm.DB {
	db, err := storm.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func (self *ORM) Where(field string, value interface{}, instance interface{}) error {
	err := self.Find(field, value, instance)
	if err == storm.ErrNotFound {
		emptySlice(instance)
		return nil
	}
	return err
}

func (self *ORM) InitBucket(model interface{}) error {
	return self.Init(model)
}

func (self *ORM) JobsWithCron() ([]Job, error) {
	jobs := []Job{}
	err := self.AllByIndex("Cron", &jobs)
	return jobs, err
}

func emptySlice(to interface{}) {
	ref := reflect.ValueOf(to)
	results := reflect.MakeSlice(reflect.Indirect(ref).Type(), 0, 0)
	reflect.Indirect(ref).Set(results)
}
