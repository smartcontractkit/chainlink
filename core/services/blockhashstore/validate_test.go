package blockhashstore

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestValidate(t *testing.T) {
	v1Coordinator := ethkey.EIP55Address("0x1F72B4A5DCf7CC6d2E38423bF2f4BFA7db97d139")
	v2Coordinator := ethkey.EIP55Address("0x2be990eE17832b59E0086534c5ea2459Aa75E38F")
	fromAddresses := []ethkey.EIP55Address{("0x469aA2CD13e037DC5236320783dCfd0e641c0559")}

	var tests = []struct {
		name      string
		toml      string
		assertion func(t *testing.T, os job.Job, err error)
	}{
		{
			name: "valid",
			toml: `
type = "blockhashstore"
name = "valid-test"
coordinatorV1Address = "0x1F72B4A5DCf7CC6d2E38423bF2f4BFA7db97d139"
coordinatorV2Address = "0x2be990eE17832b59E0086534c5ea2459Aa75E38F"
waitBlocks = 59
lookbackBlocks = 159
blockhashStoreAddress = "0x3e20Cef636EdA7ba135bCbA4fe6177Bd3cE0aB17"
pollPeriod = "23s"
runTimeout = "7s"
evmChainID = "4"
fromAddresses = ["0x469aA2CD13e037DC5236320783dCfd0e641c0559"]`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.NoError(t, err)
				require.Equal(t, job.BlockhashStore, os.Type)
				require.Equal(t, "valid-test", os.Name.String)
				require.Equal(t, &v1Coordinator,
					os.BlockhashStoreSpec.CoordinatorV1Address)
				require.Equal(t, &v2Coordinator,
					os.BlockhashStoreSpec.CoordinatorV2Address)
				require.Equal(t, int32(59), os.BlockhashStoreSpec.WaitBlocks)
				require.Equal(t, int32(159), os.BlockhashStoreSpec.LookbackBlocks)
				require.Equal(t, ethkey.EIP55Address("0x3e20Cef636EdA7ba135bCbA4fe6177Bd3cE0aB17"),
					os.BlockhashStoreSpec.BlockhashStoreAddress)
				require.Equal(t, 23*time.Second, os.BlockhashStoreSpec.PollPeriod)
				require.Equal(t, 7*time.Second, os.BlockhashStoreSpec.RunTimeout)
				require.Equal(t, utils.NewBigI(4), os.BlockhashStoreSpec.EVMChainID)
				require.Equal(t, fromAddresses,
					os.BlockhashStoreSpec.FromAddresses)
			},
		},
		{
			name: "defaults",
			toml: `
type = "blockhashstore"
name = "defaults-test"
coordinatorV1Address = "0x1F72B4A5DCf7CC6d2E38423bF2f4BFA7db97d139"
coordinatorV2Address = "0x2be990eE17832b59E0086534c5ea2459Aa75E38F"
blockhashStoreAddress = "0x3e20Cef636EdA7ba135bCbA4fe6177Bd3cE0aB17"
evmChainID = "4"`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.NoError(t, err)
				require.Equal(t, int32(100), os.BlockhashStoreSpec.WaitBlocks)
				require.Equal(t, int32(200), os.BlockhashStoreSpec.LookbackBlocks)
				require.Nil(t, os.BlockhashStoreSpec.FromAddresses)
				require.Equal(t, 30*time.Second, os.BlockhashStoreSpec.PollPeriod)
				require.Equal(t, 30*time.Second, os.BlockhashStoreSpec.RunTimeout)
			},
		},
		{
			name: "v1 only",
			toml: `
type = "blockhashstore"
name = "defaults-test"
coordinatorV1Address = "0x1F72B4A5DCf7CC6d2E38423bF2f4BFA7db97d139"
blockhashStoreAddress = "0x3e20Cef636EdA7ba135bCbA4fe6177Bd3cE0aB17"
evmChainID = "4"
fromAddresses = ["0x469aA2CD13e037DC5236320783dCfd0e641c0559"]`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.NoError(t, err)
				require.Equal(t, &v1Coordinator,
					os.BlockhashStoreSpec.CoordinatorV1Address)
				require.Nil(t, os.BlockhashStoreSpec.CoordinatorV2Address)
			},
		},
		{
			name: "v2 only",
			toml: `
type = "blockhashstore"
name = "defaults-test"
coordinatorV2Address = "0x2be990eE17832b59E0086534c5ea2459Aa75E38F"
blockhashStoreAddress = "0x3e20Cef636EdA7ba135bCbA4fe6177Bd3cE0aB17"
evmChainID = "4"
fromAddresses = ["0x469aA2CD13e037DC5236320783dCfd0e641c0559"]`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.NoError(t, err)
				require.Nil(t, os.BlockhashStoreSpec.CoordinatorV1Address)
				require.Equal(t, &v2Coordinator, os.BlockhashStoreSpec.CoordinatorV2Address)
			},
		},
		{
			name: "invalid no coordinators",
			toml: `
type = "blockhashstore"
name = "defaults-test"
blockhashStoreAddress = "0x3e20Cef636EdA7ba135bCbA4fe6177Bd3cE0aB17"
evmChainID = "4"
fromAddresses = ["0x469aA2CD13e037DC5236320783dCfd0e641c0559"]`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.EqualError(t, err, `at least one of "coordinatorV1Address" and "coordinatorV2Address" must be set`)
			},
		},
		{
			name: "invalid no blockhashstore",
			toml: `
type = "blockhashstore"
name = "defaults-test"
coordinatorV2Address = "0x2be990eE17832b59E0086534c5ea2459Aa75E38F"
evmChainID = "4"
fromAddresses = ["0x469aA2CD13e037DC5236320783dCfd0e641c0559"]`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.EqualError(t, err, `"blockhashStoreAddress" must be set`)
			},
		},
		{
			name: "invalid no chain ID",
			toml: `
type = "blockhashstore"
name = "defaults-test"
coordinatorV2Address = "0x2be990eE17832b59E0086534c5ea2459Aa75E38F"
blockhashStoreAddress = "0x3e20Cef636EdA7ba135bCbA4fe6177Bd3cE0aB17"
fromAddresses = ["0x469aA2CD13e037DC5236320783dCfd0e641c0559"]`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.EqualError(t, err, `"evmChainID" must be set`)
			},
		},
		{
			name: "invalid waitBlocks too high",
			toml: `
type = "blockhashstore"
name = "valid-test"
coordinatorV1Address = "0x1F72B4A5DCf7CC6d2E38423bF2f4BFA7db97d139"
coordinatorV2Address = "0x2be990eE17832b59E0086534c5ea2459Aa75E38F"
waitBlocks = 257
lookbackBlocks = 258
blockhashStoreAddress = "0x3e20Cef636EdA7ba135bCbA4fe6177Bd3cE0aB17"
evmChainID = "4"`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.EqualError(t, err, `"waitBlocks" must be less than 256`)
			},
		},
		{
			name: "invalid lookbackBlocks too high",
			toml: `
type = "blockhashstore"
name = "valid-test"
coordinatorV1Address = "0x1F72B4A5DCf7CC6d2E38423bF2f4BFA7db97d139"
coordinatorV2Address = "0x2be990eE17832b59E0086534c5ea2459Aa75E38F"
lookbackBlocks = 257
blockhashStoreAddress = "0x3e20Cef636EdA7ba135bCbA4fe6177Bd3cE0aB17"
evmChainID = "4"`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.EqualError(t, err, `"lookbackBlocks" must be less than 256`)
			},
		},
		{
			name: "invalid waitBlocks higher than lookbackBlocks",
			toml: `
type = "blockhashstore"
name = "valid-test"
coordinatorV1Address = "0x1F72B4A5DCf7CC6d2E38423bF2f4BFA7db97d139"
coordinatorV2Address = "0x2be990eE17832b59E0086534c5ea2459Aa75E38F"
waitBlocks = 200
lookbackBlocks = 100
blockhashStoreAddress = "0x3e20Cef636EdA7ba135bCbA4fe6177Bd3cE0aB17"
evmChainID = "4"`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.EqualError(t, err, `"waitBlocks" must be less than "lookbackBlocks"`)
			},
		},
		{
			name: "invalid toml",
			toml: `
type = invalid`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "loading toml")
			},
		},
		{
			name: "toml wrong type for spec",
			toml: `
type = 123`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "unmarshalling toml spec")
			},
		},
		{
			name: "toml wrong type for job",
			toml: `
type = "blockhashstore"
waitBlocks = "shouldBeInt"`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "unmarshalling toml job")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s, err := ValidatedSpec(test.toml)
			test.assertion(t, s, err)
		})
	}
}
