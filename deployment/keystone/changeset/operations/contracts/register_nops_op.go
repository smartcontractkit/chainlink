package contracts

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"golang.org/x/exp/maps"

	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

type RegisterNopsOpDeps struct {
	Env           *cldf.Environment
	RegistryChain *evm.Chain
	Contract      *capabilities_registry.CapabilitiesRegistry
	NopsToNodes   map[capabilities_registry.CapabilitiesRegistryNodeOperator][]string
}

type RegisterNopsOpInput struct {
	UseMCMS          bool
	RegistryChainSel uint64
}

type RegisterNopsOpOutput struct {
	Nops           []*capabilities_registry.CapabilitiesRegistryNodeOperatorAdded // if UseMCMS is false, a list of added node operators is returned
	BatchOperation *mcmstypes.BatchOperation                                      // if UseMCMS is true, a batch proposal is returned and no transaction is confirmed onchain.
}

var RegisterNopsOp = operations.NewOperation[RegisterNopsOpInput, RegisterNopsOpOutput, RegisterNopsOpDeps](
	"register-nops-op",
	semver.MustParse("1.0.0"),
	"Register Node Operators in Capabilities Registry",
	func(b operations.Bundle, deps RegisterNopsOpDeps, input RegisterNopsOpInput) (RegisterNopsOpOutput, error) {
		nopsList := maps.Keys(deps.NopsToNodes)
		nopsResp, err := internal.RegisterNOPS(b.GetContext(), b.Logger, internal.RegisterNOPSRequest{
			Env:                   deps.Env,
			RegistryChainSelector: input.RegistryChainSel,
			Nops:                  nopsList,
			UseMCMS:               input.UseMCMS,
			Registry:              deps.Contract,
			RegistryChain:         deps.RegistryChain,
		})
		if err != nil {
			return RegisterNopsOpOutput{}, fmt.Errorf("register-nops-op failed: %w", err)
		}
		b.Logger.Infow("registered node operators", "nops", nopsResp.Nops)

		return RegisterNopsOpOutput{
			Nops:           nopsResp.Nops,
			BatchOperation: nopsResp.Ops,
		}, nil
	},
)
