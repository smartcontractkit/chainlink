package solana

import (
	"errors"
	"fmt"

	"github.com/gagliardetto/solana-go"

	bin "github.com/gagliardetto/binary"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"

	df_store "github.com/smartcontractkit/chainlink-solana/contracts/generated/data_feeds_store"
)

// StoreTarget is the data feeds store a changeset configures. Feed admin instructions are signed
// by the deployer, or by the timelock signer when MCMS is set, which then has to be a feed admin.
type StoreTarget struct {
	ChainSel  uint64
	Version   string
	Qualifier string
	MCMS      *cldfproposalutils.TimelockConfig
}

func (t StoreTarget) ref() storeRef {
	return storeRef{t.ChainSel, t.Version, t.Qualifier}
}

func (t StoreTarget) verify(env cldf.Environment) error {
	return verifyStoreRequest(env, t.ref(), t.MCMS)
}

// load returns the store and the account signing as feed admin.
func (t StoreTarget) load(env cldf.Environment) (storeAddresses, solana.PublicKey, error) {
	store, err := loadStore(env.DataStore, t.ref())
	if err != nil {
		return storeAddresses{}, solana.PublicKey{}, err
	}
	feedAdmin, err := authority(env, t.ChainSel, t.MCMS)
	if err != nil {
		return storeAddresses{}, solana.PublicKey{}, err
	}
	return store, feedAdmin, nil
}

// send sends the instructions to the target store, or proposes them under MCMS.
func (t StoreTarget) send(env cldf.Environment, description string, ixs []solana.Instruction) (cldf.ChangesetOutput, error) {
	return sendInstructions(env, t.ChainSel, t.MCMS, description, ixs)
}

// buildBatched builds one instruction per batch of items.
func buildBatched[T any](env cldf.Environment, t StoreTarget, items []T,
	build func(store storeAddresses, feedAdmin solana.PublicKey, batch []T) (solana.Instruction, error),
) ([]solana.Instruction, error) {
	store, feedAdmin, err := t.load(env)
	if err != nil {
		return nil, err
	}
	var ixs []solana.Instruction
	for _, batch := range chunk(items, storeBatchSize) {
		ix, err := build(store, feedAdmin, batch)
		if err != nil {
			return nil, err
		}
		ixs = append(ixs, ix)
	}
	return ixs, nil
}

// ----------------------------------------------------------------------------------------------
// Workflow metadata
// ----------------------------------------------------------------------------------------------

type SetStoreDefaultWorkflowMetadataRequest struct {
	StoreTarget
	// Workflows allowed to write any data ID without feed-specific metadata, replacing the current list
	Workflows []StoreWorkflow
}

var _ cldf.ChangeSetV2[*SetStoreDefaultWorkflowMetadataRequest] = SetStoreDefaultWorkflowMetadata{}

type SetStoreDefaultWorkflowMetadata struct{}

func (cs SetStoreDefaultWorkflowMetadata) VerifyPreconditions(env cldf.Environment, req *SetStoreDefaultWorkflowMetadataRequest) error {
	if _, err := workflowMetadata(solana.PublicKey{}, req.Workflows); err != nil {
		return err
	}
	return req.verify(env)
}

func (cs SetStoreDefaultWorkflowMetadata) Apply(env cldf.Environment, req *SetStoreDefaultWorkflowMetadataRequest) (cldf.ChangesetOutput, error) {
	store, feedAdmin, err := req.load(env)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	metadata, err := workflowMetadata(store.ProgramID, req.Workflows)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	ix, err := storeInstruction(store.ProgramID, nil)(df_store.NewSetDefaultFeedMetadataInstruction(metadata, feedAdmin, store.State))
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	return sendInstructions(env, req.ChainSel, req.MCMS, "set data feeds store default workflow metadata", []solana.Instruction{ix})
}

type SetStoreSpecificWorkflowMetadataRequest struct {
	StoreTarget
	DataIDs []string
	// Workflows allowed to write the data IDs. Empty removes the feed-specific metadata, so the
	// data IDs fall back to the default workflows.
	Workflows []StoreWorkflow
}

