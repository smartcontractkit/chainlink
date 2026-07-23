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

	"github.com/smartcontractkit/chainlink/deployment/helpers"
)

const ReadFixtureContract datastore.ContractType = "StellarReadFixture"

var _ cldf.ChangeSetV2[*DeployReadFixtureRequest] = DeployReadFixture{}

type DeployReadFixture struct{}

type DeployReadFixtureRequest struct {
	ChainSel    uint64
	Qualifier   string
	Version     string
	LabelSet    datastore.LabelSet
	Salt        [32]byte
	BuildConfig *helpers.BuildStellarConfig
}

func (cs DeployReadFixture) VerifyPreconditions(env cldf.Environment, req *DeployReadFixtureRequest) error {
	if req == nil {
		return errors.New("request is required")
	}
	if _, ok := env.BlockChains.StellarChains()[req.ChainSel]; !ok {
		return fmt.Errorf("stellar chain not found for chain selector %d", req.ChainSel)
	}
	if req.Qualifier == "" {
		return errors.New("read fixture qualifier is required")
	}
	if req.Version == "" {
		return errors.New("read fixture version is required")
	}
	if _, err := semver.NewVersion(req.Version); err != nil {
		return fmt.Errorf("invalid read fixture version %q: %w", req.Version, err)
	}
	if req.BuildConfig == nil {
		return errors.New("build config is required to source the read fixture WASM")
	}
	return nil
}

func (cs DeployReadFixture) Apply(env cldf.Environment, req *DeployReadFixtureRequest) (cldf.ChangesetOutput, error) {
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

	fixtureAddr, err := DeployReadFixtureForChain(env.GetContext(), ch, *req.BuildConfig, req.Salt)
	if err != nil {
		return out, fmt.Errorf("failed to deploy stellar read fixture on chain selector %d: %w", req.ChainSel, err)
	}

	if err := out.DataStore.Addresses().Add(datastore.AddressRef{
		Address:       fixtureAddr,
		ChainSelector: req.ChainSel,
		Type:          ReadFixtureContract,
		Version:       version,
		Qualifier:     req.Qualifier,
		Labels:        req.LabelSet,
	}); err != nil && !errors.Is(err, datastore.ErrAddressRefExists) {
		return out, fmt.Errorf("failed to add stellar read fixture address to datastore: %w", err)
	}

	env.Logger.Infow("Deployed Stellar CRE read fixture", "chainSelector", req.ChainSel, "fixture", fixtureAddr, "qualifier", req.Qualifier)

	return out, nil
}

// DeployReadFixtureForChain deploys the CRE read fixture contract directly from
// a CLDF stellar chain. It is exposed for test helpers that do not have a full
// cldf.Environment available.
func DeployReadFixtureForChain(ctx context.Context, ch cldfstellar.Chain, buildCfg helpers.BuildStellarConfig, salt [32]byte) (string, error) {
	deployer, err := stellardeployment.NewDeployerFromChain(ch)
	if err != nil {
		return "", fmt.Errorf("failed to build stellar deployer: %w", err)
	}

	wasm, err := helpers.BuildStellar(ctx, buildCfg)
	if err != nil {
		return "", fmt.Errorf("failed to source read fixture WASM: %w", err)
	}

	return deployer.DeployContractBytes(ctx, wasm, salt)
}
