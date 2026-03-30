package aptos

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	aptoslib "github.com/aptos-labs/aptos-go-sdk"
	aptoscrypto "github.com/aptos-labs/aptos-go-sdk/crypto"
	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

type Deployer struct {
	provider   infra.Provider
	testLogger zerolog.Logger
}

func NewDeployer(testLogger zerolog.Logger, provider *infra.Provider) *Deployer {
	return &Deployer{
		provider:   *provider,
		testLogger: testLogger,
	}
}

type Blockchain struct {
	testLogger    zerolog.Logger
	chainSelector uint64
	chainID       uint64
	ctfOutput     *blockchain.Output
}

func NewBlockchain(testLogger zerolog.Logger, chainID, chainSelector uint64, ctfOutput *blockchain.Output) *Blockchain {
	return &Blockchain{
		testLogger:    testLogger,
		chainSelector: chainSelector,
		chainID:       chainID,
		ctfOutput:     ctfOutput,
	}
}

func (a *Blockchain) ChainSelector() uint64 {
	return a.chainSelector
}

func (a *Blockchain) ChainID() uint64 {
	return a.chainID
}

func (a *Blockchain) CtfOutput() *blockchain.Output {
	return a.ctfOutput
}

func (a *Blockchain) NodeURL() (string, error) {
	if a.ctfOutput == nil || len(a.ctfOutput.Nodes) == 0 {
		return "", fmt.Errorf("no nodes found for Aptos chain %s-%d", a.ChainFamily(), a.chainID)
	}
	return NormalizeNodeURL(a.ctfOutput.Nodes[0].ExternalHTTPUrl)
}

func (a *Blockchain) InternalNodeURL() (string, error) {
	if a.ctfOutput == nil || len(a.ctfOutput.Nodes) == 0 {
		return "", fmt.Errorf("no nodes found for Aptos chain %s-%d", a.ChainFamily(), a.chainID)
	}
	return NormalizeNodeURL(a.ctfOutput.Nodes[0].InternalHTTPUrl)
}

func (a *Blockchain) NodeClient() (*aptoslib.NodeClient, error) {
	nodeURL, err := a.NodeURL()
	if err != nil {
		return nil, err
	}
	chainID, err := aptosChainIDUint8(a.chainID)
	if err != nil {
		return nil, err
	}
	return aptoslib.NewNodeClient(nodeURL, chainID)
}

func (a *Blockchain) LocalDeployerAccount() (*aptoslib.Account, error) {
	var deployerPrivateKey aptoscrypto.Ed25519PrivateKey
	if err := deployerPrivateKey.FromHex(blockchain.DefaultAptosPrivateKey); err != nil {
		return nil, fmt.Errorf("failed to parse default Aptos deployer private key: %w", err)
	}
	deployerAccount, err := aptoslib.NewAccountFromSigner(&deployerPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create default Aptos deployer signer: %w", err)
	}
	return deployerAccount, nil
}

func (a *Blockchain) LocalDeploymentChain() (cldf_aptos.Chain, error) {
	nodeURL, err := a.NodeURL()
	if err != nil {
		return cldf_aptos.Chain{}, err
	}
	client, err := a.NodeClient()
	if err != nil {
		return cldf_aptos.Chain{}, err
	}
	deployerAccount, err := a.LocalDeployerAccount()
	if err != nil {
		return cldf_aptos.Chain{}, err
	}
	return cldf_aptos.Chain{
		Selector:       a.chainSelector,
		Client:         client,
		DeployerSigner: deployerAccount,
		URL:            nodeURL,
		Confirm: func(txHash string, opts ...any) error {
			tx, err := client.WaitForTransaction(txHash, opts...)
			if err != nil {
				return err
			}
			if !tx.Success {
				return fmt.Errorf("transaction failed: %s", tx.VmStatus)
			}
			return nil
		},
	}, nil
}

func (a *Blockchain) IsFamily(chainFamily string) bool {
	return strings.EqualFold(a.ctfOutput.Family, chainFamily)
}

