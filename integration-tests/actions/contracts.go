package actions

import (
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/seth"

	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

func GetLinkTokenContract(l zerolog.Logger, sethClient *seth.Client, configWithLinkToken tc.LinkTokenContractConfig) (*contracts.EthereumLinkToken, error) {
	if configWithLinkToken.UseExistingLinkTokenContract() {
		linkAddress, err := configWithLinkToken.GetLinkTokenContractAddress()
		if err != nil {
			return nil, err
		}

		return contracts.LoadLinkTokenContract(l, sethClient, linkAddress)
	}
	return contracts.DeployLinkTokenContract(l, sethClient)
}
