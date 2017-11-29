package models

import (
	"log"
)

func migrate() {
	initializeModel(&Job{})
	initializeModel(&JobRun{})
}

func initializeModel(klass interface{}) {
	err := getDB().Init(klass)
	if err != nil {
		log.Fatal(err)
	}
}
