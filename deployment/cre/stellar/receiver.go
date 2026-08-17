package stellar

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"

	cldfstellar "github.com/smartcontractkit/chainlink-deployments-framework/chain/stellar"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	stellardeployment "github.com/smartcontractkit/chainlink-stellar/deployment"
	stellarreceiver "github.com/smartcontractkit/chainlink-stellar/deployment/cre/receiver"
)

const ReceiverContract datastore.ContractType = "StellarReceiver"

var _ cldf.ChangeSetV2[*DeployReceiverRequest] = DeployReceiver{}
var _ cldf.ChangeSetV2[*ReportCountRequest] = ReportCount{}

type DeployReceiver struct{}

type DeployReceiverRequest struct {
	ChainSel  uint64
	Qualifier string
	Version   string
	LabelSet  datastore.LabelSet
	Salt      [32]byte
}

func (cs DeployReceiver) VerifyPreconditions(env cldf.Environment, req *DeployReceiverRequest) error {
	if req == nil {
		return errors.New("request is required")
	}
	if _, ok := env.BlockChains.StellarChains()[req.ChainSel]; !ok {
		return fmt.Errorf("stellar chain not found for chain selector %d", req.ChainSel)
	}
	if req.Qualifier == "" {
		return errors.New("receiver qualifier is required")
	}
	if req.Version == "" {
		return errors.New("receiver version is required")
	}
	if _, err := semver.NewVersion(req.Version); err != nil {
		return fmt.Errorf("invalid receiver version %q: %w", req.Version, err)
	}
	return nil
}

func (cs DeployReceiver) Apply(env cldf.Environment, req *DeployReceiverRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput
	out.DataStore = datastore.NewMemoryDataStore()

	version := semver.MustParse(req.Version)
	ch, ok := env.BlockChains.StellarChains()[req.ChainSel]
	if !ok {
		return out, fmt.Errorf("stellar chain not found for chain selector %d", req.ChainSel)
	}
	if ch.Signer == nil {
		return out, errors.New("stellar chain has no signer")
	}

	deployer, err := stellardeployment.NewDeployerFromChain(ch)
	if err != nil {
		return out, fmt.Errorf("failed to build stellar deployer for chain selector %d: %w", req.ChainSel, err)
	}

	wasm, err := Artifact(ReceiverWasm)
	if err != nil {
		return out, fmt.Errorf("failed to source receiver WASM: %w", err)
	}

	receiverAddr, err := stellarreceiver.DeployReceiver(env.GetContext(), deployer, wasm, req.Salt)
	if err != nil {
		return out, fmt.Errorf("failed to deploy stellar receiver on chain selector %d: %w", req.ChainSel, err)
	}

	if err := out.DataStore.Addresses().Add(datastore.AddressRef{
		Address:       receiverAddr,
		ChainSelector: req.ChainSel,
		Type:          ReceiverContract,
		Version:       version,
		Qualifier:     req.Qualifier,
		Labels:        req.LabelSet,
	}); err != nil && !errors.Is(err, datastore.ErrAddressRefExists) {
		return out, fmt.Errorf("failed to add stellar receiver address to datastore: %w", err)
	}

	env.Logger.Infow("Deployed Stellar CRE receiver", "chainSelector", req.ChainSel, "receiver", receiverAddr, "qualifier", req.Qualifier)

	return out, nil
}

type ReportCount struct{}

type ReportCountRequest struct {
	ChainSel  uint64
	Qualifier string
	Version   string
}

func (cs ReportCount) VerifyPreconditions(env cldf.Environment, req *ReportCountRequest) error {
	if req == nil {
		return errors.New("request is required")
	}
	if _, ok := env.BlockChains.StellarChains()[req.ChainSel]; !ok {
		return fmt.Errorf("stellar chain not found for chain selector %d", req.ChainSel)
	}
	if req.Qualifier == "" {
		return errors.New("receiver qualifier is required")
	}
	if req.Version == "" {
		return errors.New("receiver version is required")
	}
	if _, err := semver.NewVersion(req.Version); err != nil {
		return fmt.Errorf("invalid receiver version %q: %w", req.Version, err)
	}
	return nil
}