var _ cldf.ChangeSetV2[*SetStoreSpecificWorkflowMetadataRequest] = SetStoreSpecificWorkflowMetadata{}

// SetStoreSpecificWorkflowMetadata sets the workflows allowed to write specific data IDs.
// The data IDs need their feeds, and SVR paths their data ID configs (InitStoreFeeds).
type SetStoreSpecificWorkflowMetadata struct{}

func (cs SetStoreSpecificWorkflowMetadata) VerifyPreconditions(env cldf.Environment, req *SetStoreSpecificWorkflowMetadataRequest) error {
	if len(req.DataIDs) == 0 {
		return errors.New("DataIDs cannot be empty")
	}
	if _, err := parseDataIDs(req.DataIDs); err != nil {
		return err
	}
	if _, err := workflowMetadata(solana.PublicKey{}, req.Workflows); err != nil {
		return err
	}
	return req.verify(env)
}

func (cs SetStoreSpecificWorkflowMetadata) Apply(env cldf.Environment, req *SetStoreSpecificWorkflowMetadataRequest) (cldf.ChangesetOutput, error) {
	ids, err := parseDataIDs(req.DataIDs)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	ixs, err := buildBatched(env, req.StoreTarget, ids,
		func(store storeAddresses, feedAdmin solana.PublicKey, batch [][16]byte) (solana.Instruction, error) {
			settings, err := store.settingsMetas(batch)
			if err != nil {
				return nil, err
			}
			build := storeInstruction(store.ProgramID, settings)
			if len(req.Workflows) == 0 {
				return build(df_store.NewRemoveSpecificFeedMetadataInstruction(batch, feedAdmin, store.State, solana.SystemProgramID))
			}
			metadata, err := workflowMetadata(store.ProgramID, req.Workflows)
			if err != nil {
				return nil, err
			}
			return build(df_store.NewSetSpecificFeedMetadataInstruction(batch, metadata, feedAdmin, store.State, solana.SystemProgramID))
		})
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	return req.StoreTarget.send(env, "set data feeds store specific workflow metadata", ixs)
}

// ----------------------------------------------------------------------------------------------
// Feeds and data IDs
// ----------------------------------------------------------------------------------------------

type InitStoreFeedsRequest struct {
	StoreTarget
	// Data IDs, non-SVR or SVR. The non-SVR data ID of each gets a feed, which also holds its
	// settings, and each SVR path a data ID config.
	DataIDs []string
}

var _ cldf.ChangeSetV2[*InitStoreFeedsRequest] = InitStoreFeeds{}

// InitStoreFeeds creates the feed accounts and SVR path data ID configs reports need. The feed
// admin pays their rent. Existing accounts are left as they are.
type InitStoreFeeds struct{}

func (cs InitStoreFeeds) VerifyPreconditions(env cldf.Environment, req *InitStoreFeedsRequest) error {
	if len(req.DataIDs) == 0 {
		return errors.New("DataIDs cannot be empty")
	}
	if _, err := parseDataIDs(req.DataIDs); err != nil {
		return err
	}
	return req.verify(env)
}

func (cs InitStoreFeeds) Apply(env cldf.Environment, req *InitStoreFeedsRequest) (cldf.ChangesetOutput, error) {
	ids, err := parseDataIDs(req.DataIDs)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	var feeds, svrPaths [][16]byte
	seen := map[[16]byte]bool{}
	for _, id := range ids {
		nonSvr := nonSvrDataID(id)
		if !seen[nonSvr] {
			seen[nonSvr] = true
			feeds = append(feeds, nonSvr)
		}
		if isSvrDataID(id) {
			svrPaths = append(svrPaths, id)
		}
	}

	ixs, err := buildBatched(env, req.StoreTarget, feeds,
		func(store storeAddresses, feedAdmin solana.PublicKey, batch [][16]byte) (solana.Instruction, error) {
			accounts, err := store.metas(feedSeed, idKeys(batch), true)
			if err != nil {
				return nil, err
			}
			return storeInstruction(store.ProgramID, accounts)(df_store.NewInitFeedsInstruction(batch, feedAdmin, store.State, solana.SystemProgramID))
		})
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	more, err := buildBatched(env, req.StoreTarget, svrPaths,
		func(store storeAddresses, feedAdmin solana.PublicKey, batch [][16]byte) (solana.Instruction, error) {
			accounts, err := store.metas(dataIDSeed, idKeys(batch), true)
			if err != nil {
				return nil, err
			}
			return storeInstruction(store.ProgramID, accounts)(df_store.NewInitDataIdsInstruction(batch, feedAdmin, store.State, solana.SystemProgramID))
		})
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	return req.send(env, "init data feeds store feeds and data IDs", append(ixs, more...))
}

