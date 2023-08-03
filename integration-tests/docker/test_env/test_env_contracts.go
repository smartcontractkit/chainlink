package test_env

import (
	"math/big"
)

func (te *CLClusterTestEnv) WaitForEvents() error {
	return te.Geth.EthClient.WaitForEvents()
}

func (te *CLClusterTestEnv) DeployLINKToken() error {
	linkToken, err := te.Geth.ContractDeployer.DeployLinkTokenContract()
	if err != nil {
		return err
	}
	te.LinkToken = linkToken
	return err
}

func (te *CLClusterTestEnv) DeployMockETHLinkFeed(answer *big.Int) error {
	mockETHLINKFeed, err := te.Geth.ContractDeployer.DeployMockETHLINKFeed(answer)
	if err != nil {
		return err
	}
	te.MockETHLinkFeed = mockETHLINKFeed
	return err
}
