package offchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"strconv"
	"text/template"

	"github.com/google/uuid"
	chainsel "github.com/smartcontractkit/chain-selectors"
)

const (
	bootstrapPath = "bootstrap.tmpl"
)

type BootstrapJobCfg struct {
	JobName       string
	ExternalJobID string
	ContractID    string
	ChainID       string
}

func createBoostrapExternalJobID(donID uint64, evmChainSel uint64) (string, error) {
	in := []byte(strconv.FormatUint(donID, 10) + "-ocr3-capability-job-spec")
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, evmChainSel)
	in = append(in, b...)
	sha256Hash := sha256.New()
	sha256Hash.Write(in)
	in = sha256Hash.Sum(nil)[:16]
	// tag as valid UUID v4 https://github.com/google/uuid/blob/0f11ee6918f41a04c201eceeadf612a377bc7fbc/version4.go#L53-L54
	in[6] = (in[6] & 0x0f) | 0x40 // Version 4
	in[8] = (in[8] & 0x3f) | 0x80 // Variant is 10

	id, err := uuid.FromBytes(in)
	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func JobSpecFromBootstrap(donID uint64, chainSelector uint64, bootstrapJobName string, contract string) (string, error) {
	externalID, err := createBoostrapExternalJobID(donID, chainSelector)
	if err != nil {
		return "", fmt.Errorf("failed to get bootstrap external job id: %w", err)
	}

	chainID, err := chainsel.GetChainIDFromSelector(chainSelector)
	if err != nil {
		return "", fmt.Errorf("failed to get chain ID from selector: %w", err)
	}

	btCfg := BootstrapJobCfg{
		JobName:       bootstrapJobName,
		ExternalJobID: externalID,
		ContractID:    contract,
		ChainID:       chainID,
	}

	bootstrapJobSpec, err := btCfg.createSpec()
	if err != nil {
		return "", fmt.Errorf("failed to create bootstrap job spec: %w", err)
	}
	return bootstrapJobSpec, nil
}

func (btCfg *BootstrapJobCfg) createSpec() (string, error) {
	t, err := template.New("s").ParseFS(offchainFs, bootstrapPath)
	if err != nil {
		return "", err
	}

	b := &bytes.Buffer{}
	err = t.ExecuteTemplate(b, bootstrapPath, btCfg)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}
