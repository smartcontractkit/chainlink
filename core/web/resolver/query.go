package resolver

import (
	"context"
	"database/sql"
	"fmt"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/chains"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/vrfkey"
	evmrelay "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/utils/stringutils"
	"github.com/smartcontractkit/chainlink/v2/core/web/loader"
)

// Bridge retrieves a bridges by name.
func (r *Resolver) Bridge(ctx context.Context, args struct{ ID graphql.ID }) (*BridgePayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	name, err := bridges.ParseBridgeName(string(args.ID))
	if err != nil {
		return nil, err
	}

	bridge, err := r.App.BridgeORM().FindBridge(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewBridgePayload(bridge, err), nil
		}

		return nil, err
	}

	return NewBridgePayload(bridge, nil), nil
}

// Bridges retrieves a paginated list of bridges.
func (r *Resolver) Bridges(ctx context.Context, args struct {
	Offset *int32
	Limit  *int32
}) (*BridgesPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	offset := pageOffset(args.Offset)
	limit := pageLimit(args.Limit)

	brdgs, count, err := r.App.BridgeORM().BridgeTypes(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	return NewBridgesPayload(brdgs, int32(count)), nil
}

// Chain retrieves a chain by id.
func (r *Resolver) Chain(ctx context.Context,
	args struct {
		ID      graphql.ID
		Network *string
	}) (*ChainPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	// fall back to original behaviour if network is not provided
	if args.Network == nil {
		id, err := loader.GetChainByID(ctx, string(args.ID))
		if err != nil {
			if errors.Is(err, chains.ErrNotFound) {
				return NewChainPayload(chainlink.NetworkChainStatus{}, chains.ErrNotFound), nil
			}
			return nil, err
		}
		return NewChainPayload(*id, nil), nil
	}

	relayID := types.NewRelayID(*args.Network, string(args.ID))
	id, err := loader.GetChainByRelayID(ctx, relayID.Name())
	if err != nil {
		if errors.Is(err, chains.ErrNotFound) {
			return NewChainPayload(chainlink.NetworkChainStatus{}, chains.ErrNotFound), nil
		}
		return nil, err
	}

	return NewChainPayload(*id, nil), nil
}

// Chains retrieves a paginated list of chains.
func (r *Resolver) Chains(ctx context.Context, args struct {
	Offset *int32
	Limit  *int32
}) (*ChainsPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	offset := pageOffset(args.Offset)
	limit := pageLimit(args.Limit)

	relayersMap := r.App.GetRelayers().GetIDToRelayerMap()

	chains := make([]chainlink.NetworkChainStatus, 0, len(relayersMap))
	for k, v := range relayersMap {
		s, err := v.GetChainStatus(ctx)
		if err != nil {
			return nil, err
		}

		chains = append(chains, chainlink.NetworkChainStatus{
			ChainStatus: s,
			Network:     k.Network,
		})
	}

	count := len(chains)
	if count == 0 {
		// No chains are configured, return an empty ChainsPayload, so we don't break the UI
		return NewChainsPayload(nil, 0), nil
	}

	// bound the chain results
	if offset >= count {
		return nil, fmt.Errorf("offset %d out of range", offset)
	}
	end := count
	if limit > 0 && offset+limit < end {
		end = offset + limit
	}

	sortByNetworkAndID(chains)
	return NewChainsPayload(chains[offset:end], int32(count)), nil
}

func sortByNetworkAndID(chains []chainlink.NetworkChainStatus) {
	sort.SliceStable(chains, func(i, j int) bool {
		if chains[i].Network == chains[j].Network {
			return chains[i].ID < chains[j].ID
		}
		return chains[i].Network < chains[j].Network
	})
}

// FeedsManager retrieves a feeds manager by id.
func (r *Resolver) FeedsManager(ctx context.Context, args struct{ ID graphql.ID }) (*FeedsManagerPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	id, err := stringutils.ToInt64(string(args.ID))
	if err != nil {
		return nil, err
	}

	mgr, err := r.App.GetFeedsService().GetManager(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewFeedsManagerPayload(nil, err), nil
		}

		return nil, err
	}

	return NewFeedsManagerPayload(mgr, nil), nil
}

