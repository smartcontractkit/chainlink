package stellar

import (
	"context"
	"fmt"
	"math/big"

	"github.com/stellar/go-stellar-sdk/xdr"

	cache "github.com/smartcontractkit/chainlink-stellar/bindings/contracts/data_feeds_cache"
	"github.com/smartcontractkit/chainlink-stellar/bindings/scval"
	stellardeployment "github.com/smartcontractkit/chainlink-stellar/deployment"
	datafeeds "github.com/smartcontractkit/chainlink-stellar/deployment/data-feeds"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	stellchain "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/stellar"
)

func StellarForwarderAddress(creEnv *cre.Environment, chainSelector uint64) string {
	return mustForwarderAddress(creEnv.CldfEnvironment.DataStore, chainSelector)
}

type DFReportEntry struct {
	DataID    [32]byte
	Answer    *big.Int
	Timestamp uint64
}

func DeployAndConfigureStellarDFCache(
	ctx context.Context,
	chain *stellchain.Blockchain,
	dataID [16]byte,
	description string,
	forwarderAddress string,
	workflowOwner [20]byte,
	workflowName [10]byte,
) (string, error) {
	stellarChain, err := stellarCldfChain(chain)
	if err != nil {
		return "", err
	}
	owner := stellarChain.Signer.Address()
	if fundErr := chain.Fund(ctx, owner, 0); fundErr != nil {
		return "", fmt.Errorf("failed to fund stellar deployer %s via friendbot: %w", owner, fundErr)
	}
	deployer, err := stellardeployment.NewDeployerFromChain(stellarChain)
	if err != nil {
		return "", fmt.Errorf("failed to build stellar deployer: %w", err)
	}

	wasm, err := datafeeds.Artifact(datafeeds.DataFeedsCacheWasm)
	if err != nil {
		return "", fmt.Errorf("failed to source data feeds cache WASM: %w", err)
	}

	var salt [32]byte
	cacheID, err := deployer.DeployContractBytesWithArgs(ctx, wasm, salt, []xdr.ScVal{scval.AddressToScVal(owner)})
	if err != nil {
		return "", fmt.Errorf("failed to deploy data feeds cache: %w", err)
	}

	c := cache.NewDataFeedsCacheClient(deployer, cacheID)
	if err := c.AddFeedAdmin(ctx, owner); err != nil {
		return "", fmt.Errorf("failed to register feed admin on %s: %w", cacheID, err)
	}
	entries := []cache.FeedConfigEntry{{
		DataId: dataID,
		Config: cache.FeedConfig{
			Description: description,
			WorkflowPermissions: []cache.WorkflowPermission{{
				AllowedSender:        forwarderAddress,
				AllowedWorkflowOwner: workflowOwner,
				AllowedWorkflowName:  workflowName,
			}},
		},
	}}
	if err := c.SetFeedConfigs(ctx, owner, entries); err != nil {
		return "", fmt.Errorf("failed to configure DF cache feed on %s: %w", cacheID, err)
	}
	return cacheID, nil
}

func DFCacheLatestRound(ctx context.Context, chain *stellchain.Blockchain, contractID string, dataID [16]byte) (*cache.RoundData, error) {
	stellarChain, err := stellarCldfChain(chain)
	if err != nil {
		return nil, err
	}
	deployer, err := stellardeployment.NewDeployerFromChain(stellarChain)
	if err != nil {
		return nil, fmt.Errorf("failed to build stellar deployer: %w", err)
	}
	return cache.NewDataFeedsCacheClient(deployer, contractID).LatestRound(ctx, dataID)
}

func BuildDFReportPayload(entries []DFReportEntry) ([]byte, error) {
	vec := make(xdr.ScVec, 0, len(entries))
	for _, e := range entries {
		answer, err := scval.I256ToScVal(e.Answer)
		if err != nil {
			return nil, fmt.Errorf("encode answer for feed %x: %w", e.DataID, err)
		}
		dataID := e.DataID
		m := &xdr.ScMap{
			structEntry("answer", answer),
			structEntry("data_id", xdr.ScVal{Type: xdr.ScValTypeScvBytes, Bytes: bytesPtr(dataID[:])}),
			structEntry("timestamp", xdr.ScVal{Type: xdr.ScValTypeScvU64, U64: u64Ptr(xdr.Uint64(e.Timestamp))}),
		}
		vec = append(vec, xdr.ScVal{Type: xdr.ScValTypeScvMap, Map: &m})
	}
	vecPtr := &vec
	val := xdr.ScVal{Type: xdr.ScValTypeScvVec, Vec: &vecPtr}
	return val.MarshalBinary()
}

func structEntry(key string, val xdr.ScVal) xdr.ScMapEntry {
	sym := xdr.ScSymbol(key)
	return xdr.ScMapEntry{
		Key: xdr.ScVal{Type: xdr.ScValTypeScvSymbol, Sym: &sym},
		Val: val,
	}
}

func bytesPtr(b []byte) *xdr.ScBytes {
	sb := xdr.ScBytes(b)
	return &sb
}

func u64Ptr(v xdr.Uint64) *xdr.Uint64 {
	return &v
}