type StoreMaxStoredReports struct {
	DataID           string // non-SVR
	MaxStoredReports uint64
}

type SetStoreMaxStoredReportsRequest struct {
	StoreTarget
	Feeds []StoreMaxStoredReports
}

var _ cldf.ChangeSetV2[*SetStoreMaxStoredReportsRequest] = SetStoreMaxStoredReports{}

// SetStoreMaxStoredReports sets feed ring sizes. The first configuration is free, after that a
// ring only grows, and only while it has never wrapped.
type SetStoreMaxStoredReports struct{}

func (cs SetStoreMaxStoredReports) VerifyPreconditions(env cldf.Environment, req *SetStoreMaxStoredReportsRequest) error {
	if len(req.Feeds) == 0 {
		return errors.New("feeds cannot be empty")
	}
	if _, err := nonSvrIDs(req.Feeds, func(f StoreMaxStoredReports) string { return f.DataID }); err != nil {
		return err
	}
	for _, f := range req.Feeds {
		if f.MaxStoredReports < minMaxStoredReports {
			return fmt.Errorf("%s: MaxStoredReports must be at least %d", f.DataID, minMaxStoredReports)
		}
	}
	return req.verify(env)
}

func (cs SetStoreMaxStoredReports) Apply(env cldf.Environment, req *SetStoreMaxStoredReportsRequest) (cldf.ChangesetOutput, error) {
	ixs, err := buildBatched(env, req.StoreTarget, req.Feeds,
		func(store storeAddresses, feedAdmin solana.PublicKey, batch []StoreMaxStoredReports) (solana.Instruction, error) {
			ids, err := nonSvrIDs(batch, func(f StoreMaxStoredReports) string { return f.DataID })
			if err != nil {
				return nil, err
			}
			sizes := make([]uint64, len(batch))
			for i, f := range batch {
				sizes[i] = f.MaxStoredReports
			}
			accounts, err := store.metas(feedSeed, idKeys(ids), true)
			if err != nil {
				return nil, err
			}
			return storeInstruction(store.ProgramID, accounts)(df_store.NewSetMaxStoredReportsInstruction(ids, sizes, feedAdmin, store.State, solana.SystemProgramID))
		})
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	return req.StoreTarget.send(env, "set data feeds store max stored reports", ixs)
}

type StoreMinMaxAnswer struct {
	DataID string // non-SVR
	// Accepted answer range at 18 decimals, base 10
	Min, Max string
}

type SetStoreMinMaxAnswersRequest struct {
	StoreTarget
	Feeds []StoreMinMaxAnswer
}

var _ cldf.ChangeSetV2[*SetStoreMinMaxAnswersRequest] = SetStoreMinMaxAnswers{}

// SetStoreMinMaxAnswers sets the answer range a feed accepts, reports outside it are skipped.
type SetStoreMinMaxAnswers struct{}

func (cs SetStoreMinMaxAnswers) VerifyPreconditions(env cldf.Environment, req *SetStoreMinMaxAnswersRequest) error {
	if len(req.Feeds) == 0 {
		return errors.New("feeds cannot be empty")
	}
	if _, err := nonSvrIDs(req.Feeds, func(f StoreMinMaxAnswer) string { return f.DataID }); err != nil {
		return err
	}
	for _, f := range req.Feeds {
		minAnswer, err := parseAnswer(f.Min)
		if err != nil {
			return err
		}
		maxAnswer, err := parseAnswer(f.Max)
		if err != nil {
			return err
		}
		if minAnswer.Cmp(maxAnswer) > 0 {
			return fmt.Errorf("%s: min answer %s above max answer %s", f.DataID, f.Min, f.Max)
		}
	}
	return req.verify(env)
}

