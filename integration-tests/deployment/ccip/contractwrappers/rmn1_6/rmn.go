package rmn1_6

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_remote"
)

type RMN struct {
	*rmn_remote.RMNRemote
}

func (rmn *RMN) GetVersionedConfig(opts *bind.CallOpts) (view.RMNRemoteVersionedConfig, error) {
	config, err := rmn.RMNRemote.GetVersionedConfig(opts)
	if err != nil {
		return view.RMNRemoteVersionedConfig{}, err
	}
	var signers []view.RMNRemoteSigner
	for _, signer := range config.Config.Signers {
		signers = append(signers, view.RMNRemoteSigner{
			OnchainPublicKey: signer.OnchainPublicKey.Hex(),
			NodeIndex:        signer.NodeIndex,
		})
	}
	return view.RMNRemoteVersionedConfig{
		Version:    config.Version,
		MinSigners: config.Config.MinSigners,
		Signers:    signers,
	}, nil

}

func New(rmn *rmn_remote.RMNRemote) *RMN {
	return &RMN{rmn}
}