func (a *Blockchain) ChainFamily() string {
	return a.ctfOutput.Family
}

func (a *Blockchain) Fund(ctx context.Context, address string, amount uint64) error {
	client, err := a.NodeClient()
	if err != nil {
		return fmt.Errorf("cannot fund Aptos address %s: create node client: %w", address, err)
	}

	var account aptoslib.AccountAddress
	if parseErr := account.ParseStringRelaxed(address); parseErr != nil {
		return fmt.Errorf("cannot fund Aptos address %q: parse error: %w", address, parseErr)
	}

	faucetURL, err := a.faucetURL()
	if err != nil {
		return fmt.Errorf("failed to derive Aptos faucet URL for %s: %w", address, err)
	}
	faucetClient, err := aptoslib.NewFaucetClient(client, faucetURL)
	if err != nil {
		return fmt.Errorf("failed to create Aptos faucet client for %s: %w", address, err)
	}
	if err := faucetClient.Fund(account, amount); err != nil {
		return fmt.Errorf("failed to fund Aptos address %s via host faucet: %w", address, err)
	}
	if err := waitForAptosAccountVisible(ctx, client, account, 15*time.Second); err != nil {
		return fmt.Errorf("aptos funding request completed but account is still not visible: %w", err)
	}

	a.testLogger.Info().Msgf("Funded Aptos account %s via host faucet (%d octas)", account.StringLong(), amount)
	return nil
}

// ToCldfChain returns the chainlink-deployments-framework aptos.Chain for this blockchain
// so that BlockChains.AptosChains() and saved state work like EVM/Solana.
func (a *Blockchain) ToCldfChain() (cldf_chain.BlockChain, error) {
	nodeURL, err := a.NodeURL()
	if err != nil {
		return nil, fmt.Errorf("invalid Aptos ExternalHTTPUrl for chain %d: %w", a.chainID, err)
	}
	if nodeURL == "" {
		return nil, fmt.Errorf("aptos node has no ExternalHTTPUrl for chain %d", a.chainID)
	}
	client, err := a.NodeClient()
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "create Aptos RPC client for chain %d", a.chainID)
	}
	return cldf_aptos.Chain{
		Selector:       a.chainSelector,
		Client:         client,
		DeployerSigner: nil, // CRE read-only use; deployer not required for View calls
		URL:            nodeURL,
		Confirm: func(txHash string, opts ...any) error {
			tx, err := client.WaitForTransaction(txHash, opts...)
			if err != nil {
				return err
			}
			if !tx.Success {
				return fmt.Errorf("transaction failed: %s", tx.VmStatus)
			}
			return nil
		},
	}, nil
}

func (a *Deployer) Deploy(ctx context.Context, input *blockchain.Input) (blockchains.Blockchain, error) {
	var bcOut *blockchain.Output
	var err error

	switch {
	case a.provider.IsKubernetes():
		if err = blockchains.ValidateKubernetesBlockchainOutput(input); err != nil {
			return nil, err
		}
		a.testLogger.Info().Msgf("Using configured Kubernetes blockchain URLs for %s (chain_id: %s)", input.Type, input.ChainID)
		bcOut = input.Out
	case input.Out != nil:
		bcOut = input.Out
	default:
		bcOut, err = blockchain.NewWithContext(ctx, input)
		if err != nil {
			return nil, pkgerrors.Wrapf(err, "failed to deploy blockchain %s chainID: %s", input.Type, input.ChainID)
		}
	}

	// Framework Aptos output may have empty ChainID; use config input.ChainID (e.g. "4" for local devnet)
	chainIDStr := bcOut.ChainID
	if chainIDStr == "" {
		chainIDStr = input.ChainID
	}
	if chainIDStr == "" {
		return nil, pkgerrors.New("aptos chain id is empty (set chain_id in [[blockchains]] in TOML)")
	}
	chainID, err := strconv.ParseUint(chainIDStr, 10, 64)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "failed to parse chain id %s", chainIDStr)
	}

	selector, err := aptosChainSelector(chainIDStr, chainID)
	if err != nil {
		return nil, err
	}

	// Ensure ctfOutput has ChainID set for downstream (e.g. findAptosChains)
	bcOut.ChainID = chainIDStr

	return NewBlockchain(a.testLogger, chainID, selector, bcOut), nil
}

