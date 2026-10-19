package solana

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/Masterminds/semver/v3"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	df_store "github.com/smartcontractkit/chainlink-solana/contracts/generated/data_feeds_store"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/solana/sequence/operation"
	"github.com/smartcontractkit/chainlink/deployment/helpers"
)

// Data feeds store (programs/data-feeds-store in chainlink-solana), the Solana port of the EVM
// DataFeedsStore. It replaces the never-deployed data feeds cache.
const (
	StoreContract datastore.ContractType = "DataFeedsStore"
	StoreState    datastore.ContractType = "DataFeedsStoreState"
)

const (
	// storedDecimals is the precision of every answer in the data feeds store.
	storedDecimals = 18
	// maxDefaultStoredReports keeps a new feed account within the 10 KiB account creation limit:
	// 8 byte discriminator + 1,376 byte header + 32 bytes per stored report.
	maxDefaultStoredReports = (10_240 - 8 - 1_376) / 32
	minMaxStoredReports     = 2
	// maxDescriptionLen is the byte length of a feed or SVR path description.
	maxDescriptionLen = 64
	// maxCallbackAccounts is the number of accounts a report callback may take.
	maxCallbackAccounts = 4
	// maxWorkflowMetadatas is the number of workflows a default or specific allow list holds.
	maxWorkflowMetadatas = 16
	// storeBatchSize caps the data IDs configured per instruction, so a transaction (or an MCMS
	// operation) with its remaining accounts stays under the 1232 byte limit.
	storeBatchSize = 5
)

// PDA seeds of the data feeds store.
const (
	storeAuthoritySeed     = "store_authority"
	feedSeed               = "feed"
	dataIDSeed             = "data_id"
	proxySeed              = "proxy"
	feedDescriptionSeed    = "feed_description"
	svrPathDescriptionSeed = "svr_path_description"
)

// Sender identifies a keystone forwarder instance: its program and its state account.
type Sender struct {
	ProgramID solana.PublicKey
	StateID   solana.PublicKey
}

// StoreWorkflow is a workflow allowed to write to the data feeds store.
type StoreWorkflow struct {
	// Forwarder the workflow reports through, its forwarder authority for the store is the allowed sender
	Forwarder Sender
	Owner     [20]uint8
	Name      [10]uint8
}

// storeRef identifies a data feeds store deployment in the datastore.
type storeRef struct {
	ChainSel  uint64
	Version   string
	Qualifier string
}

type storeAddresses struct {
	ProgramID solana.PublicKey
	State     solana.PublicKey
}

func loadStore(ds datastore.DataStore, ref storeRef) (storeAddresses, error) {
	version, err := semver.NewVersion(ref.Version)
	if err != nil {
		return storeAddresses{}, err
	}
	program, err := ds.Addresses().Get(datastore.NewAddressRefKey(ref.ChainSel, StoreContract, version, ref.Qualifier))
	if err != nil {
		return storeAddresses{}, fmt.Errorf("failed to load data feeds store program: %w", err)
	}
	state, err := ds.Addresses().Get(datastore.NewAddressRefKey(ref.ChainSel, StoreState, version, ref.Qualifier))
	if err != nil {
		return storeAddresses{}, fmt.Errorf("failed to load data feeds store state: %w", err)
	}
	programID, err := solana.PublicKeyFromBase58(program.Address)
	if err != nil {
		return storeAddresses{}, fmt.Errorf("invalid data feeds store program address: %w", err)
	}
	stateID, err := solana.PublicKeyFromBase58(state.Address)
	if err != nil {
		return storeAddresses{}, fmt.Errorf("invalid data feeds store state address: %w", err)
	}
	return storeAddresses{ProgramID: programID, State: stateID}, nil
}

// verifyStoreRequest checks the chain and the store deployment exist, and that a timelock is
// deployed when the change goes through MCMS.
func verifyStoreRequest(env cldf.Environment, ref storeRef, mcmsCfg *cldfproposalutils.TimelockConfig) error {
	if _, ok := env.BlockChains.SolanaChains()[ref.ChainSel]; !ok {
		return fmt.Errorf("solana chain not found for chain selector %d", ref.ChainSel)
	}
	if _, err := loadStore(env.DataStore, ref); err != nil {
		return err
	}
	return verifyMCMS(env, ref.ChainSel, mcmsCfg)
}

func verifyMCMS(env cldf.Environment, chainSel uint64, mcmsCfg *cldfproposalutils.TimelockConfig) error {
	if mcmsCfg == nil {
		return nil
	}
	if _, err := helpers.FetchTimelockSigner(env.DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(chainSel))); err != nil {
		return fmt.Errorf("failed to fetch timelock signer: %w", err)
	}
	return nil
}

