package solclient

import (
	"github.com/gagliardetto/solana-go"
)

type LinkToken struct {
	Client        *Client
	Mint          *solana.Wallet
	MintAuthority *solana.Wallet
}

func (l *LinkToken) Address() string {
	return l.Mint.PublicKey().String()
}
