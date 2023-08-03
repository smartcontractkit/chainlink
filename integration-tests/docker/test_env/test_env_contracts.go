package test_env

import (
	"math/big"
)

func (m *CLClusterTestEnv) WaitForEvents() error {
	return m.Geth.EthClient.WaitForEvents()
}

func (m *CLClusterTestEnv) DeployLINKToken() error {
	linkToken, err := m.Geth.ContractDeployer.DeployLinkTokenContract()
	if err != nil {
		return err
	}
	m.LinkToken = linkToken
	return err
}

func (m *CLClusterTestEnv) DeployMockETHLinkFeed(answer *big.Int) error {
	mockETHLINKFeed, err := m.Geth.ContractDeployer.DeployMockETHLINKFeed(answer)
	if err != nil {
		return err
	}
	m.MockETHLinkFeed = mockETHLINKFeed
	return err
}
