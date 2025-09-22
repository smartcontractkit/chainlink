package jobs

import (
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"

	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
)

func verifyEVMJobSpecInputs(inputs job_types.JobSpecInput) error {
	if v, ok := inputs["command"]; !ok {
		return errors.New("command is required and must be a string")
	} else if s, ok := v.(string); !ok || strings.TrimSpace(s) == "" {
		return errors.New("command is required and must be a string")
	}

	if v, ok := inputs["config"]; !ok {
		return errors.New("config is required and must be a string")
	} else if s, ok := v.(string); !ok || strings.TrimSpace(s) == "" {
		return errors.New("config is required and must be a string")
	}

	of, err := job_types.DecodeOracleFactory(inputs)
	if err != nil {
		return err
	}

	if !of.Enabled {
		return errors.New("oracleFactory.enabled must be true for EVM jobs")
	}

	if len(of.BootstrapPeers) == 0 {
		return errors.New("oracleFactory.bootstrapPeers is required")
	}
	if _, err := ocrcommon.ParseBootstrapPeers(of.BootstrapPeers); err != nil {
		return errors.New("oracleFactory.bootstrapPeers is invalid: " + err.Error())
	}

	if strings.TrimSpace(of.OCRContractAddress) == "" {
		return errors.New("oracleFactory.ocrContractAddress is required")
	}

	if !common.IsHexAddress(of.OCRContractAddress) {
		return errors.New("oracleFactory.ocrContractAddress is invalid: not a checksummed hex address")
	}

	if strings.TrimSpace(of.OCRKeyBundleID) == "" {
		return errors.New("oracleFactory.ocrKeyBundleID is required")
	}

	if strings.TrimSpace(of.ChainID) == "" {
		return errors.New("oracleFactory.chainID is required")
	}

	if _, err := chain_selectors.GetChainDetailsByChainIDAndFamily(of.ChainID, chain_selectors.FamilyEVM); err != nil {
		return errors.New("oracleFactory.chainID is invalid: " + err.Error())
	}

	if strings.TrimSpace(of.TransmitterID) == "" {
		return errors.New("oracleFactory.transmitterID is required")
	}

	if strings.TrimSpace(of.OnchainSigningStrategy.StrategyName) == "" {
		return errors.New("oracleFactory.onchainSigningStrategy.strategyName is required")
	}

	if of.OnchainSigningStrategy.Config == nil {
		return errors.New("oracleFactory.onchainSigningStrategy.config is required")
	}
	if v, ok := of.OnchainSigningStrategy.Config["evm"]; !ok || strings.TrimSpace(v) == "" {
		return errors.New("oracleFactory.onchainSigningStrategy.config.evm is required")
	}

	return nil
}