func (cs SetStoreMinMaxAnswers) Apply(env cldf.Environment, req *SetStoreMinMaxAnswersRequest) (cldf.ChangesetOutput, error) {
	ixs, err := buildBatched(env, req.StoreTarget, req.Feeds,
		func(store storeAddresses, feedAdmin solana.PublicKey, batch []StoreMinMaxAnswer) (solana.Instruction, error) {
			ids, err := nonSvrIDs(batch, func(f StoreMinMaxAnswer) string { return f.DataID })
			if err != nil {
				return nil, err
			}
			mins := make([]bin.Int128, len(batch))
			maxs := make([]bin.Int128, len(batch))
			for i, f := range batch {
				minAnswer, err := parseAnswer(f.Min)
				if err != nil {
					return nil, err
				}
				maxAnswer, err := parseAnswer(f.Max)
				if err != nil {
					return nil, err
				}
				mins[i], maxs[i] = toInt128(minAnswer), toInt128(maxAnswer)
			}
			accounts, err := store.metas(feedSeed, idKeys(ids), true)
			if err != nil {
				return nil, err
			}
			return storeInstruction(store.ProgramID, accounts)(df_store.NewSetMinMaxAnswersInstruction(ids, mins, maxs, feedAdmin, store.State, solana.SystemProgramID))
		})
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	return req.StoreTarget.send(env, "set data feeds store min/max answers", ixs)
}

type StoreSvrDisabled struct {
	DataID   string // SVR
	Disabled bool
}

type SetStoreSvrDisabledRequest struct {
	StoreTarget
	Paths []StoreSvrDisabled
}

var _ cldf.ChangeSetV2[*SetStoreSvrDisabledRequest] = SetStoreSvrDisabled{}

// SetStoreSvrDisabled disables or re-enables SVR paths. A disabled path reads the non-SVR feed.
type SetStoreSvrDisabled struct{}

func (cs SetStoreSvrDisabled) VerifyPreconditions(env cldf.Environment, req *SetStoreSvrDisabledRequest) error {
	if len(req.Paths) == 0 {
		return errors.New("paths cannot be empty")
	}
	if _, err := svrIDs(req.Paths); err != nil {
		return err
	}
	return req.verify(env)
}

func (cs SetStoreSvrDisabled) Apply(env cldf.Environment, req *SetStoreSvrDisabledRequest) (cldf.ChangesetOutput, error) {
	ixs, err := buildBatched(env, req.StoreTarget, req.Paths,
		func(store storeAddresses, feedAdmin solana.PublicKey, batch []StoreSvrDisabled) (solana.Instruction, error) {
			ids, err := svrIDs(batch)
			if err != nil {
				return nil, err
			}
			disabled := make([]bool, len(batch))
			for i, p := range batch {
				disabled[i] = p.Disabled
			}
			accounts, err := store.metas(dataIDSeed, idKeys(ids), true)
			if err != nil {
				return nil, err
			}
			return storeInstruction(store.ProgramID, accounts)(df_store.NewSetSvrDisabledInstruction(ids, disabled, feedAdmin, store.State, solana.SystemProgramID))
		})
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	return req.StoreTarget.send(env, "set data feeds store SVR disabled", ixs)
}

// ----------------------------------------------------------------------------------------------
// Proxies and descriptions
// ----------------------------------------------------------------------------------------------

type StoreProxyConfig struct {
	ProxyID           string // 32 bytes, hex
	DataID            string // non-SVR or SVR
	Decimals          uint8
	FallbackTime      uint32 // at most the feed's ring size
	MaxSyncIterations uint32
}

type SetStoreProxyConfigsRequest struct {
	StoreTarget
	Proxies        []StoreProxyConfig
	RemoveProxyIDs []string
}

var _ cldf.ChangeSetV2[*SetStoreProxyConfigsRequest] = SetStoreProxyConfigs{}

