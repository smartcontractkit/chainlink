package main

import (
	"context"
	"log"
	"os"

	"github.com/gobuffalo/packr/builder"
)

// main builds packr boxes to serve static assets.
func main() {
	b := builder.New(context.Background(), os.Args[1])
	b.Compress = true

	err := b.Run()
	if err != nil {
		log.Fatal(err)
	}
}