// authority returns the account that signs: the timelock signer under MCMS, else the deployer.
func authority(env cldf.Environment, chainSel uint64, mcmsCfg *cldfproposalutils.TimelockConfig) (solana.PublicKey, error) {
	if mcmsCfg == nil {
		return env.BlockChains.SolanaChains()[chainSel].DeployerKey.PublicKey(), nil
	}
	return helpers.FetchTimelockSigner(env.DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(chainSel)))
}

// sendInstructions confirms the instructions with the deployer, or builds a timelock proposal
// for them under MCMS.
func sendInstructions(
	env cldf.Environment,
	chainSel uint64,
	mcmsCfg *cldfproposalutils.TimelockConfig,
	description string,
	ixs []solana.Instruction,
) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput
	if len(ixs) == 0 {
		env.Logger.Infow("nothing to do", "description", description)
		return out, nil
	}

	input := operation.SendInstructionsInput{
		ChainSel:     chainSel,
		MCMS:         mcmsCfg,
		ContractType: cldf.ContractType(StoreContract),
		Description:  description,
	}
	for _, ix := range ixs {
		data, err := ix.Data()
		if err != nil {
			return out, fmt.Errorf("failed to encode instruction: %w", err)
		}
		input.Instructions = append(input.Instructions, operation.Instruction{
			ProgramID: ix.ProgramID(),
			Accounts:  ix.Accounts(),
			Data:      data,
		})
	}

	deps := operation.Deps{
		Datastore: env.DataStore,
		Env:       env,
		Chain:     env.BlockChains.SolanaChains()[chainSel],
	}
	report, err := operations.ExecuteOperation(env.OperationsBundle, operation.SendInstructionsOp, deps, input)
	if err != nil {
		return out, err
	}
	out.MCMSTimelockProposals = report.Output.Proposals
	return out, nil
}

// storeInstruction retargets a generated instruction at the deployed store program and appends
// its remaining accounts.
func storeInstruction(programID solana.PublicKey, remaining []*solana.AccountMeta) func(solana.Instruction, error) (solana.Instruction, error) {
	return func(ix solana.Instruction, err error) (solana.Instruction, error) {
		if err != nil {
			return nil, err
		}
		data, err := ix.Data()
		if err != nil {
			return nil, err
		}
		accounts := append(solana.AccountMetaSlice(ix.Accounts()), remaining...)
		return solana.NewInstruction(programID, accounts, data), nil
	}
}

// ParseDataID parses a 16 byte data ID written as 32 hex characters, with or without 0x.
func ParseDataID(s string) ([16]byte, error) {
	var out [16]byte
	b, err := parseHex(s, len(out))
	if err != nil {
		return out, fmt.Errorf("invalid data ID %q: %w", s, err)
	}
	copy(out[:], b)
	return out, nil
}

func parseDataIDs(ids []string) ([][16]byte, error) {
	out := make([][16]byte, 0, len(ids))
	seen := make(map[[16]byte]bool, len(ids))
	for _, s := range ids {
		id, err := ParseDataID(s)
		if err != nil {
			return nil, err
		}
		if id == ([16]byte{}) {
			return nil, errors.New("data ID cannot be zero")
		}
		if seen[id] {
			return nil, fmt.Errorf("duplicate data ID %s", s)
		}
		seen[id] = true
		out = append(out, id)
	}
	return out, nil
}

func parseHex(s string, n int) ([]byte, error) {
	b, err := hex.DecodeString(strings.TrimPrefix(strings.TrimPrefix(s, "0x"), "0X"))
	if err != nil {
		return nil, err
	}
	if len(b) != n {
		return nil, fmt.Errorf("expected %d bytes, got %d", n, len(b))
	}
	return b, nil
}

// Data ID layout, as DataId.sol:
// | CompactFeedId (0..8) | Index (8) | SvrPathId (9..11) | Instance (11) | Reserved (12..16) |

func isSvrDataID(id [16]byte) bool {
	return id[9] != 0 || id[10] != 0
}

func svrPathID(id [16]byte) [2]byte {
	return [2]byte{id[9], id[10]}
}

func nonSvrDataID(id [16]byte) [16]byte {
	id[9], id[10] = 0, 0
	return id
}

func indexedFeedID(id [16]byte) [16]byte {
	id = nonSvrDataID(id)
	id[11] = 0
	clear(id[12:])
	return id
}

func storePDA(programID, state solana.PublicKey, seed string, key []byte) (solana.PublicKey, error) {
	seeds := [][]byte{[]byte(seed), state.Bytes()}
	if key != nil {
		seeds = append(seeds, key)
	}
	pda, _, err := solana.FindProgramAddress(seeds, programID)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("failed to derive %s PDA: %w", seed, err)
	}
	return pda, nil
}

// StoreAuthority is the data feeds store PDA that signs report callbacks, the writer of cut over
// DF1 feeds.
func StoreAuthority(programID, state solana.PublicKey) (solana.PublicKey, error) {
	return storePDA(programID, state, storeAuthoritySeed, nil)
}

