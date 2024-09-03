package solclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	confirm "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/gagliardetto/solana-go/text"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

var (
	DefaultPrivateKeysSolValidator = []string{
		"57qbvFjTChfNwQxqkFZwjHp7xYoPZa7f9ow6GA59msfCH1g6onSjKUTrrLp4w1nAwbwQuit8YgJJ2AwT9BSwownC",
		"2tye1GyG7wwTUS2T8puXSErDyzQcBxpgwRN5R2MMy5osJKjQF6ZoeYTTpeHaAxpuiE1G4Pnq4sTa4YCWx3RcXb4Y",
		"5aRBAnU3NBymRyMtrRjPLZ3erZNgTZBhEHszsXF8kTwbGLz8q5FYgKicJ7AFifrFitvJB2NS5jbyQohSJtvkgPER",
		"2MYG6HKpMuGEo3qErj4pAF2Gney6Yb6jgjTc6TZCuu7fiLAVQekTd3HbsT9ienzGHpKwA7Ekj2TGuMHPUB6EHJ8P",
	}
)

// Accounts is a shared state between contracts in which data is stored in Solana
type Accounts struct {
	// OCR OCR program state account
	OCR *solana.Wallet
	// Store store program state account
	Store *solana.Wallet
	// OCRVault OCR program account to hold LINK
	OCRVault                 *solana.Wallet
	OCRVaultAssociatedPubKey solana.PublicKey
	// Transmissions OCR transmissions state account
	Feed *solana.Wallet
	// Authorities authorities used to sign on-chain, used by programs
	Authorities map[string]*Authority
	// Owner is the owner of all programs
	Owner *solana.Wallet
	// Mint LINK mint state account
	Mint *solana.Wallet
	// OCR2 Proposal account
	Proposal *solana.Wallet
	// MintAuthority LINK mint authority
	MintAuthority *solana.Wallet
}

type SolNetwork struct {
	External          bool          `mapstructure:"external" yaml:"external"`
	ContractsDeployed bool          `mapstructure:"contracts_deployed" yaml:"contracts_deployed"`
	Name              string        `mapstructure:"name" yaml:"name"`
	ID                string        `mapstructure:"id" yaml:"id"`
	ChainID           int64         `mapstructure:"chain_id" yaml:"chain_id"`
	URL               string        `mapstructure:"url" yaml:"url"`
	URLs              []string      `mapstructure:"urls" yaml:"urls"`
	Type              string        `mapstructure:"type" yaml:"type"`
	PrivateKeys       []string      `mapstructure:"private_keys" yaml:"private_keys"`
	Timeout           time.Duration `mapstructure:"transaction_timeout" yaml:"transaction_timeout"`
}

// Client implements BlockchainClient
type Client struct {
	Accounts *Accounts
	Config   *SolNetwork
	// Wallets lamport wallets
	Wallets []*solana.Wallet
	// ProgramWallets program wallets by key filename
	ProgramWallets    map[string]*solana.Wallet
	DefaultWallet     *solana.Wallet
	txErrGroup        errgroup.Group
	queueTransactions bool
	// RPC rpc client
	RPC *rpc.Client
	// WS websocket client
	WS        *ws.Client
	LinkToken *LinkToken
}

func (c *Client) GetNetworkType() string {
	return c.Config.Type
}

func (c *Client) ContractsDeployed() bool {
	return c.Config.ContractsDeployed
}

func (c *Client) EstimateCostForChainlinkOperations(amountOfOperations int) (*big.Float, error) {
	panic("implement me")
}

// NewClient creates new Solana client both for RPC ans WS
func NewClient(cfg *SolNetwork) (*Client, error) {
	c := rpc.New(cfg.URLs[0])
	wsc, err := ws.Connect(context.Background(), cfg.URLs[1])
	if err != nil {
		return nil, err
	}
	client := &Client{
		Config:         cfg,
		RPC:            c,
		WS:             wsc,
		ProgramWallets: make(map[string]*solana.Wallet),
		txErrGroup:     errgroup.Group{},
	}
	if err := client.LoadWallets(cfg); err != nil {
		return nil, err
	}
	return client, nil
}

// CreateAccInstr creates instruction for account creation of particular size
func (c *Client) CreateAccInstr(acc solana.PublicKey, accSize uint64, ownerPubKey solana.PublicKey) (solana.Instruction, error) {
	payer := c.DefaultWallet
	rentMin, err := c.RPC.GetMinimumBalanceForRentExemption(
		context.TODO(),
		accSize,
		rpc.CommitmentConfirmed,
	)
	if err != nil {
		return nil, err
	}
	return system.NewCreateAccountInstruction(
		rentMin,
		accSize,
		ownerPubKey,
		payer.PublicKey(),
		acc,
	).Build(), nil
}

