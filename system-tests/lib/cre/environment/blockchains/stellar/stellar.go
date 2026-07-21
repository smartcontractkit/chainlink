package stellar

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_stellar_provider "github.com/smartcontractkit/chainlink-deployments-framework/chain/stellar/provider"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

// friendbot funding retry bounds, matching the CLDF stellar reference helper
// (require.Eventually 60s / 1s): the RPC + friendbot can lag container readiness.
const (
	friendbotFundTimeout  = 60 * time.Second
	friendbotFundInterval = 1 * time.Second
)

// Deployer starts (Docker) or reuses (K8s) a local Stellar/Soroban node via the
// Chainlink Testing Framework (image stellar/quickstart, RPC on /rpc, Friendbot
// faucet on /friendbot). Mirrors the Solana/Aptos deployers.
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
	testLogger        zerolog.Logger
	chainSelector     uint64
	stellarChainID    string // Stellar network id (SHA-256 of the passphrase), a string
	friendbotURL      string
	networkPassphrase string
	ctfOutput         *blockchain.Output
}

func NewBlockchain(testLogger zerolog.Logger, stellarChainID string, chainSelector uint64, friendbotURL, networkPassphrase string, ctfOutput *blockchain.Output) *Blockchain {
	return &Blockchain{
		testLogger:        testLogger,
		chainSelector:     chainSelector,
		stellarChainID:    stellarChainID,
		friendbotURL:      friendbotURL,
		networkPassphrase: networkPassphrase,
		ctfOutput:         ctfOutput,
	}
}

func (s *Blockchain) ChainSelector() uint64 { return s.chainSelector }

// ChainID returns 0: Stellar identifies networks by a string passphrase/network
// id, not a numeric chain id (same convention as the Solana blockchain layer).
// Use StellarChainID() for the string network id.
func (s *Blockchain) ChainID() uint64 { return 0 }

// StellarChainID returns the Stellar network id string used to scope the
// capability (e.g. the STELLAR_LOCALNET network id).
func (s *Blockchain) StellarChainID() string { return s.stellarChainID }

func (s *Blockchain) CtfOutput() *blockchain.Output { return s.ctfOutput }

func (s *Blockchain) ChainFamily() string { return s.ctfOutput.Family }

func (s *Blockchain) IsFamily(chainFamily string) bool {
	return strings.EqualFold(s.ctfOutput.Family, chainFamily)
}

// NodeURL returns the external Soroban RPC URL (…/rpc), for callers on the host.
func (s *Blockchain) NodeURL() (string, error) {
	if s.ctfOutput == nil || len(s.ctfOutput.Nodes) == 0 {
		return "", fmt.Errorf("no nodes found for Stellar chain %s-%s", s.ChainFamily(), s.stellarChainID)
	}
	return s.ctfOutput.Nodes[0].ExternalHTTPUrl, nil
}

// InternalNodeURL returns the in-network Soroban RPC URL (…/rpc), for node TOML.
func (s *Blockchain) InternalNodeURL() (string, error) {
	if s.ctfOutput == nil || len(s.ctfOutput.Nodes) == 0 {
		return "", fmt.Errorf("no nodes found for Stellar chain %s-%s", s.ChainFamily(), s.stellarChainID)
	}
	return s.ctfOutput.Nodes[0].InternalHTTPUrl, nil
}

// Fund funds a Stellar account via Friendbot, the funding faucet the local
// stellar/quickstart standalone network runs (funded from the standalone root
// account) — the analog of Aptos's local faucet, not a remote testnet service.
// Friendbot funds a fixed lumen amount and creates the account, so the amount
// argument is ignored (kept for interface parity with the EVM/Solana/Aptos
// funders). Polls until the faucet returns 200, mirroring the CLDF reference.
//
// NOTE: only exercised on the WRITE path. FundNodes (environment/dons.go) skips
// any chain unless flags.RequiresForwarderContract is true (or Solana), and
// Stellar reads (SimulateTransaction/View) submit no transactions, so this is
// not called during the read milestone.
func (s *Blockchain) Fund(ctx context.Context, address string, _ uint64) error {
	if s.friendbotURL == "" {
		return fmt.Errorf("no friendbot URL available for stellar chain %s", s.stellarChainID)
	}

	fundURL := s.friendbotURL + "?addr=" + url.QueryEscape(address)
	s.testLogger.Info().Msgf("Funding Stellar account %s via friendbot", address)

	deadline := time.Now().Add(friendbotFundTimeout)
	var lastErr error
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		status, body, err := friendbotRequest(ctx, fundURL)
		switch {
		case err != nil:
			lastErr = err
		case status == http.StatusOK:
			s.testLogger.Info().Msgf("Successfully funded Stellar account %s", address)
			return nil
		default:
			lastErr = fmt.Errorf("friendbot returned status %d: %s", status, body)
		}

		if time.Now().After(deadline) {
			return fmt.Errorf("friendbot funding failed for %s after %s: %w", address, friendbotFundTimeout, lastErr)
		}
		time.Sleep(friendbotFundInterval)
	}
}

