package mercury_common_test

import (
	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

var (
	RelayArgs = types.RelayArgs{
		ExternalJobID: uuid.MustParse("f1928153-d0b4-474b-9fd7-fd0985d0e7ca"),
		JobID:         42,
		ContractID:    "mercury-testcontract",
		New:           true,
		RelayConfig:   []byte{1: 4, 36: 101},
		ProviderType:  string(types.Mercury),
	}
	PluginArgs = types.PluginArgs{
		TransmitterID: "mercury-testtransmitter",
		PluginConfig:  []byte{133: 79},
	}
)
