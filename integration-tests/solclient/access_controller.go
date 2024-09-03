package solclient

import (
	"github.com/gagliardetto/solana-go"

	access_controller2 "github.com/smartcontractkit/chainlink-solana/contracts/generated/access_controller"
)

type AccessController struct {
	Client        *Client
	State         *solana.Wallet
	Owner         *solana.Wallet
	ProgramWallet *solana.Wallet
}

func (s *AccessController) Address() string {
	return s.State.PublicKey().String()
}

func (s *AccessController) AddAccess(addr string) error {
	payer := s.Client.DefaultWallet
	validatorPubKey, err := solana.PublicKeyFromBase58(addr)
	if err != nil {
		return nil
	}
	err = s.Client.TXAsync(
		"Add validator access",
		[]solana.Instruction{
			access_controller2.NewAddAccessInstruction(
				s.State.PublicKey(),
				s.Owner.PublicKey(),
				validatorPubKey,
			).Build(),
		},
		func(key solana.PublicKey) *solana.PrivateKey {
			if key.Equals(s.Owner.PublicKey()) {
				return &s.Owner.PrivateKey
			}
			if key.Equals(payer.PublicKey()) {
				return &payer.PrivateKey
			}
			return nil
		},
		payer.PublicKey(),
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *AccessController) RemoveAccess(addr string) error {
	panic("implement me")
}

func (s *AccessController) HasAccess(to string) (bool, error) {
	panic("implement me")
}