// aptosChainSelector returns the chain selector for the given Aptos chain ID.
// Uses chain-selectors when available; falls back to known Aptos localnet selector for chain_id 4.
func aptosChainSelector(chainIDStr string, chainID uint64) (uint64, error) {
	chainDetails, err := chainselectors.GetChainDetailsByChainIDAndFamily(chainIDStr, chainselectors.FamilyAptos)
	if err == nil {
		return chainDetails.ChainSelector, nil
	}
	// Fallback: Aptos local devnet (aptos node run-local-testnet) uses chain_id 4 and this selector
	if chainID == 4 {
		const aptosLocalnetSelector = 4457093679053095497
		return aptosLocalnetSelector, nil
	}
	return 0, pkgerrors.Wrapf(err, "failed to get chain selector for Aptos chain id %s", chainIDStr)
}

func aptosNodeURLWithV1(rawURL string) (string, error) {
	u, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return "", err
	}
	if u.Scheme == "" || u.Host == "" {
		return "", fmt.Errorf("invalid url %q", rawURL)
	}
	path := strings.TrimRight(u.Path, "/")
	if path == "" || path != "/v1" {
		u.Path = "/v1"
	}
	u.RawPath = ""
	u.RawQuery = ""
	u.Fragment = ""
	return u.String(), nil
}

func NormalizeNodeURL(rawURL string) (string, error) {
	return aptosNodeURLWithV1(rawURL)
}

func aptosFaucetURLFromNodeURL(nodeURL string) (string, error) {
	u, err := url.Parse(nodeURL)
	if err != nil {
		return "", err
	}
	host := u.Hostname()
	if host == "" {
		return "", fmt.Errorf("empty host in node url %q", nodeURL)
	}
	u.Host = host + ":8081"
	u.Path = ""
	u.RawPath = ""
	u.RawQuery = ""
	u.Fragment = ""
	return u.String(), nil
}

func FaucetURLFromNodeURL(nodeURL string) (string, error) {
	return aptosFaucetURLFromNodeURL(nodeURL)
}

func (a *Blockchain) faucetURL() (string, error) {
	if a.ctfOutput == nil || len(a.ctfOutput.Nodes) == 0 {
		return "", errors.New("missing chain nodes output")
	}
	nodeURL, err := NormalizeNodeURL(a.ctfOutput.Nodes[0].ExternalHTTPUrl)
	if err != nil {
		return "", err
	}
	return FaucetURLFromNodeURL(nodeURL)
}

func waitForAptosAccountVisible(ctx context.Context, client *aptoslib.NodeClient, account aptoslib.AccountAddress, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	var lastErr error
	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		_, accountErr := client.Account(account)
		if accountErr == nil {
			return nil
		}
		lastErr = accountErr
		time.Sleep(1 * time.Second)
	}
	if lastErr != nil {
		return fmt.Errorf("account %s not visible after funding attempt: %w", account.StringLong(), lastErr)
	}
	return fmt.Errorf("account %s not visible after funding attempt", account.StringLong())
}

func aptosChainIDUint8(chainID uint64) (uint8, error) {
	if chainID > uint64(^uint8(0)) {
		return 0, fmt.Errorf("aptos chain id %d does not fit in uint8", chainID)
	}

	return uint8(chainID), nil
}

func ChainIDUint8(chainID uint64) (uint8, error) {
	return aptosChainIDUint8(chainID)
}

func WaitForTransactionSuccess(client *aptoslib.NodeClient, txHash, label string) error {
	tx, err := client.WaitForTransaction(txHash)
	if err != nil {
		return fmt.Errorf("failed waiting for Aptos tx %s: %w", label, err)
	}
	if !tx.Success {
		return fmt.Errorf("aptos tx failed: %s vm_status=%s", label, tx.VmStatus)
	}
	return nil
}
