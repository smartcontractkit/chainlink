package gogauntlet

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

type NextAction string

var (
	ApproveNextAction NextAction = "approve"
	ExecuteNextAction NextAction = "execute"
)

type BigNumber struct {
	Type string `json:"type"`
	Hex  string `json:"hex"`
}

// SafeState is the state of the multisig safe and the proposal associated with the :safe call.
type SafeState struct {
	Safe struct {
		Address   string    `json:"address"`
		Threshold BigNumber `json:"threshold"`
		Signers   []string  `json:"signers"`
		Nonce     BigNumber `json:"nonce"`
	} `json:"safe"`
	Proposal struct {
		ID                     string     `json:"id"`
		Confirmations          uint64     `json:"confirmations"`
		SelfApproved           bool       `json:"selfApproved"`
		HasEnoughConfirmations bool       `json:"hasEnoughConfirmations"`
		NextAction             NextAction `json:"nextAction"`
		Approvers              []string   `json:"approvers"`
	} `json:"proposal"`
}

type Report struct {
	Responses []struct {
		Tx struct {
			Hash    string `json:"hash"`
			Address string `json:"address"`
			Status  string `json:"status"`
			Tx      struct {
				Data string `json:"data"`
				To   string `json:"to"`
			} `json:"tx"`
		} `json:"tx"`
		Contract string `json:"contract"`
	} `json:"responses"`
	Data struct {
		State    SafeState `json:"state,omitempty"` // State is only returned for :safe calls.
		Messages []struct {
			To    string `json:"to"`
			Data  string `json:"data"`
			Value uint16 `json:"value"`
		} `json:"messages"`
	} `json:"data"`
}

func (g Gauntlet) parseJsonReport(reportName string) (Report, error) {
	jsonFile, err := os.Open(fmt.Sprintf("%s.json", reportName))
	if err != nil {
		return Report{}, err
	}
	defer jsonFile.Close()

	var report Report
	byteValue, _ := io.ReadAll(jsonFile)
	if err = json.Unmarshal(byteValue, &report); err != nil {
		return Report{}, err
	}

	return report, nil
}

func (r Report) GetExportData() (*string, error) {
	if r.Data.Messages != nil && r.Data.Messages[0].Data != "" {
		return &(r.Data.Messages[0].Data), nil
	}
	return nil, errors.New("could not extract export data from report")
}
