package handler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/olekukonko/tablewriter"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"

	ocr2keepers20config "github.com/smartcontractkit/chainlink-automation/pkg/v2/config"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/keeper_registry_wrapper2_0"
)

func OCR2GetConfig(hdlr *baseHandler, registry_addr string) error {
	b, err := common.ParseHexOrString(registry_addr)
	if err != nil {
		return fmt.Errorf("failed to parse address hash: %w", err)
	}

	addr := common.BytesToAddress(b)
	registry, err := keeper_registry_wrapper2_0.NewKeeperRegistry(addr, hdlr.client)
	if err != nil {
		return fmt.Errorf("failed to create caller for address and backend: %w", err)
	}

	log.Printf("getting config details from contract: %s\n", addr.Hex())
	detail, err := registry.LatestConfigDetails(nil)
	if err != nil {
		return fmt.Errorf("failed to get latest config detail from contract: %w", err)
	}

	block, err := hdlr.client.BlockByNumber(context.Background(), big.NewInt(int64(detail.BlockNumber)))
	if err != nil {
		return fmt.Errorf("failed to get block at number %d: %w", detail.BlockNumber, err)
	}

	config, err := configFromBlock(block, addr, detail)
	if err != nil {
		return fmt.Errorf("failed to get config from block: %w", err)
	}

	printConfigValues(config)
	return nil
}

func configFromBlock(bl *types.Block, addr common.Address, detail keeper_registry_wrapper2_0.LatestConfigDetails) (*confighelper.PublicConfig, error) {
	for _, tx := range bl.Transactions() {
		if tx.To() != nil && bytes.Equal(tx.To()[:], addr[:]) {
			// this is our transaction
			// txRes, txErr, err := getTransactionDetailForHashes(hdlr, []string{tx})
			ocr2Tx, err := NewBaseOCR2Tx(tx)
			if err != nil {
				log.Printf("failed to create set config transaction: %s", err)
				continue
			}

			method, err := ocr2Tx.Method()
			if err != nil {
				log.Printf("failed to parse method signature: %s", err)
				continue
			}

			if method.Name == "setConfig" {
				log.Printf("found transaction for last config update: %s", ocr2Tx.Hash())
				confTx, err := NewOCR2SetConfigTx(tx)
				if err != nil {
					log.Printf("failed to create conf tx: %s", err)
					continue
				}

				conf, err := confTx.Config()
				if err != nil {
					log.Printf("failed to parse transaction config: %s", err)
				}
				conf.ConfigCount = uint64(detail.ConfigCount)
				conf.ConfigDigest = detail.ConfigDigest

				pubConf, err := confighelper.PublicConfigFromContractConfig(true, conf)
				if err != nil {
					log.Printf("failed to parse public config: %s", err)
				}

				return &pubConf, nil
			}
		}
	}

	return nil, errors.New("public config not found")
}

func printConfigValues(config *confighelper.PublicConfig) {
	data := [][]string{}

	data = append(data, []string{"DeltaProgress", config.DeltaProgress.String()})
	data = append(data, []string{"DeltaResend", config.DeltaResend.String()})
	data = append(data, []string{"DeltaRound", config.DeltaRound.String()})
	data = append(data, []string{"DeltaGrace", config.DeltaGrace.String()})
	data = append(data, []string{"DeltaStage", config.DeltaStage.String()})
	data = append(data, []string{"RMax", strconv.FormatUint(uint64(config.RMax), 10)})
	data = append(data, []string{"S", fmt.Sprintf("%v", config.S)})
	data = append(data, []string{"MaxDurationQuery", config.MaxDurationQuery.String()})
	data = append(data, []string{"MaxDurationObservation", config.MaxDurationObservation.String()})
	data = append(data, []string{"MaxDurationReport", config.MaxDurationReport.String()})
	data = append(data, []string{"MaxDurationShouldAcceptFinalizedReport", config.MaxDurationShouldAcceptFinalizedReport.String()})
	data = append(data, []string{"MaxDurationShouldTransmitAcceptedReport", config.MaxDurationShouldTransmitAcceptedReport.String()})
	data = append(data, []string{"F", strconv.Itoa(config.F)})

	if offConf, err := ocr2keepers20config.DecodeOffchainConfig(config.ReportingPluginConfig); err == nil {
		data = append(data, []string{"", ""})
		data = append(data, []string{"TargetProbability", offConf.TargetProbability})
		data = append(data, []string{"GasLimitPerReport", strconv.FormatUint(uint64(offConf.GasLimitPerReport), 10)})
		data = append(data, []string{"GasOverheadPerUpkeep", strconv.FormatUint(uint64(offConf.GasOverheadPerUpkeep), 10)})
		data = append(data, []string{"MinConfirmations", strconv.Itoa(offConf.MinConfirmations)})
		data = append(data, []string{"PerformLockoutWindow", strconv.FormatInt(offConf.PerformLockoutWindow, 10)})
		data = append(data, []string{"SamplingJobDuration", strconv.FormatInt(offConf.SamplingJobDuration, 10)})
		data = append(data, []string{"TargetInRounds", strconv.Itoa(offConf.TargetInRounds)})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Field", "Value"})
	// table.SetFooter([]string{"", "", "Total", "$146.93"}) // Add Footer
	table.SetBorder(false) // Set Border to false
	table.AppendBulk(data) // Add Bulk Data
	table.Render()
}