func (r *Resolver) FeedsManagers(ctx context.Context) (*FeedsManagersPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	mgrs, err := r.App.GetFeedsService().ListManagers(ctx)
	if err != nil {
		return nil, err
	}

	return NewFeedsManagersPayload(mgrs), nil
}

// Job retrieves a job by id.
func (r *Resolver) Job(ctx context.Context, args struct{ ID graphql.ID }) (*JobPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	id, err := stringutils.ToInt32(string(args.ID))
	if err != nil {
		return nil, err
	}

	j, err := r.App.JobORM().FindJobWithoutSpecErrors(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewJobPayload(r.App, nil, err), nil
		}

		// We still need to show the job in UI/CLI even if the chain id is disabled
		if errors.Is(err, chains.ErrNoSuchChainID) {
			return NewJobPayload(r.App, &j, err), nil
		}

		return nil, err
	}

	return NewJobPayload(r.App, &j, nil), nil
}

// Jobs fetches a paginated list of jobs
func (r *Resolver) Jobs(ctx context.Context, args struct {
	Offset *int32
	Limit  *int32
}) (*JobsPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	offset := pageOffset(args.Offset)
	limit := pageLimit(args.Limit)

	jobs, count, err := r.App.JobORM().FindJobs(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	return NewJobsPayload(r.App, jobs, int32(count)), nil
}

func (r *Resolver) OCRKeyBundles(ctx context.Context) (*OCRKeyBundlesPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	ocrKeyBundles, err := r.App.GetKeyStore().OCR().GetAll()
	if err != nil {
		return nil, err
	}

	return NewOCRKeyBundlesPayloadResolver(ocrKeyBundles), nil
}

func (r *Resolver) CSAKeys(ctx context.Context) (*CSAKeysPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	keys, err := r.App.GetKeyStore().CSA().GetAll()
	if err != nil {
		return nil, err
	}

	return NewCSAKeysResolver(keys), nil
}

// Features retrieves each featured enabled by boolean mapping
func (r *Resolver) Features(ctx context.Context) (*FeaturesPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	return NewFeaturesPayloadResolver(r.App.GetConfig().Feature()), nil
}

// Node retrieves a node by ID (Name)
func (r *Resolver) Node(ctx context.Context, args struct{ ID graphql.ID }) (*NodePayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}
	r.App.GetLogger().Debug("resolver Node args %v", args)
	name := string(args.ID)
	r.App.GetLogger().Debug("resolver Node name %s", name)

	for _, relayer := range r.App.GetRelayers().Slice() {
		statuses, _, _, err := relayer.ListNodeStatuses(ctx, 0, "")
		if err != nil {
			return nil, err
		}
		for i, s := range statuses {
			if s.Name == name {
				npr, err2 := NewNodePayloadResolver(&statuses[i], nil)
				if err2 != nil {
					return nil, err2
				}
				return npr, nil
			}
		}
	}

	r.App.GetLogger().Errorw("resolver getting node status", "err", chains.ErrNotFound)
	return NewNodePayloadResolver(nil, chains.ErrNotFound)
}

func (r *Resolver) P2PKeys(ctx context.Context) (*P2PKeysPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	p2pKeys, err := r.App.GetKeyStore().P2P().GetAll()
	if err != nil {
		return nil, err
	}

	return NewP2PKeysPayload(p2pKeys), nil
}

// VRFKeys fetches all VRF keys.
func (r *Resolver) VRFKeys(ctx context.Context) (*VRFKeysPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	keys, err := r.App.GetKeyStore().VRF().GetAll()
	if err != nil {
		return nil, err
	}

	return NewVRFKeysPayloadResolver(keys), nil
}

