package main

import (
	"github.com/smartcontractkit/chainlink-go/orm"
	"github.com/smartcontractkit/chainlink-go/web"
	"log"
)

func main() {
	orm.Init()
	defer orm.Close()
	r := web.Router()
	log.Fatal(r.Run())
}
