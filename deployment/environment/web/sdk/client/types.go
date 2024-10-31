package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

type JobDistributorInput struct {
	Name      string `json:"name"`
	Uri       string `json:"uri"`
	PublicKey string `json:"publicKey"`
}

type JobDistributorChainConfigInput struct {
	JobDistributorID     string `json:"feedsManagerID"`
	ChainID              string `json:"chainID"`
	ChainType            string `json:"chainType"`
	AccountAddr          string `json:"accountAddr"`
	AccountAddrPubKey    string `json:"accountAddrPubKey"`
	AdminAddr            string `json:"adminAddr"`
	FluxMonitorEnabled   bool   `json:"fluxMonitorEnabled"`
	Ocr1Enabled          bool   `json:"ocr1Enabled"`
	Ocr1IsBootstrap      bool   `json:"ocr1IsBootstrap"`
	Ocr1Multiaddr        string `json:"ocr1Multiaddr"`
	Ocr1P2PPeerID        string `json:"ocr1P2PPeerID"`
	Ocr1KeyBundleID      string `json:"ocr1KeyBundleID"`
	Ocr2Enabled          bool   `json:"ocr2Enabled"`
	Ocr2IsBootstrap      bool   `json:"ocr2IsBootstrap"`
	Ocr2Multiaddr        string `json:"ocr2Multiaddr"`
	Ocr2ForwarderAddress string `json:"ocr2ForwarderAddress"`
	Ocr2P2PPeerID        string `json:"ocr2P2PPeerID"`
	Ocr2KeyBundleID      string `json:"ocr2KeyBundleID"`
	Ocr2Plugins          string `json:"ocr2Plugins"`
}

type JobProposalApprovalSuccessSpec struct {
	Id              string `json:"id"`
	Definition      string `json:"definition"`
	Version         int    `json:"version"`
	Status          string `json:"status"`
	StatusUpdatedAt string `json:"statusUpdatedAt"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

func DecodeInput(in, out any) error {
	if reflect.TypeOf(out).Kind() != reflect.Ptr || reflect.ValueOf(out).IsNil() {
		return fmt.Errorf("out type must be a non-nil pointer")
	}
	jsonBytes, err := json.Marshal(in)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonBytes))
	decoder.DisallowUnknownFields()
	return decoder.Decode(out)
}
