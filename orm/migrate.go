package orm

import (
	"github.com/smartcontractkit/chainlink-go/models"
	"log"
)

func migrate() {
	err := GetDB().Init(&models.Job{})
	if err != nil {
		log.Fatal(err)
	}
}