// VRFKey fetches the VRF key with the given ID.
func (r *Resolver) VRFKey(ctx context.Context, args struct {
	ID graphql.ID
}) (*VRFKeyPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	key, err := r.App.GetKeyStore().VRF().Get(string(args.ID))
	if err != nil {
		if errors.Is(errors.Cause(err), keystore.ErrMissingVRFKey) {
			return NewVRFKeyPayloadResolver(vrfkey.KeyV2{}, err), nil
		}
		return nil, err
	}

	return NewVRFKeyPayloadResolver(key, nil), err
}

// JobProposal retrieves a job proposal by ID
func (r *Resolver) JobProposal(ctx context.Context, args struct {
	ID graphql.ID
}) (*JobProposalPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	id, err := stringutils.ToInt64(string(args.ID))
	if err != nil {
		return nil, err
	}

	jp, err := r.App.GetFeedsService().GetJobProposal(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewJobProposalPayload(nil, err), nil
		}

		return nil, err
	}

	return NewJobProposalPayload(jp, err), nil
}

// Nodes retrieves a paginated list of nodes.
func (r *Resolver) Nodes(ctx context.Context, args struct {
	Offset *int32
	Limit  *int32
}) (*NodesPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	offset := pageOffset(args.Offset)
	limit := pageLimit(args.Limit)
	r.App.GetLogger().Debugw("resolver Nodes query", "offset", offset, "limit", limit)
	allNodes, total, err := r.App.GetRelayers().NodeStatuses(ctx, offset, limit)
	r.App.GetLogger().Debugw("resolver Nodes query result", "nodes", allNodes, "total", total, "err", err)

	if err != nil {
		r.App.GetLogger().Errorw("Error creating get nodes status from app", "err", err)
		return nil, err
	}
	npr, warn := NewNodesPayload(allNodes, int32(total))
	if warn != nil {
		r.App.GetLogger().Warnw("Error creating NodesPayloadResolver", "err", warn)
	}
	return npr, nil
}

func (r *Resolver) JobRuns(ctx context.Context, args struct {
	Offset *int32
	Limit  *int32
}) (*JobRunsPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	limit := pageLimit(args.Limit)
	offset := pageOffset(args.Offset)

	runs, count, err := r.App.JobORM().PipelineRuns(ctx, nil, offset, limit)
	if err != nil {
		return nil, err
	}

	return NewJobRunsPayload(runs, int32(count), r.App), nil
}

func (r *Resolver) JobRun(ctx context.Context, args struct {
	ID graphql.ID
}) (*JobRunPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	id, err := stringutils.ToInt64(string(args.ID))
	if err != nil {
		return nil, err
	}

	jr, err := r.App.JobORM().FindPipelineRunByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewJobRunPayload(nil, r.App, err), nil
		}

		return nil, err
	}

	return NewJobRunPayload(&jr, r.App, err), nil
}

func (r *Resolver) ETHKeys(ctx context.Context) (*ETHKeysPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	ks := r.App.GetKeyStore().Eth()

	keys, err := ks.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting unlocked keys: %w", err)
	}

	states, err := ks.GetStatesForKeys(ctx, keys)
	if err != nil {
		return nil, fmt.Errorf("error getting key states: %w", err)
	}

	var ethKeys []ETHKey

	for _, state := range states {
		k, err := ks.Get(ctx, state.Address.Hex())
		if err != nil {
			return nil, err
		}

		chain, err := r.App.GetRelayers().LegacyEVMChains().Get(state.EVMChainID.String())
		if errors.Is(errors.Cause(err), evmrelay.ErrNoChains) {
			ethKeys = append(ethKeys, ETHKey{
				addr:  k.EIP55Address,
				state: state,
			})

			continue
		}
		// Don't include keys without valid chain.
		// OperatorUI fails to show keys where chains are not in the config.
		if err == nil {
			ethKeys = append(ethKeys, ETHKey{
				addr:  k.EIP55Address,
				state: state,
				chain: chain,
			})
		}
	}
	// Put disabled keys to the end
	sort.SliceStable(ethKeys, func(i, j int) bool {
		return !states[i].Disabled && states[j].Disabled
	})

	return NewETHKeysPayload(ethKeys), nil
}

