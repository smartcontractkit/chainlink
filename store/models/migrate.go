package models

import (
	"log"
)

func (orm ORM) migrate() {
	orm.initializeModel(&Job{})
	orm.initializeModel(&JobRun{})
	orm.initializeModel(&Initiator{})
	orm.initializeModel(&Tx{})
	orm.initializeModel(&TxAttempt{})
}

func (orm ORM) initializeModel(klass interface{}) {
	err := orm.InitBucket(klass)
	if err != nil {
		log.Fatal(err)
	}
}
