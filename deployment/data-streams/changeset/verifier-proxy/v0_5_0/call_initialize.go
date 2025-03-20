package v0_5_0

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier_proxy_v0_5_0"
)

var InitializeChangeset = deployment.CreateChangeSet(initializeLogic, initializePrecondition)

type InitializeConfig struct {
	ConfigsByChain map[uint64][]Initialize
}

type Initialize struct {
	VerifierAddress      common.Address
	VerifierProxyAddress common.Address
}

type VerifierState struct {
	VerifierProxy *verifier_proxy_v0_5_0.VerifierProxy
}

func (cfg InitializeConfig) Validate() error {
	if len(cfg.ConfigsByChain) == 0 {
		return errors.New("ConfigsByChain cannot be empty")
	}
	return nil
}

func initializePrecondition(_ deployment.Environment, cc InitializeConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid Initialize config: %w", err)
	}
	return nil
}

func initializeLogic(e deployment.Environment, cfg InitializeConfig) (deployment.ChangesetOutput, error) {
	for chainSelector, configs := range cfg.ConfigsByChain {
		chain, chainExists := e.Chains[chainSelector]
		if !chainExists {
			return deployment.ChangesetOutput{}, fmt.Errorf("chain %d not found in environment", chainSelector)
		}

		opts := getTransactOptsInitialize(e, chainSelector)
		for _, config := range configs {
			confState, err := maybeLoadVerifierProxy(e, chainSelector, config.VerifierProxyAddress.String())
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}

			_, err = initializeOrBuildTx(e, confState.VerifierProxy, config, opts, chain)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
		}
	}

	return deployment.ChangesetOutput{}, nil
}

func initializeOrBuildTx(
	e deployment.Environment,
	verifierProxyContract *verifier_proxy_v0_5_0.VerifierProxy,
	cfg Initialize,
	opts *bind.TransactOpts,
	chain deployment.Chain,
) (*ethTypes.Transaction, error) {
	proxyAddr := verifierProxyContract.Address()
	parsedABI, err := abi.JSON(strings.NewReader(verifier_proxy_v0_5_0.VerifierProxyABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	data, err := parsedABI.Pack("initializeVerifier", cfg.VerifierAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to pack call data: %w", err)
	}

	callMsg := ethereum.CallMsg{From: opts.From, To: &proxyAddr, Data: data}

	ctx := context.Background()

	result, callErr := chain.Client.CallContract(ctx, callMsg, nil)

	if callErr != nil {
		if revertData, found := ethclient.RevertErrorData(callErr); found {
			fmt.Println("Revert reason (raw hex):", hexutil.Encode(revertData))
		} else {
			fmt.Println("CallContract error without revert data:", callErr)
		}
	} else {
		fmt.Println("CallContract simulation succeeded, result:", result)
	}

	tx, err := verifierProxyContract.InitializeVerifier(
		opts,
		cfg.VerifierAddress,
	)
	if err != nil {
		return nil, fmt.Errorf("error packing initialize tx data: %w", err)
	}

	if _, err := deployment.ConfirmIfNoError(chain, tx, err); err != nil {
		e.Logger.Errorw("Failed to confirm initialize tx", "chain", chain.String(), "err", err)
		return nil, err
	}

	return tx, nil
}

func maybeLoadVerifierProxy(e deployment.Environment, chainSel uint64, contractAddr string) (*VerifierState, error) {
	chain, ok := e.Chains[chainSel]
	if !ok {
		return nil, fmt.Errorf("chain %d not found", chainSel)
	}

	addresses, err := e.ExistingAddresses.AddressesForChain(chainSel)
	if err != nil {
		return nil, err
	}

	tv, found := addresses[contractAddr]
	if !found {
		return nil, fmt.Errorf("unable to find VerifierProxy contract on chain %s (selector %d)", chain.Name(), chain.Selector)
	}
	if tv.Type != types.VerifierProxy || tv.Version != deployment.Version0_5_0 {
		return nil, fmt.Errorf("unexpected contract type %s for Verifier on chain %s (selector %d)", tv, chain.Name(), chain.Selector)
	}

	conf, err := verifier_proxy_v0_5_0.NewVerifierProxy(common.HexToAddress(contractAddr), chain.Client)
	if err != nil {
		return nil, err
	}

	return &VerifierState{VerifierProxy: conf}, nil
}

func getTransactOptsInitialize(e deployment.Environment, chainSel uint64) *bind.TransactOpts {
	return e.Chains[chainSel].DeployerKey
}