// SetStoreProxyConfigs creates, updates or removes proxy (reader) configs.
type SetStoreProxyConfigs struct{}

func (cs SetStoreProxyConfigs) VerifyPreconditions(env cldf.Environment, req *SetStoreProxyConfigsRequest) error {
	if len(req.Proxies)+len(req.RemoveProxyIDs) == 0 {
		return errors.New("no proxy configs to update or remove")
	}
	seen := map[[32]byte]bool{}
	for _, p := range req.Proxies {
		proxyID, args, err := p.parse()
		if err != nil {
			return err
		}
		if args.Decimals > storedDecimals {
			return fmt.Errorf("proxy %s: decimals %d above %d", p.ProxyID, p.Decimals, storedDecimals)
		}
		if seen[proxyID] {
			return fmt.Errorf("proxy %s listed twice", p.ProxyID)
		}
		seen[proxyID] = true
	}
	for _, s := range req.RemoveProxyIDs {
		proxyID, err := parseProxyID(s)
		if err != nil {
			return err
		}
		if seen[proxyID] {
			return fmt.Errorf("proxy %s listed twice", s)
		}
		seen[proxyID] = true
	}
	return req.verify(env)
}

func (cs SetStoreProxyConfigs) Apply(env cldf.Environment, req *SetStoreProxyConfigsRequest) (cldf.ChangesetOutput, error) {
	ixs, err := buildBatched(env, req.StoreTarget, req.Proxies,
		func(store storeAddresses, feedAdmin solana.PublicKey, batch []StoreProxyConfig) (solana.Instruction, error) {
			proxyIDs := make([][32]byte, len(batch))
			configs := make([]df_store.ProxyConfigArgs, len(batch))
			var accounts []*solana.AccountMeta
			for i, p := range batch {
				proxyID, args, err := p.parse()
				if err != nil {
					return nil, err
				}
				proxyIDs[i], configs[i] = proxyID, args
				proxy, err := storePDA(store.ProgramID, store.State, proxySeed, proxyID[:])
				if err != nil {
					return nil, err
				}
				feed, err := store.feed(args.DataId)
				if err != nil {
					return nil, err
				}
				accounts = append(accounts, solana.Meta(proxy).WRITE(), solana.Meta(feed))
			}
			return storeInstruction(store.ProgramID, accounts)(df_store.NewUpdateConfigsForProxiesInstruction(proxyIDs, configs, feedAdmin, store.State, solana.SystemProgramID))
		})
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	more, err := buildBatched(env, req.StoreTarget, req.RemoveProxyIDs,
		func(store storeAddresses, feedAdmin solana.PublicKey, batch []string) (solana.Instruction, error) {
			proxyIDs := make([][32]byte, len(batch))
			keys := make([][]byte, len(batch))
			for i, s := range batch {
				proxyID, err := parseProxyID(s)
				if err != nil {
					return nil, err
				}
				proxyIDs[i] = proxyID
				keys[i] = proxyIDs[i][:]
			}
			accounts, err := store.metas(proxySeed, keys, true)
			if err != nil {
				return nil, err
			}
			return storeInstruction(store.ProgramID, accounts)(df_store.NewRemoveConfigsForProxiesInstruction(proxyIDs, feedAdmin, store.State, solana.SystemProgramID))
		})
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	return req.send(env, "update data feeds store proxy configs", append(ixs, more...))
}

func parseProxyID(s string) ([32]byte, error) {
	var out [32]byte
	b, err := parseHex(s, len(out))
	if err != nil {
		return out, fmt.Errorf("invalid proxy ID %q: %w", s, err)
	}
	copy(out[:], b)
	if out == ([32]byte{}) {
		return out, errors.New("proxy ID cannot be zero")
	}
	return out, nil
}

