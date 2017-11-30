package main

import (
	"github.com/smartcontractkit/chainlink-go/models"
	"github.com/smartcontractkit/chainlink-go/scheduler"
	"github.com/smartcontractkit/chainlink-go/web"
	"log"
)

func main() {
	models.InitDB()
	defer models.CloseDB()
	sched, err := scheduler.Start()
	if err != nil {
		log.Fatal(err)
	}
	defer sched.Stop()
	r := web.Router()

	log.Fatal(r.Run())
}
