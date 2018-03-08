package models

import (
	"log"
)

func (orm ORM) migrate() {
	orm.initializeModel(&JobSpec{})
	orm.initializeModel(&JobRun{})
	orm.initializeModel(&Initiator{})
	orm.initializeModel(&Tx{})
	orm.initializeModel(&TxAttempt{})
	orm.initializeModel(&BridgeType{})
	orm.initializeModel(&IndexableBlockNumber{})
}

func (orm ORM) initializeModel(klass interface{}) {
	err := orm.InitBucket(klass)
	if err != nil {
		log.Fatal(err)
	}
}
