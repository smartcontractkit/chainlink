package models

import (
	"fmt"
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
	orm.initializeModel(&User{})
}

func (orm ORM) initializeModel(klass interface{}) {
	err := orm.InitBucket(klass)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to migrate %T: %+v", klass, err))
	}
}