func (p StoreProxyConfig) parse() ([32]byte, df_store.ProxyConfigArgs, error) {
	proxyID, err := parseProxyID(p.ProxyID)
	if err != nil {
		return proxyID, df_store.ProxyConfigArgs{}, err
	}
	dataID, err := ParseDataID(p.DataID)
	if err != nil {
		return proxyID, df_store.ProxyConfigArgs{}, err
	}
	if dataID == ([16]byte{}) {
		return proxyID, df_store.ProxyConfigArgs{}, errors.New("data ID cannot be zero")
	}
	return proxyID, df_store.ProxyConfigArgs{
		DataId:            dataID,
		Decimals:          p.Decimals,
		FallbackTime:      p.FallbackTime,
		MaxSyncIterations: p.MaxSyncIterations,
	}, nil
}

type StoreFeedDescription struct {
	DataID      string // any data ID of the feed, the description is keyed by its indexed feed ID
	Description string // at most 64 bytes
}

type StoreSvrPathDescription struct {
	SvrPathID   string // 2 bytes, hex
	Description string // at most 64 bytes
}

type SetStoreDescriptionsRequest struct {
	StoreTarget
	Feeds    []StoreFeedDescription
	SvrPaths []StoreSvrPathDescription
}

var _ cldf.ChangeSetV2[*SetStoreDescriptionsRequest] = SetStoreDescriptions{}

// SetStoreDescriptions sets feed and SVR path descriptions.
type SetStoreDescriptions struct{}

func (cs SetStoreDescriptions) VerifyPreconditions(env cldf.Environment, req *SetStoreDescriptionsRequest) error {
	if len(req.Feeds)+len(req.SvrPaths) == 0 {
		return errors.New("no descriptions to set")
	}
	if _, err := feedDescriptionIDs(req.Feeds); err != nil {
		return err
	}
	if _, err := svrPathIDs(req.SvrPaths); err != nil {
		return err
	}
	return req.verify(env)
}

func (cs SetStoreDescriptions) Apply(env cldf.Environment, req *SetStoreDescriptionsRequest) (cldf.ChangesetOutput, error) {
	ixs, err := buildBatched(env, req.StoreTarget, req.Feeds,
		func(store storeAddresses, feedAdmin solana.PublicKey, batch []StoreFeedDescription) (solana.Instruction, error) {
			ids, err := feedDescriptionIDs(batch)
			if err != nil {
				return nil, err
			}
			descriptions := make([]string, len(batch))
			for i, f := range batch {
				descriptions[i] = f.Description
			}
			accounts, err := store.metas(feedDescriptionSeed, idKeys(ids), true)
			if err != nil {
				return nil, err
			}
			return storeInstruction(store.ProgramID, accounts)(df_store.NewSetFeedDescriptionsInstruction(ids, descriptions, feedAdmin, store.State, solana.SystemProgramID))
		})
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	more, err := buildBatched(env, req.StoreTarget, req.SvrPaths,
		func(store storeAddresses, feedAdmin solana.PublicKey, batch []StoreSvrPathDescription) (solana.Instruction, error) {
			ids, err := svrPathIDs(batch)
			if err != nil {
				return nil, err
			}
			descriptions := make([]string, len(batch))
			keys := make([][]byte, len(batch))
			for i, p := range batch {
				descriptions[i] = p.Description
				keys[i] = ids[i][:]
			}
			accounts, err := store.metas(svrPathDescriptionSeed, keys, true)
			if err != nil {
				return nil, err
			}
			return storeInstruction(store.ProgramID, accounts)(df_store.NewSetSvrPathDescriptionsInstruction(ids, descriptions, feedAdmin, store.State, solana.SystemProgramID))
		})
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	return req.send(env, "set data feeds store descriptions", append(ixs, more...))
}

func feedDescriptionIDs(feeds []StoreFeedDescription) ([][16]byte, error) {
	out := make([][16]byte, len(feeds))
	seen := map[[16]byte]bool{}
	for i, f := range feeds {
		id, err := ParseDataID(f.DataID)
		if err != nil {
			return nil, err
		}
		if len(f.Description) > maxDescriptionLen {
			return nil, fmt.Errorf("%s: description longer than %d bytes", f.DataID, maxDescriptionLen)
		}
		out[i] = indexedFeedID(id)
		if out[i] == ([16]byte{}) {
			return nil, fmt.Errorf("%s: indexed feed ID is zero", f.DataID)
		}
		if seen[out[i]] {
			return nil, fmt.Errorf("%s: feed described twice", f.DataID)
		}
		seen[out[i]] = true
	}
	return out, nil
}