func friendbotRequest(ctx context.Context, fundURL string) (int, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fundURL, nil)
	if err != nil {
		return 0, "", fmt.Errorf("failed to build friendbot request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(body), nil
}

// ToCldfChain builds a CLDF Stellar chain over the already-running RPC endpoint
// using a random deployer keypair (sufficient for read/View; the write milestone
// can swap in a funded KeypairGenerator for forwarder/receiver deployment). This
// mirrors how the Solana/Aptos blockchains expose a cldf chain for
// BlockChains.<Family>Chains() and saved state.
func (s *Blockchain) ToCldfChain() (cldf_chain.BlockChain, error) {
	nodeURL, err := s.NodeURL()
	if err != nil {
		return nil, err
	}
	prov := cldf_stellar_provider.NewRPCChainProvider(s.chainSelector, cldf_stellar_provider.RPCChainProviderConfig{
		SorobanRPCURL:      nodeURL,
		NetworkPassphrase:  s.networkPassphrase,
		FriendbotURL:       s.friendbotURL,
		DeployerKeypairGen: cldf_stellar_provider.KeypairRandom(),
	})
	return prov.Initialize(context.Background())
}

func (s *Deployer) Deploy(ctx context.Context, input *blockchain.Input) (blockchains.Blockchain, error) {
	var bcOut *blockchain.Output
	var err error

	switch {
	case s.provider.IsKubernetes():
		if err = blockchains.ValidateKubernetesBlockchainOutput(input); err != nil {
			return nil, err
		}
		s.testLogger.Info().Msgf("Using configured Kubernetes blockchain URLs for %s (chain_id: %s)", input.Type, input.ChainID)
		bcOut = input.Out
	case input.Out != nil:
		bcOut = input.Out
	default:
		bcOut, err = blockchain.NewWithContext(ctx, input)
		if err != nil {
			return nil, pkgerrors.Wrapf(err, "failed to deploy blockchain %s chainID: %s", input.Type, input.ChainID)
		}
	}

	// Framework Stellar output sets ChainID from input; fall back to config input.
	chainIDStr := bcOut.ChainID
	if chainIDStr == "" {
		chainIDStr = input.ChainID
	}
	if chainIDStr == "" {
		return nil, pkgerrors.New("stellar chain id is empty (set chain_id to the network id in [[blockchains]] in TOML)")
	}

	selector, ok := chainselectors.StellarChainIdToChainSelector()[chainIDStr]
	if !ok {
		return nil, fmt.Errorf("selector not found for stellar chainID '%s'", chainIDStr)
	}

	if len(bcOut.Nodes) == 0 {
		return nil, fmt.Errorf("no nodes found in stellar blockchain output for chainID '%s'", chainIDStr)
	}

	networkPassphrase := blockchain.DefaultStellarNetworkPassphrase
	var friendbotURL string
	if bcOut.NetworkSpecificData != nil && bcOut.NetworkSpecificData.StellarNetwork != nil {
		if bcOut.NetworkSpecificData.StellarNetwork.NetworkPassphrase != "" {
			networkPassphrase = bcOut.NetworkSpecificData.StellarNetwork.NetworkPassphrase
		}
		friendbotURL = bcOut.NetworkSpecificData.StellarNetwork.FriendbotURL
	}

	// Ensure ctfOutput has ChainID set for downstream (e.g. findStellarChains).
	bcOut.ChainID = chainIDStr

	return NewBlockchain(s.testLogger, chainIDStr, selector, friendbotURL, networkPassphrase, bcOut), nil
}
