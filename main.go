package main

import (
	"github.com/smartcontractkit/chainlink-go/web"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/assignments", &web.Assignments{})
	server := http.Server{Handler: mux, Addr: "localhost:6688"}
	log.Fatal(server.ListenAndServe())
}
