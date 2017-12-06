package main

import (
	"github.com/smartcontractkit/chainlink-go/logger"
	"github.com/smartcontractkit/chainlink-go/store"
	"github.com/smartcontractkit/chainlink-go/web"
)

func main() {
	defer logger.Sync()
	store := store.New()
	r := web.Router(store)
	err := store.Start()
	if err != nil {
		logger.Fatal(err)
	}
	defer store.Close()
	logger.Fatal(r.Run())
}
