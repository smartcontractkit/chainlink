package models

import (
	"log"
)

func (self ORM) migrate() {
	self.initializeModel(&Job{})
	self.initializeModel(&JobRun{})
	self.initializeModel(&Initiator{})
	self.initializeModel(&Tx{})
	self.initializeModel(&TxAttempt{})
}

func (self ORM) initializeModel(klass interface{}) {
	err := self.InitBucket(klass)
	if err != nil {
		log.Fatal(err)
	}
}
