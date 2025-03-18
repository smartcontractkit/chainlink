package v0_5_0

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
)

func TestCallPayRecipients(t *testing.T) {
	e := testutil.NewMemoryEnv(t, true)

	e, rewardManagerAddr, _ := DeployRewardManagerAndLink(t, e)

	var poolId [32]byte
	copy(poolId[:], []byte("poolId"))

	_, err := commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			PayRecipientsChangeset,
			PayRecipientsConfig{
				ConfigsByChain: map[uint64][]PayRecipients{
					testutil.TestChain.Selector: {PayRecipients{
						RewardManagerAddress: rewardManagerAddr,
						PoolId:               poolId,
						Recipients:           []common.Address{},
					}},
				},
			},
		),
	)
	// Need Configured Fee Manager For PayRecipients Event
	require.NoError(t, err)
}
