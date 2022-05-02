package chains_test

import (
	"time"

	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

func ExampleNewORM() {
	type Config struct {
		Foo null.String
	}
	type Node = struct {
		ID             int32
		Name           string
		ExampleChainID string
		URL            string
		Bar            null.Int
		CreatedAt      time.Time
		UpdatedAt      time.Time
	}
	var q pg.Q
	_ = chains.NewORM[string, Config, Node](q, "example", "url", "bar")

	// Output:
}
