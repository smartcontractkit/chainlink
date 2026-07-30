package stellar

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	stellardeployment "github.com/smartcontractkit/chainlink-stellar/deployment"
	stellarforwarder "github.com/smartcontractkit/chainlink-stellar/deployment/cre/forwarder"

	"github.com/smartcontractkit/chainlink/deployment/helpers"
)

var _ cldf.ChangeSetV2[*DeployForwarderRequest] = DeployForwarder{}

type DeployForwarder struct{}

type DeployForwarderRequest struct {
	ChainSel    uint64
	Qualifier   string
	Version     string
	LabelSet    datastore.LabelSet
	Salt        [32]byte
	BuildConfig *helpers.BuildStellarConfig
}

func (cs DeployForwarder) VerifyPreconditions(env cldf.Environment, req *DeployForwarderRequest) error {
	if req == nil {
		return errors.New("request is required")
	}
	if _, ok := env.BlockChains.StellarChains()[req.ChainSel]; !ok {
		return fmt.Errorf("stellar chain not found for chain selector %d", req.ChainSel)
	}
	if req.Qualifier == "" {
		return errors.New("forwarder qualifier is required")
	}
	if req.Version == "" {
		return errors.New("forwarder version is required")
	}
	if _, err := semver.NewVersion(req.Version); err != nil {
		return fmt.Errorf("invalid forwarder version %q: %w", req.Version, err)
	}
	if req.BuildConfig == nil {
		return errors.New("build config is required to source the forwarder WASM")
	}
	return nil
}

func (cs DeployForwarder) Apply(env cldf.Environment, req *DeployForwarderRequest) (cldf.ChangesetOutput, error) {
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
	owner := ch.Signer.Address()

	// Soroban requires a salt for create_contract (address = hash(networkID ‖
	// deployer ‖ salt)). Default a zero salt to a value namespaced by the contract
	// type so this deploy does not collide with other contracts the same deployer
	// instantiates with a zero salt.
	salt := req.Salt
	if salt == ([32]byte{}) {
		salt = stellardeployment.GenerateDeterministicSalt(owner, string(ForwarderContract))
	}

	deployer, err := stellardeployment.NewDeployerFromChain(ch)
	if err != nil {
		return out, fmt.Errorf("failed to build stellar deployer for chain selector %d: %w", req.ChainSel, err)
	}

	wasm, err := helpers.BuildStellar(env.GetContext(), *req.BuildConfig)
	if err != nil {
		return out, fmt.Errorf("failed to source forwarder WASM: %w", err)
	}
	forwarderAddr, err := stellarforwarder.DeployForwarder(env.GetContext(), deployer, owner, wasm, salt)
	if err != nil {
		return out, fmt.Errorf("failed to deploy stellar forwarder on chain selector %d: %w", req.ChainSel, err)
	}

	if err := out.DataStore.Addresses().Add(datastore.AddressRef{
		Address:       forwarderAddr,
		ChainSelector: req.ChainSel,
		Type:          ForwarderContract,
		Version:       version,
		Qualifier:     req.Qualifier,
		Labels:        req.LabelSet,
	}); err != nil && !errors.Is(err, datastore.ErrAddressRefExists) {
		return out, fmt.Errorf("failed to add stellar forwarder address to datastore: %w", err)
	}

	env.Logger.Infow("Deployed Stellar CRE forwarder", "chainSelector", req.ChainSel, "forwarder", forwarderAddr, "owner", owner, "qualifier", req.Qualifier)

	return out, nil
}
