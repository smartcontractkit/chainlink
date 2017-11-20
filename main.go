package main

import (
	"github.com/smartcontractkit/chainlink-go/web"
	"log"
)

func main() {
	r := web.Router()
	log.Fatal(r.Run())
}