func (s storeAddresses) feed(id [16]byte) (solana.PublicKey, error) {
	nonSvr := nonSvrDataID(id)
	return storePDA(s.ProgramID, s.State, feedSeed, nonSvr[:])
}

func (s storeAddresses) dataIDConfig(id [16]byte) (solana.PublicKey, error) {
	return storePDA(s.ProgramID, s.State, dataIDSeed, id[:])
}

// settings is the account holding a data ID's settings (callback, feed-specific workflow
// metadata): its feed for a non-SVR data ID, its data ID config for an SVR path.
func (s storeAddresses) settings(id [16]byte) (solana.PublicKey, error) {
	if isSvrDataID(id) {
		return s.dataIDConfig(id)
	}
	return s.feed(id)
}

func (s storeAddresses) settingsMetas(ids [][16]byte) ([]*solana.AccountMeta, error) {
	out := make([]*solana.AccountMeta, 0, len(ids))
	for _, id := range ids {
		address, err := s.settings(id)
		if err != nil {
			return nil, err
		}
		out = append(out, solana.Meta(address).WRITE())
	}
	return out, nil
}

func (s storeAddresses) authority() (solana.PublicKey, error) {
	return StoreAuthority(s.ProgramID, s.State)
}

func (s storeAddresses) metas(seed string, keys [][]byte, writable bool) ([]*solana.AccountMeta, error) {
	out := make([]*solana.AccountMeta, 0, len(keys))
	for _, key := range keys {
		pda, err := storePDA(s.ProgramID, s.State, seed, key)
		if err != nil {
			return nil, err
		}
		meta := solana.Meta(pda)
		if writable {
			meta = meta.WRITE()
		}
		out = append(out, meta)
	}
	return out, nil
}

func idKeys(ids [][16]byte) [][]byte {
	out := make([][]byte, len(ids))
	for i := range ids {
		out[i] = ids[i][:]
	}
	return out
}

// workflowMetadata converts workflows to the store's metadata, the allowed sender being the
// forwarder authority PDA ["forwarder", forwarder_state, store_program] of each forwarder.
func workflowMetadata(storeProgram solana.PublicKey, workflows []StoreWorkflow) ([]df_store.WorkflowMetadata, error) {
	if len(workflows) > maxWorkflowMetadatas {
		return nil, fmt.Errorf("at most %d workflows, got %d", maxWorkflowMetadatas, len(workflows))
	}
	out := make([]df_store.WorkflowMetadata, len(workflows))
	for i, w := range workflows {
		if w.Forwarder.ProgramID.IsZero() || w.Forwarder.StateID.IsZero() {
			return nil, fmt.Errorf("workflow %d: forwarder program and state are required", i)
		}
		sender, err := deriveForwarderAuthority(w.Forwarder.StateID, storeProgram, w.Forwarder.ProgramID)
		if err != nil {
			return nil, fmt.Errorf("failed to derive forwarder authority: %w", err)
		}
		out[i] = df_store.WorkflowMetadata{
			AllowedSender:        sender,
			AllowedWorkflowOwner: w.Owner,
			AllowedWorkflowName:  w.Name,
		}
	}
	return out, nil
}

func deriveForwarderAuthority(forwarderState solana.PublicKey, receiverProgram solana.PublicKey, forwarderProgram solana.PublicKey) (solana.PublicKey, error) {
	seeds := [][]byte{
		[]byte("forwarder"),
		forwarderState[:],
		receiverProgram[:],
	}
	ret, _, err := solana.FindProgramAddress(seeds, forwarderProgram)
	return ret, err
}

var (
	maxInt128 = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 127), big.NewInt(1))
	minInt128 = new(big.Int).Neg(new(big.Int).Lsh(big.NewInt(1), 127))
)

// parseAnswer parses a base 10 answer at the store's 18 decimals.
func parseAnswer(s string) (*big.Int, error) {
	v, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return nil, fmt.Errorf("invalid answer %q", s)
	}
	if v.Cmp(maxInt128) > 0 || v.Cmp(minInt128) < 0 {
		return nil, fmt.Errorf("answer %s does not fit an i128", s)
	}
	return v, nil
}

func toInt128(v *big.Int) bin.Int128 {
	twos := new(big.Int).Set(v)
	if v.Sign() < 0 {
		twos.Add(twos, new(big.Int).Lsh(big.NewInt(1), 128))
	}
	mask := new(big.Int).SetUint64(^uint64(0))
	return bin.Int128{
		Lo: new(big.Int).And(twos, mask).Uint64(),
		Hi: new(big.Int).Rsh(twos, 64).Uint64(),
	}
}

func chunk[T any](s []T, size int) [][]T {
	var out [][]T
	for size < len(s) {
		out = append(out, s[:size])
		s = s[size:]
	}
	if len(s) > 0 {
		out = append(out, s)
	}
	return out
}
