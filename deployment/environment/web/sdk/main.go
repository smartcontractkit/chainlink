package main

import (
	"fmt"
	"os"

	"github.com/Khan/genqlient/generate"

	"github.com/smartcontractkit/chainlink/v2/core/web/schema"
)

func main() {
	schema := schema.MustGetRootSchema()

	if err := os.WriteFile("./internal/schema.graphql", []byte(schema), 0600); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	generate.Main()
}
