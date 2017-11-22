package orm

import (
	"github.com/smartcontractkit/chainlink-go/models"
	"log"
)

func migrate() {
	err := db.Init(&models.Job{})
	if err != nil {
		log.Fatal(err)
	}
}
