package v1_6

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	ccipseq "github.com/smartcontractkit/chainlink/deployment/ccip/sequence/evm/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/ccip_home"
)

var _ cldf.ChangeSet[ccipseq.DeployChainContractsConfig] = DeployChainContractsChangeset

// DeployChainContracts deploys all new CCIP v1.6 or later contracts for the given chains.
// It returns the new addresses for the contracts.
// DeployChainContractsChangeset is idempotent. If there is an error, it will return the successfully deployed addresses and the error so that the caller can call the
// changeset again with the same input to retry the failed deployment.
// Caller should update the environment's address book with the returned addresses.
// Points to note :
// In case of migrating from legacy ccip to 1.6, the previous RMN address should be set while deploying RMNRemote.
// if there is no existing RMN address found, RMNRemote will be deployed with 0x0 address for previous RMN address
// which will set RMN to 0x0 address immutably in RMNRemote.
func DeployChainContractsChangeset(env cldf.Environment, c ccipseq.DeployChainContractsConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid DeployChainContractsConfig: %w", err)
	}
	report, err := deployChainContractsForChains(env, c.HomeChainSelector, c)
	if err != nil {
		return cldf.ChangesetOutput{
			Reports: report.ExecutionReports,
		}, fmt.Errorf("failed to deploy CCIP contracts: %w", err)
	}
	addressBook := cldf.NewMemoryAddressBook()
	for chainSel, addresses := range report.Output {
		for address, typeAndVersion := range addresses {
			err := addressBook.Save(chainSel, address, cldf.MustTypeAndVersionFromString(typeAndVersion))
			if err != nil {
				return cldf.ChangesetOutput{
					Reports: report.ExecutionReports,
				}, fmt.Errorf("failed to save address %s for chain %d: %w", address, chainSel, err)
			}
		}
	}
	return cldf.ChangesetOutput{
		Reports:     report.ExecutionReports,
		AddressBook: addressBook,
	}, nil
}

func ValidateHomeChainState(e cldf.Environment, homeChainSel uint64, existingState stateview.CCIPOnChainState) error {
	capReg := existingState.Chains[homeChainSel].CapabilityRegistry
	if capReg == nil {
		e.Logger.Errorw("Failed to get capability registry")
		return errors.New("capability registry not found")
	}
	cr, err := capReg.GetHashedCapabilityId(
		&bind.CallOpts{}, shared.CapabilityLabelledName, shared.CapabilityVersion)
	if err != nil {
		e.Logger.Errorw("Failed to get hashed capability id", "err", err)
		return err
	}
	if cr != shared.CCIPCapabilityID {
		return fmt.Errorf("unexpected mismatch between calculated ccip capability id (%s) and expected ccip capability id constant (%s)",
			hexutil.Encode(cr[:]),
			hexutil.Encode(shared.CCIPCapabilityID[:]))
	}
	capability, err := capReg.GetCapability(nil, shared.CCIPCapabilityID)
	if err != nil {
		e.Logger.Errorw("Failed to get capability", "err", err)
		return err
	}
	ccipHome, err := ccip_home.NewCCIPHome(capability.ConfigurationContract, e.BlockChains.EVMChains()[homeChainSel].Client)
	if err != nil {
		e.Logger.Errorw("Failed to get ccip config", "err", err)
		return err
	}
	if ccipHome.Address() != existingState.Chains[homeChainSel].CCIPHome.Address() {
		return errors.New("ccip home address mismatch")
	}
	rmnHome := existingState.Chains[homeChainSel].RMNHome
	if rmnHome == nil {
		e.Logger.Errorw("Failed to get rmn home", "err", err)
		return errors.New("rmn home not found")
	}
	return nil
}

func deployChainContractsForChains(
	e cldf.Environment,
	homeChainSel uint64,
	c ccipseq.DeployChainContractsConfig,
) (operations.SequenceReport[ccipseq.DeployChainContractsSeqConfig, map[uint64]map[string]string], error) {
	existingState, err := stateview.LoadOnchainState(e)
	if err != nil {
		e.Logger.Errorw("Failed to load existing onchain state", "err", err)
		return operations.SequenceReport[ccipseq.DeployChainContractsSeqConfig, map[uint64]map[string]string]{}, err
	}

	err = ValidateHomeChainState(e, homeChainSel, existingState)
	if err != nil {
		return operations.SequenceReport[ccipseq.DeployChainContractsSeqConfig, map[uint64]map[string]string]{}, err
	}

	addresses := make(map[uint64]ccipseq.CCIPAddresses)
	for chainSel, params := range c.ContractParamsPerChain {
		linkToken, err := existingState.Chains[chainSel].LinkTokenAddress()
		if err != nil {
			return operations.SequenceReport[ccipseq.DeployChainContractsSeqConfig, map[uint64]map[string]string]{}, err
		}

		fq := existingState.Chains[chainSel].FeeQuoter
		fqVersion := existingState.Chains[chainSel].FeeQuoterVersion
		if fq != nil && fqVersion != nil &&
			params.FeeQuoterOpts != nil &&
			params.FeeQuoterOpts.Version != nil &&
			params.FeeQuoterOpts.Version.GreaterThan(fqVersion) {
			fq = nil // Deploy a new FeeQuoter if the version in params is greater than existing version
		}

		addresses[chainSel] = ccipseq.CCIPAddresses{
			LegacyRMNAddress:          getAddressSafely(existingState.Chains[chainSel].RMN),
			RMNProxyAddress:           getAddressSafely(existingState.Chains[chainSel].RMNProxy),
			WrappedNativeAddress:      getAddressSafely(existingState.Chains[chainSel].Weth9),
			TimelockAddress:           getAddressSafely(existingState.Chains[chainSel].Timelock),
			LinkAddress:               linkToken,
			FeeAggregatorAddress:      existingState.Chains[chainSel].FeeAggregator,
			TokenAdminRegistryAddress: getAddressSafely(existingState.Chains[chainSel].TokenAdminRegistry),
			OnRampAddress:             getAddressSafely(existingState.Chains[chainSel].OnRamp),
			TestRouterAddress:         getAddressSafely(existingState.Chains[chainSel].TestRouter),
			OffRampAddress:            getAddressSafely(existingState.Chains[chainSel].OffRamp),
			NonceManagerAddress:       getAddressSafely(existingState.Chains[chainSel].NonceManager),
			FeeQuoterAddress:          getAddressSafely(fq),
			RMNRemoteAddress:          getAddressSafely(existingState.Chains[chainSel].RMNRemote),
		}
	}

	report, err := operations.ExecuteSequence(
		e.OperationsBundle,
		ccipseq.DeployChainContractsSeq,
		e.BlockChains.EVMChains(),
		ccipseq.DeployChainContractsSeqConfig{
			RMNHomeAddress:             getAddressSafely(existingState.Chains[homeChainSel].RMNHome),
			DeployChainContractsConfig: c,
			AddressesPerChain:          addresses,
			GasBoostConfigPerChain:     c.GasBoostConfigPerChain,
		},
	)
	if err != nil {
		return report, fmt.Errorf("failed to deploy chain contracts: %w", err)
	}

	return report, nil
}

type addressable interface {
	Address() common.Address
}

func getAddressSafely(a addressable) common.Address {
	if a == nil || reflect.ValueOf(a).IsNil() { // assumes 'a' is a pointer type
		return common.Address{}
	}
	return a.Address()
}