// TXSync executes tx synchronously with specified commitment (defaults to finalized)
func (c *Client) TXSync(name string, commitment rpc.CommitmentType, instr []solana.Instruction, signerFunc func(key solana.PublicKey) *solana.PrivateKey, payer solana.PublicKey) error {
	recent, err := c.RPC.GetRecentBlockhash(context.Background(), rpc.CommitmentFinalized)
	if err != nil {
		return err
	}
	tx, err := solana.NewTransaction(
		instr,
		recent.Value.Blockhash,
		solana.TransactionPayer(payer),
	)
	if err != nil {
		return err
	}
	if _, err = tx.EncodeTree(text.NewTreeEncoder(os.Stdout, name)); err != nil {
		return err
	}
	if _, err = tx.Sign(signerFunc); err != nil {
		return err
	}
	sig, err := c.RPC.SendTransactionWithOpts(
		context.Background(),
		tx,
		rpc.TransactionOpts{
			SkipPreflight:       false,
			PreflightCommitment: commitment,
		},
	)
	if err != nil {
		return err
	}
	log.Info().Interface("Sig", sig).Msg("TX committed")
	sub, err := c.WS.SignatureSubscribe(
		sig,
		commitment,
	)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()
	res, err := sub.Recv()
	if err != nil {
		return err
	}
	log.Debug().Interface("TX", res).Msg("TX response")
	return nil
}

func (c *Client) queueTX(sig solana.Signature, commitment rpc.CommitmentType) {
	c.txErrGroup.Go(func() error {
		sub, err := c.WS.SignatureSubscribe(
			sig,
			commitment,
		)
		if err != nil {
			return err
		}
		defer sub.Unsubscribe()
		res, err := sub.Recv()
		if err != nil {
			return err
		}
		if res.Value.Err != nil {
			return fmt.Errorf("transaction confirmation failed: %v", res.Value.Err)
		}
		return nil
	})
}

// TXAsync executes tx async, need to block on WaitForEvents after
func (c *Client) TXAsync(name string, instr []solana.Instruction, signerFunc func(key solana.PublicKey) *solana.PrivateKey, payer solana.PublicKey) error {
	recent, err := c.RPC.GetRecentBlockhash(context.Background(), rpc.CommitmentFinalized)
	if err != nil {
		return err
	}
	tx, err := solana.NewTransaction(
		instr,
		recent.Value.Blockhash,
		solana.TransactionPayer(payer),
	)
	if err != nil {
		return err
	}
	if _, err = tx.EncodeTree(text.NewTreeEncoder(os.Stdout, name)); err != nil {
		return err
	}
	if _, err = tx.Sign(signerFunc); err != nil {
		return err
	}
	sig, err := c.RPC.SendTransactionWithOpts(
		context.Background(),
		tx,
		rpc.TransactionOpts{
			PreflightCommitment: rpc.CommitmentConfirmed,
		},
	)

	if err != nil {
		return err
	}
	c.queueTX(sig, rpc.CommitmentConfirmed)
	log.Info().Interface("Sig", sig).Msg("TX send")
	return nil
}

// LoadWallet loads wallet from path
func (c *Client) LoadWallet(path string) (*solana.Wallet, error) {
	pk, err := solana.PrivateKeyFromBase58(path)
	if err != nil {
		return nil, err
	}
	log.Debug().
		Str("PrivKey", pk.String()).
		Str("PubKey", pk.PublicKey().String()).
		Msg("Loaded wallet")
	return &solana.Wallet{PrivateKey: pk}, nil
}

// Airdrop airdrops a wallet with lamports
func (c *Client) Airdrop(wpk solana.PublicKey, solAmount uint64) error {
	txHash, err := c.RPC.RequestAirdrop(
		context.Background(),
		wpk,
		solana.LAMPORTS_PER_SOL*solAmount,
		rpc.CommitmentConfirmed,
	)
	if err != nil {
		return err
	}
	log.Debug().
		Str("PublicKey", wpk.String()).
		Str("TX", txHash.String()).
		Msg("Airdropping account")
	c.queueTX(txHash, rpc.CommitmentConfirmed)
	return nil
}

func (c *Client) AirdropAddresses(addr []string, solAmount uint64) error {
	for _, a := range addr {
		pubKey, err := solana.PublicKeyFromBase58(a)
		if err != nil {
			return err
		}
		if err := c.Airdrop(pubKey, solAmount); err != nil {
			return err
		}
	}
	return c.WaitForEvents()
}

// ListDirFilenamesByExt returns all the filenames inside a dir for file with particular extension, for ex. ".json"
func (c *Client) ListDirFilenamesByExt(dir string, ext string) ([]string, error) {
	keyFiles := make([]string, 0)
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ext {
			keyFiles = append(keyFiles, info.Name())
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return keyFiles, nil
}

// LoadWallets loads wallets from config
func (c *Client) LoadWallets(nc interface{}) error {
	cfg := nc.(*SolNetwork)
	for _, pkString := range cfg.PrivateKeys {
		w, err := c.LoadWallet(pkString)
		if err != nil {
			return err
		}
		c.Wallets = append(c.Wallets, w)
	}
	addresses := make([]string, 0)
	for _, w := range c.Wallets {
		addresses = append(addresses, w.PublicKey().String())
	}
	if err := c.AirdropAddresses(addresses, 500); err != nil {
		return err
	}
	if err := c.SetWallet(1); err != nil {
		return err
	}
	log.Debug().Interface("Wallets", c.Wallets).Msg("Common wallets")
	return nil
}

// SetWallet sets default client
func (c *Client) SetWallet(num int) error {
	c.DefaultWallet = c.Wallets[num]
	return nil
}

func (c *Client) CalculateTXSCost(txs int64) (*big.Float, error) {
	panic("implement me")
}

func (c *Client) CalculateTxGas(gasUsedValue *big.Int) (*big.Float, error) {
	panic("implement me")
}

func (c *Client) Get() interface{} {
	return c
}

func (c *Client) GetNetworkName() string {
	return c.Config.Name
}

func (c *Client) SwitchNode(node int) error {
	panic("implement me")
}

func (c *Client) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	panic("implement me")
}