// ConfigV2 retrieves the Chainlink node's configuration (V2 mode)
func (r *Resolver) ConfigV2(ctx context.Context) (*ConfigV2PayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	cfg := r.App.GetConfig()
	return NewConfigV2Payload(cfg.ConfigTOML()), nil
}

func (r *Resolver) EthTransaction(ctx context.Context, args struct {
	Hash graphql.ID
}) (*EthTransactionPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	hash := common.HexToHash(string(args.Hash))
	etx, err := r.App.TxmStorageService().FindTxByHash(ctx, hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewEthTransactionPayload(nil, err), nil
		}

		return nil, err
	}

	return NewEthTransactionPayload(etx, err), nil
}

func (r *Resolver) EthTransactions(ctx context.Context, args struct {
	Offset *int32
	Limit  *int32
}) (*EthTransactionsPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	offset := pageOffset(args.Offset)
	limit := pageLimit(args.Limit)

	txs, count, err := r.App.TxmStorageService().Transactions(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	return NewEthTransactionsPayload(txs, int32(count)), nil
}

func (r *Resolver) EthTransactionsAttempts(ctx context.Context, args struct {
	Offset *int32
	Limit  *int32
}) (*EthTransactionsAttemptsPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	offset := pageOffset(args.Offset)
	limit := pageLimit(args.Limit)

	attempts, count, err := r.App.TxmStorageService().TxAttempts(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	return NewEthTransactionsAttemptsPayload(attempts, int32(count)), nil
}

func (r *Resolver) GlobalLogLevel(ctx context.Context) (*GlobalLogLevelPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	logLevel := r.App.GetConfig().Log().Level().String()

	return NewGlobalLogLevelPayload(logLevel), nil
}

func (r *Resolver) SolanaKeys(ctx context.Context) (*SolanaKeysPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	keys, err := r.App.GetKeyStore().Solana().GetAll()
	if err != nil {
		return nil, err
	}

	return NewSolanaKeysPayload(keys), nil
}

func (r *Resolver) AptosKeys(ctx context.Context) (*AptosKeysPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	keys, err := r.App.GetKeyStore().Aptos().GetAll()
	if err != nil {
		return nil, err
	}

	return NewAptosKeysPayload(keys), nil
}

func (r *Resolver) CosmosKeys(ctx context.Context) (*CosmosKeysPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}
	keys, err := r.App.GetKeyStore().Cosmos().GetAll()
	if err != nil {
		return nil, err
	}

	return NewCosmosKeysPayload(keys), nil
}

func (r *Resolver) StarkNetKeys(ctx context.Context) (*StarkNetKeysPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}
	keys, err := r.App.GetKeyStore().StarkNet().GetAll()
	if err != nil {
		return nil, err
	}

	return NewStarkNetKeysPayload(keys), nil
}

func (r *Resolver) TronKeys(ctx context.Context) (*TronKeysPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	keys, err := r.App.GetKeyStore().Tron().GetAll()
	if err != nil {
		return nil, err
	}

	return NewTronKeysPayload(keys), nil
}

func (r *Resolver) TONKeys(ctx context.Context) (*TONKeysPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	keys, err := r.App.GetKeyStore().TON().GetAll()
	if err != nil {
		return nil, err
	}

	return NewTONKeysPayload(keys), nil
}

func (r *Resolver) SQLLogging(ctx context.Context) (*GetSQLLoggingPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	enabled := r.App.GetConfig().Database().LogSQL()

	return NewGetSQLLoggingPayload(enabled), nil
}

// OCR2KeyBundles resolves the list of OCR2 key bundles
func (r *Resolver) OCR2KeyBundles(ctx context.Context) (*OCR2KeyBundlesPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	ekbs, err := r.App.GetKeyStore().OCR2().GetAll()
	if err != nil {
		return nil, err
	}

	return NewOCR2KeyBundlesPayload(ekbs), nil
}
