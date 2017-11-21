package main

import (
	"github.com/smartcontractkit/chainlink-go/db"
	"github.com/smartcontractkit/chainlink-go/web"
	"log"
)

func main() {
	db.Init()
	defer db.Close()
	r := web.Router()
	log.Fatal(r.Run())
}
