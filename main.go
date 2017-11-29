package main

import (
	"github.com/smartcontractkit/chainlink-go/models"
	"github.com/smartcontractkit/chainlink-go/web"
	"log"
)

func main() {
	models.InitDB()
	defer models.CloseDB()
	r := web.Router()
	log.Fatal(r.Run())
}
