package keeper

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

type registry struct {
	pgo  uint32
	mpds uint32
}

func (r *registry) CheckGasOverhead() uint32   { return uint32(0) }
func (r *registry) PerformGasOverhead() uint32 { return r.pgo }
func (r *registry) MaxPerformDataSize() uint32 { return r.mpds }

func TestBuildJobSpec(t *testing.T) {
	from := types.EIP55Address(testutils.NewAddress().Hex())
	contract := types.EIP55Address(testutils.NewAddress().Hex())
	chainID := "250"
	jb := job.Job{
		ID: 10,
		KeeperSpec: &job.KeeperSpec{
			FromAddress:     from,
			ContractAddress: contract,
		}}

	upkeepID := big.NewI(4)
	upkeep := UpkeepRegistration{
		Registry: Registry{
			FromAddress:     from,
			ContractAddress: contract,
			CheckGas:        11,
		},
		UpkeepID:   upkeepID,
		ExecuteGas: 12,
	}
	gasPrice := assets.NewWeiI(24)
	gasTipCap := assets.NewWeiI(48)
	gasFeeCap := assets.NewWeiI(72)

	r := &registry{
		pgo:  uint32(9),
		mpds: uint32(1000),
	}

	spec := buildJobSpec(jb, jb.KeeperSpec.FromAddress.Address(), upkeep, r, gasPrice, gasTipCap, gasFeeCap, chainID)

	expected := map[string]interface{}{
		"jobSpec": map[string]interface{}{
			"jobID":                  int32(10),
			"fromAddress":            from.String(),
			"effectiveKeeperAddress": jb.KeeperSpec.FromAddress.String(),
			"contractAddress":        contract.String(),
			"upkeepID":               "4",
			"prettyID":               fmt.Sprintf("UPx%064d", 4),
			"pipelineSpec": &pipeline.Spec{
				ForwardingAllowed: false,
			},
			"performUpkeepGasLimit": uint32(5_000_000 + 9),
			"maxPerformDataSize":    uint32(1000),
			"gasPrice":              gasPrice.ToInt(),
			"gasTipCap":             gasTipCap.ToInt(),
			"gasFeeCap":             gasFeeCap.ToInt(),
			"evmChainID":            "250",
		},
	}

	require.Equal(t, expected, spec)
}