func (cs ReportCount) Apply(env cldf.Environment, req *ReportCountRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput

	version := semver.MustParse(req.Version)
	ch, ok := env.BlockChains.StellarChains()[req.ChainSel]
	if !ok {
		return out, fmt.Errorf("stellar chain not found for chain selector %d", req.ChainSel)
	}

	refKey := datastore.NewAddressRefKey(req.ChainSel, ReceiverContract, version, req.Qualifier)
	addrRef, err := env.DataStore.Addresses().Get(refKey)
	if err != nil {
		return out, fmt.Errorf("failed to get stellar receiver for ref key %s: %w", refKey, err)
	}

	count, err := ReportCountForChain(env.GetContext(), ch, addrRef.Address)
	if err != nil {
		return out, fmt.Errorf("failed to read report_count on stellar receiver %s (chain selector %d): %w", addrRef.Address, req.ChainSel, err)
	}

	env.Logger.Infow("Read Stellar CRE receiver report_count", "chainSelector", req.ChainSel, "receiver", addrRef.Address, "count", count)

	return out, nil
}

// DeployReceiverForChain deploys the CRE test receiver directly from a CLDF
// stellar chain. It is exposed for test helpers that do not have a full
// cldf.Environment available.
func DeployReceiverForChain(ctx context.Context, ch cldfstellar.Chain, salt [32]byte) (string, error) {
	deployer, err := stellardeployment.NewDeployerFromChain(ch)
	if err != nil {
		return "", fmt.Errorf("failed to build stellar deployer: %w", err)
	}

	wasm, err := Artifact(ReceiverWasm)
	if err != nil {
		return "", fmt.Errorf("failed to source receiver WASM: %w", err)
	}

	return stellarreceiver.DeployReceiver(ctx, deployer, wasm, salt)
}

// ReportCountForChain reads the report_count from a CRE test receiver directly
// from a CLDF stellar chain. It is exposed for test helpers that do not have a
// full cldf.Environment available.
func ReportCountForChain(ctx context.Context, ch cldfstellar.Chain, contractID string) (uint32, error) {
	deployer, err := stellardeployment.NewDeployerFromChain(ch)
	if err != nil {
		return 0, fmt.Errorf("failed to build stellar deployer: %w", err)
	}

	return stellarreceiver.ReportCount(ctx, deployer, contractID)
}

// DeployRejectingReceiverForChain deploys the CRE rejecting test receiver (always
// returns Err from on_report) directly from a CLDF stellar chain. Used by write
// regression tests to exercise the forwarder's TransmissionState::Failed path.
func DeployRejectingReceiverForChain(ctx context.Context, ch cldfstellar.Chain, salt [32]byte) (string, error) {
	deployer, err := stellardeployment.NewDeployerFromChain(ch)
	if err != nil {
		return "", fmt.Errorf("failed to build stellar deployer: %w", err)
	}

	wasm, err := Artifact(RejectingReceiverWasm)
	if err != nil {
		return "", fmt.Errorf("failed to source rejecting receiver WASM: %w", err)
	}

	contractID, err := deployer.DeployContractBytes(ctx, wasm, salt)
	if err != nil {
		return "", fmt.Errorf("failed to deploy rejecting receiver: %w", err)
	}
	return contractID, nil
}

// ReceiverLastValueU64ForChain reads the last_value_u64 from a CRE test receiver
// directly from a CLDF stellar chain. Used by the write→read roundtrip test to
// assert payload integrity (the first 8 bytes of the report payload as u64).
func ReceiverLastValueU64ForChain(ctx context.Context, ch cldfstellar.Chain, contractID string) (uint64, error) {
	deployer, err := stellardeployment.NewDeployerFromChain(ch)
	if err != nil {
		return 0, fmt.Errorf("failed to build stellar deployer: %w", err)
	}

	sv, err := deployer.SimulateContract(ctx, contractID, "last_value_u64", nil)
	if err != nil {
		return 0, fmt.Errorf("failed to simulate last_value_u64 on %s: %w", contractID, err)
	}
	if sv == nil {
		return 0, fmt.Errorf("nil return from last_value_u64 on %s", contractID)
	}
	u64, ok := sv.GetU64()
	if !ok {
		return 0, fmt.Errorf("last_value_u64 did not return u64 (got scval type %v)", sv.Type)
	}
	return uint64(u64), nil
}
