package clo

import (
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/clo/models"
)

//type DonEnvironment deployment.Environment

type DonEnvConfig struct {
	DonName string
	Chains  map[uint64]deployment.Chain
	Logger  logger.Logger
	Nops    []*models.NodeOperator
}

func NewDonEnv(cfg DonEnvConfig) *deployment.Environment {
	out := deployment.Environment{
		Name:     cfg.DonName,
		Offchain: NewJobClient(cfg.Logger, cfg.Nops),
		NodeIDs:  make([]string, 0),
		Chains:   cfg.Chains,
		Logger:   cfg.Logger,
	}
	// assume that all the nodes in the provided input nops are part of the don
	for _, nop := range cfg.Nops {
		for _, node := range nop.Nodes {
			out.NodeIDs = append(out.NodeIDs, node.ID)
		}
	}

	return &out
}