func svrPathIDs(paths []StoreSvrPathDescription) ([][2]byte, error) {
	out := make([][2]byte, len(paths))
	seen := map[[2]byte]bool{}
	for i, p := range paths {
		b, err := parseHex(p.SvrPathID, 2)
		if err != nil {
			return nil, fmt.Errorf("invalid SVR path ID %q: %w", p.SvrPathID, err)
		}
		copy(out[i][:], b)
		if out[i] == ([2]byte{}) {
			return nil, errors.New("SVR path ID cannot be zero")
		}
		if len(p.Description) > maxDescriptionLen {
			return nil, fmt.Errorf("SVR path %s: description longer than %d bytes", p.SvrPathID, maxDescriptionLen)
		}
		if seen[out[i]] {
			return nil, fmt.Errorf("SVR path %s described twice", p.SvrPathID)
		}
		seen[out[i]] = true
	}
	return out, nil
}

func nonSvrIDs[T any](items []T, dataID func(T) string) ([][16]byte, error) {
	strs := make([]string, len(items))
	for i, item := range items {
		strs[i] = dataID(item)
	}
	ids, err := parseDataIDs(strs)
	if err != nil {
		return nil, err
	}
	for i, id := range ids {
		if isSvrDataID(id) {
			return nil, fmt.Errorf("%s is an SVR data ID, expected a non-SVR one", strs[i])
		}
	}
	return ids, nil
}

func svrIDs(paths []StoreSvrDisabled) ([][16]byte, error) {
	strs := make([]string, len(paths))
	for i, p := range paths {
		strs[i] = p.DataID
	}
	ids, err := parseDataIDs(strs)
	if err != nil {
		return nil, err
	}
	for i, id := range ids {
		if !isSvrDataID(id) {
			return nil, fmt.Errorf("%s is not an SVR data ID", strs[i])
		}
	}
	return ids, nil
}

// ----------------------------------------------------------------------------------------------
// Report callbacks
// ----------------------------------------------------------------------------------------------

type StoreReportCallback struct {
	DataID string
	// Program called with every report applied for the data ID, zero disables the callback
	Program         solana.PublicKey
	Accounts        []df_store.CallbackAccount // at most 4
	RevertOnFailure bool
}

type SetStoreReportCallbacksRequest struct {
	StoreTarget
	Callbacks []StoreReportCallback
}

var _ cldf.ChangeSetV2[*SetStoreReportCallbacksRequest] = SetStoreReportCallbacks{}

// SetStoreReportCallbacks configures, per data ID, the program called before a report is stored.
type SetStoreReportCallbacks struct{}

func (cs SetStoreReportCallbacks) VerifyPreconditions(env cldf.Environment, req *SetStoreReportCallbacksRequest) error {
	if len(req.Callbacks) == 0 {
		return errors.New("callbacks cannot be empty")
	}
	if err := verifyCallbacks(req.Callbacks); err != nil {
		return err
	}
	return req.verify(env)
}

func (cs SetStoreReportCallbacks) Apply(env cldf.Environment, req *SetStoreReportCallbacksRequest) (cldf.ChangesetOutput, error) {
	return applyCallbacks(env, req.StoreTarget, req.Callbacks)
}

// DF1Feed is a DF1 store program feed (Transmissions account) and the non-SVR data ID that
// writes to it.
type DF1Feed struct {
	DataID string // non-SVR
	Feed   solana.PublicKey
}

type ConfigureDF1ReportCallbacksRequest struct {
	StoreTarget
	DF1StoreProgramID solana.PublicKey
	Feeds             []DF1Feed
}

var _ cldf.ChangeSetV2[*ConfigureDF1ReportCallbacksRequest] = ConfigureDF1ReportCallbacks{}

// ConfigureDF1ReportCallbacks points each feed's non-SVR data ID at the DF1 store program's
// on_report_callback with its DF1 feed account, so reports are written through to the DF1 feed.
// The DF1 feed only takes them once its writer is the store authority (CutoverDF1Feeds).
type ConfigureDF1ReportCallbacks struct{}