func (c *Client) HeaderHashByNumber(ctx context.Context, bn *big.Int) (string, error) {
	panic("implement me")
}

func (c *Client) BlockNumber(ctx context.Context) (uint64, error) {
	panic("implement me")
}

func (c *Client) HeaderTimestampByNumber(ctx context.Context, bn *big.Int) (uint64, error) {
	panic("implement me")
}

func (c *Client) Fund(toAddress string, amount *big.Float) error {
	pubKey, err := solana.PublicKeyFromBase58(toAddress)
	if err != nil {
		return err
	}
	a, _ := amount.Uint64()
	txHash, err := c.RPC.RequestAirdrop(
		context.Background(),
		pubKey,
		solana.LAMPORTS_PER_SOL*a,
		rpc.CommitmentConfirmed,
	)
	if err != nil {
		return err
	}
	log.Debug().
		Str("PublicKey", pubKey.String()).
		Str("TX", txHash.String()).
		Msg("Airdropping account")
	c.queueTX(txHash, rpc.CommitmentConfirmed)
	return nil
}

func (cfg *SolNetwork) Default() *SolNetwork {
	return &SolNetwork{
		Name:              "sol",
		Type:              "solana",
		ContractsDeployed: false,
		PrivateKeys:       DefaultPrivateKeysSolValidator,
		URLs: []string{
			"http://127.0.0.1:8899",
			"ws://127.0.0.1:8900",
		},
	}
}

func (c *Client) ParallelTransactions(enabled bool) {
	c.queueTransactions = enabled
}

func (c *Client) Close() error {
	c.WS.Close()
	return nil
}

func (c *Client) EstimateTransactionGasCost() (*big.Int, error) {
	panic("implement me")
}

func (c *Client) DeleteHeaderEventSubscription(key string) {
	panic("implement me")
}

func (c *Client) WaitForEvents() error {
	return c.txErrGroup.Wait()
}

func (c *Client) GetChainID() *big.Int {
	panic("implement me")
}

func (c *Client) GetClients() []blockchain.EVMClient {
	panic("implement me")
}

func (c *Client) GetDefaultWallet() *blockchain.EthereumWallet {
	panic("implement me")
}

func (c *Client) GetWallets() []*blockchain.EthereumWallet {
	panic("implement me")
}

func (c *Client) SetID(id int) {
	panic("implement me")
}

func (c *Client) SetDefaultWallet(num int) error {
	panic("implement me")
}

func (c *Client) SetWallets(wallets []*blockchain.EthereumWallet) {
	panic("implement me")
}

func (c *Client) LatestBlockNumber(ctx context.Context) (uint64, error) {
	panic("implement me")
}

func (c *Client) GasStats() *blockchain.GasStats {
	panic("implement me")
}

func (c *Client) AddHeaderEventSubscription(key string, subscriber blockchain.HeaderEventSubscription) {
	panic("implement me")
}

func SendFunds(senderPrivateKey string, receiverPublicKey string, lamports uint64, rpcClient *rpc.Client, wsClient *ws.Client) error {
	// Convert the private key string to a byte slice
	var privateKeyBytes []byte
	err := json.Unmarshal([]byte(senderPrivateKey), &privateKeyBytes)
	if err != nil {
		return err
	}

	accountFrom := solana.PrivateKey(privateKeyBytes)
	accountTo := solana.MustPublicKeyFromBase58(receiverPublicKey)

	// Get recent blockhash
	recent, err := rpcClient.GetRecentBlockhash(context.Background(), rpc.CommitmentFinalized)
	if err != nil {
		return err
	}

	// Create a transfer transaction
	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			system.NewTransferInstruction(
				lamports,
				accountFrom.PublicKey(),
				accountTo,
			).Build(),
		},
		recent.Value.Blockhash,
		solana.TransactionPayer(accountFrom.PublicKey()),
	)
	if err != nil {
		return err
	}

	// Sign the transaction
	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if accountFrom.PublicKey().Equals(key) {
				return &accountFrom
			}
			return nil
		},
	)
	if err != nil {
		return err
	}

	// Send transaction, and wait for confirmation:
	_, err = confirm.SendAndConfirmTransaction(
		context.Background(),
		rpcClient,
		wsClient,
		tx,
	)
	if err != nil {
		return err
	}

	return nil
}
