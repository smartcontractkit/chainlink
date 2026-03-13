package aptos

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	aptoslib "github.com/aptos-labs/aptos-go-sdk"
	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
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

func (a *Blockchain) ChainSelector() uint64 {
	return a.chainSelector
}

func (a *Blockchain) ChainID() uint64 {
	return a.chainID
}

func (a *Blockchain) CtfOutput() *blockchain.Output {
	return a.ctfOutput
}

func (a *Blockchain) IsFamily(chainFamily string) bool {
	return strings.EqualFold(a.ctfOutput.Family, chainFamily)
}

func (a *Blockchain) ChainFamily() string {
	return a.ctfOutput.Family
}

func (a *Blockchain) Fund(ctx context.Context, address string, amount uint64) error {
	if a.ctfOutput == nil || len(a.ctfOutput.Nodes) == 0 {
		return fmt.Errorf("cannot fund Aptos address %s: missing chain nodes output", address)
	}

	nodeURL, err := aptosNodeURLWithV1(a.ctfOutput.Nodes[0].ExternalHTTPUrl)
	if err != nil {
		return fmt.Errorf("cannot fund Aptos address %s: invalid node URL: %w", address, err)
	}

	client, err := aptoslib.NewNodeClient(nodeURL, uint8(a.chainID))
	if err != nil {
		return fmt.Errorf("cannot fund Aptos address %s: create node client: %w", address, err)
	}

	var account aptoslib.AccountAddress
	if err := account.ParseStringRelaxed(address); err != nil {
		return fmt.Errorf("cannot fund Aptos address %q: parse error: %w", address, err)
	}

	// Fast path: already visible, still top up because fees are variable.
	faucetURL, ferr := aptosFaucetURLFromNodeURL(nodeURL)
	if ferr == nil {
		if faucetClient, cErr := aptoslib.NewFaucetClient(client, faucetURL); cErr == nil {
			if fundErr := faucetClient.Fund(account, amount); fundErr == nil {
				if waitErr := waitForAptosAccountVisible(ctx, client, account, 15*time.Second); waitErr == nil {
					a.testLogger.Info().Msgf("Funded Aptos account %s via host faucet (%d octas)", account.StringLong(), amount)
					return nil
				}
			}
		}
	}

	containerName := strings.TrimSpace(a.ctfOutput.ContainerName)
	if containerName == "" {
		return fmt.Errorf("failed to fund Aptos address %s via host faucet and no container fallback available", address)
	}

	dc, err := framework.NewDockerClient()
	if err != nil {
		return fmt.Errorf("failed to create docker client for Aptos funding fallback: %w", err)
	}
	_, execErr := dc.ExecContainerWithContext(ctx, containerName, []string{
		"aptos", "account", "fund-with-faucet",
		"--account", account.StringLong(),
		"--amount", strconv.FormatUint(amount, 10),
	})
	if execErr != nil {
		return fmt.Errorf("failed to fund Aptos address %s via container faucet fallback: %w", address, execErr)
	}
	if waitErr := waitForAptosAccountVisible(ctx, client, account, 20*time.Second); waitErr != nil {
		return fmt.Errorf("Aptos funding fallback completed but account still not visible: %w", waitErr)
	}

	a.testLogger.Info().Msgf("Funded Aptos account %s via container faucet fallback (%d octas)", account.StringLong(), amount)
	return nil
}

// ToCldfChain returns the chainlink-deployments-framework aptos.Chain for this blockchain
// so that BlockChains.AptosChains() and saved state work like EVM/Solana.
func (a *Blockchain) ToCldfChain() (cldf_chain.BlockChain, error) {
	if a.ctfOutput == nil || len(a.ctfOutput.Nodes) == 0 {
		return nil, fmt.Errorf("no nodes found for Aptos chain %s-%d", a.ChainFamily(), a.chainID)
	}
	url := a.ctfOutput.Nodes[0].ExternalHTTPUrl
	if url == "" {
		return nil, fmt.Errorf("Aptos node has no ExternalHTTPUrl for chain %d", a.chainID)
	}
	// Aptos chain IDs are small (e.g. 1=mainnet, 2=testnet, 4=local devnet).
	client, err := aptoslib.NewNodeClient(url, uint8(a.chainID))
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "create Aptos RPC client for chain %d", a.chainID)
	}
	return cldf_aptos.Chain{
		Selector:       a.chainSelector,
		Client:         client,
		DeployerSigner: nil, // CRE read-only use; deployer not required for View calls
		URL:            url,
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

	return &Blockchain{
		testLogger:    a.testLogger,
		chainSelector: selector,
		chainID:       chainID,
		ctfOutput:     bcOut,
	}, nil
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

func aptosFaucetURLFromNodeURL(nodeURL string) (string, error) {
	u, err := url.Parse(nodeURL)
	if err != nil {
		return "", err
	}
	host := u.Hostname()
	if host == "" {
		return "", fmt.Errorf("empty host in node url %q", nodeURL)
	}
	u.Host = fmt.Sprintf("%s:8081", host)
	u.Path = ""
	u.RawPath = ""
	u.RawQuery = ""
	u.Fragment = ""
	return u.String(), nil
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
		if _, err := client.Account(account); err == nil {
			return nil
		} else {
			lastErr = err
		}
		time.Sleep(1 * time.Second)
	}
	if lastErr != nil {
		return fmt.Errorf("account %s not visible after funding attempt: %w", account.StringLong(), lastErr)
	}
	return fmt.Errorf("account %s not visible after funding attempt", account.StringLong())
}
