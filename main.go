package main

import (
	"github.com/smartcontractkit/chainlink-go/store"
	"github.com/smartcontractkit/chainlink-go/web"
	"log"
)

func main() {
	store := store.New()
	r := web.Router(store)
	err := store.Start()
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()
	log.Fatal(r.Run())
}
