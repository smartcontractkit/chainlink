package deployment

import (
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

// TonChain represents a TON chain.
type TonChain struct {
	Selector      uint64           // Canonical chain identifier
	Client        *ton.APIClient   // RPC client via Lite Server
	Wallet        *wallet.Wallet   // Wallet abstraction (signing, sending)
	WalletAddress *address.Address // Address of deployer wallet
	URL           string           // Liteserver URL
	DeployerSeed  string           // Optional: mnemonic or raw seed
}
