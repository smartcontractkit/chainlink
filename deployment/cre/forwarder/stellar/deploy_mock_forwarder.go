package stellar

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	stellardeployment "github.com/smartcontractkit/chainlink-stellar/deployment"
	stellarforwarder "github.com/smartcontractkit/chainlink-stellar/deployment/cre/forwarder"

	crestellar "github.com/smartcontractkit/chainlink/deployment/cre/stellar"
)

const MockForwarderContract datastore.ContractType = "StellarMockForwarder"

const DefaultMockForwarderVersion = "1.0.0"

var _ cldf.ChangeSetV2[*DeployMockForwarderRequest] = DeployMockForwarder{}

// DeployMockForwarder deploys the CRE mock forwarder
// (chainlink-stellar/contracts/cre/mock_forwarder). Never point a production
// receiver at it.
type DeployMockForwarder struct{}

type DeployMockForwarderRequest struct {
	ChainSel  uint64
	Qualifier string
	Version   string
	LabelSet  datastore.LabelSet
	Salt      [32]byte
}

func (cs DeployMockForwarder) VerifyPreconditions(env cldf.Environment, req *DeployMockForwarderRequest) error {
	if req == nil {
		return errors.New("request is required")
	}
	if _, ok := env.BlockChains.StellarChains()[req.ChainSel]; !ok {
		return fmt.Errorf("stellar chain not found for chain selector %d", req.ChainSel)
	}
	if req.Qualifier == "" {
		return errors.New("mock forwarder qualifier is required")
	}
	if req.Version == "" {
		return errors.New("mock forwarder version is required")
	}
	if _, err := semver.NewVersion(req.Version); err != nil {
		return fmt.Errorf("invalid mock forwarder version %q: %w", req.Version, err)
	}
	return nil
}

func (cs DeployMockForwarder) Apply(env cldf.Environment, req *DeployMockForwarderRequest) (cldf.ChangesetOutput, error) {
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
	deployerAddr := ch.Signer.Address()

	// Namespace a zero salt by contract type so the mock never lands on the
	// address of another contract the same deployer creates with a zero salt.
	salt := req.Salt
	if salt == ([32]byte{}) {
		salt = stellardeployment.GenerateDeterministicSalt(deployerAddr, string(MockForwarderContract))
	}

	deployer, err := stellardeployment.NewDeployerFromChain(ch)
	if err != nil {
		return out, fmt.Errorf("failed to build stellar deployer for chain selector %d: %w", req.ChainSel, err)
	}

	wasm, err := crestellar.Artifact(crestellar.MockForwarderWasm)
	if err != nil {
		return out, fmt.Errorf("failed to source mock forwarder WASM: %w", err)
	}

	mockAddr, err := stellarforwarder.DeployMockForwarder(env.GetContext(), deployer, wasm, salt)
	if err != nil {
		if mockAddr == "" {
			return out, fmt.Errorf("failed to deploy stellar mock forwarder on chain selector %d: %w", req.ChainSel, err)
		}
		env.Logger.Warnw("Deployed Stellar CRE mock forwarder but failed to extend its TTL; extend it with the stellar_extend_contract_ttl changeset",
			"chainSelector", req.ChainSel, "mockForwarder", mockAddr, "err", err)
	}

	if err := out.DataStore.Addresses().Add(datastore.AddressRef{
		Address:       mockAddr,
		ChainSelector: req.ChainSel,
		Type:          MockForwarderContract,
		Version:       version,
		Qualifier:     req.Qualifier,
		Labels:        req.LabelSet,
	}); err != nil && !errors.Is(err, datastore.ErrAddressRefExists) {
		return out, fmt.Errorf("failed to add stellar mock forwarder address to datastore: %w", err)
	}

	env.Logger.Infow("Deployed Stellar CRE mock forwarder", "chainSelector", req.ChainSel, "mockForwarder", mockAddr, "qualifier", req.Qualifier)

	return out, nil
}
