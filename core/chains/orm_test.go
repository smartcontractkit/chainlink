package chains_test

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

type Config struct {
	Foo null.String
}

func (c *Config) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, c)
}

func (c *Config) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func ExampleNewORM() {
	type Node = struct {
		ID             int32
		Name           string
		ExampleChainID string
		URL            string
		Bar            null.Int
	}
	var q pg.Q
	_ = chains.NewORM[string, *Config, Node](q, "example", "url", "bar")

	// Output:
}