func (cs ConfigureDF1ReportCallbacks) VerifyPreconditions(env cldf.Environment, req *ConfigureDF1ReportCallbacksRequest) error {
	if req.DF1StoreProgramID.IsZero() {
		return errors.New("DF1StoreProgramID cannot be zero")
	}
	if err := verifyDF1Feeds(req.Feeds); err != nil {
		return err
	}
	return req.verify(env)
}

func (cs ConfigureDF1ReportCallbacks) Apply(env cldf.Environment, req *ConfigureDF1ReportCallbacksRequest) (cldf.ChangesetOutput, error) {
	return applyCallbacks(env, req.StoreTarget, df1Callbacks(req.DF1StoreProgramID, req.Feeds))
}

func df1Callbacks(df1StoreProgram solana.PublicKey, feeds []DF1Feed) []StoreReportCallback {
	out := make([]StoreReportCallback, len(feeds))
	for i, f := range feeds {
		out[i] = StoreReportCallback{
			DataID:   f.DataID,
			Program:  df1StoreProgram,
			Accounts: []df_store.CallbackAccount{{Pubkey: f.Feed, IsWritable: true}},
			// rounds the DF1 feed can't take are skipped by the callback, never rejected
			RevertOnFailure: false,
		}
	}
	return out
}

func verifyDF1Feeds(feeds []DF1Feed) error {
	if len(feeds) == 0 {
		return errors.New("feeds cannot be empty")
	}
	if _, err := nonSvrIDs(feeds, func(f DF1Feed) string { return f.DataID }); err != nil {
		return err
	}
	seen := map[solana.PublicKey]bool{}
	for _, f := range feeds {
		if f.Feed.IsZero() {
			return fmt.Errorf("%s: DF1 feed cannot be zero", f.DataID)
		}
		if seen[f.Feed] {
			return fmt.Errorf("DF1 feed %s listed twice", f.Feed)
		}
		seen[f.Feed] = true
	}
	return nil
}

func verifyCallbacks(callbacks []StoreReportCallback) error {
	strs := make([]string, len(callbacks))
	for i, c := range callbacks {
		strs[i] = c.DataID
		if len(c.Accounts) > maxCallbackAccounts {
			return fmt.Errorf("%s: at most %d callback accounts", c.DataID, maxCallbackAccounts)
		}
		if c.Program.IsZero() && (len(c.Accounts) > 0 || c.RevertOnFailure) {
			return fmt.Errorf("%s: a disabled callback takes no accounts", c.DataID)
		}
	}
	_, err := parseDataIDs(strs)
	return err
}

func applyCallbacks(env cldf.Environment, t StoreTarget, callbacks []StoreReportCallback) (cldf.ChangesetOutput, error) {
	ixs, err := buildBatched(env, t, callbacks,
		func(store storeAddresses, feedAdmin solana.PublicKey, batch []StoreReportCallback) (solana.Instruction, error) {
			ids := make([][16]byte, len(batch))
			args := make([]df_store.CallbackConfigArgs, len(batch))
			var programs []*solana.AccountMeta
			for i, c := range batch {
				id, err := ParseDataID(c.DataID)
				if err != nil {
					return nil, err
				}
				ids[i] = id
				args[i] = df_store.CallbackConfigArgs{
					Program:         c.Program,
					Accounts:        c.Accounts,
					RevertOnFailure: c.RevertOnFailure,
				}
				if args[i].Accounts == nil {
					args[i].Accounts = []df_store.CallbackAccount{}
				}
				if !c.Program.IsZero() {
					programs = append(programs, solana.Meta(c.Program))
				}
			}
			settings, err := store.settingsMetas(ids)
			if err != nil {
				return nil, err
			}
			return storeInstruction(store.ProgramID, append(settings, programs...))(
				df_store.NewSetReportCallbackConfigsInstruction(ids, args, feedAdmin, store.State, solana.SystemProgramID))
		})
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	return t.send(env, "set data feeds store report callbacks", ixs)
}
