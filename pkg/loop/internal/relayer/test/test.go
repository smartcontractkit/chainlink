package test

import (
	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

const ConfigTOML = `[Foo]
Bar = "Baz"
`

var (
	nodes = []types.NodeStatus{{
		ChainID: "foo",
		State:   "Alive",
		Config: `Name = 'bar'
URL = 'http://example.com'
`}, {
		ChainID: "foo",
		State:   "Alive",
		Config: `Name = 'baz'
URL = 'https://test.url'
`}}
	PluginArgs = types.PluginArgs{
		TransmitterID: "testtransmitter",
		PluginConfig:  []byte{100: 88},
	}
	RelayArgs = types.RelayArgs{
		ExternalJobID: uuid.MustParse("1051429b-aa66-11ed-b0d2-5cff35dfbe67"),
		JobID:         123,
		ContractID:    "testcontract",
		New:           true,
		RelayConfig:   []byte{42: 11},
		ProviderType:  string(types.Median),
	}
)
