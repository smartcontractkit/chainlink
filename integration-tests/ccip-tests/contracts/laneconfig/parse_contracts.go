package laneconfig

import (
	_ "embed"
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/multierr"
)

var (
	//go:embed contracts.json
	ExistingContracts []byte
	laneMu            = &sync.Mutex{}
)

type CommonContracts struct {
	IsNativeFeeToken bool     `json:"is_native_fee_token,omitempty"`
	IsMockARM        bool     `json:"is_mock_arm,omitempty"`
	FeeToken         string   `json:"fee_token"`
	BridgeTokens     []string `json:"bridge_tokens"`
	BridgeTokenPools []string `json:"bridge_tokens_pools"`
	ARM              string   `json:"arm"`
	Router           string   `json:"router"`
	PriceRegistry    string   `json:"price_registry"`
	WrappedNative    string   `json:"wrapped_native"`
}

type SourceContracts struct {
	OnRamp     string `json:"on_ramp"`
	DepolyedAt uint64 `json:"deployed_at"`
}

type DestContracts struct {
	OffRamp      string `json:"off_ramp"`
	CommitStore  string `json:"commit_store"`
	ReceiverDapp string `json:"receiver_dapp"`
}

type LaneConfig struct {
	CommonContracts
	SrcContracts  map[uint64]SourceContracts `json:"src_contracts"`  // key destination chain id
	DestContracts map[uint64]DestContracts   `json:"dest_contracts"` // key source chain id
}

func (l *LaneConfig) Validate() error {
	var laneConfigError error

	if l.ARM == "" || !common.IsHexAddress(l.ARM) {
		laneConfigError = multierr.Append(laneConfigError, errors.New("must set proper address for arm"))
	}

	if l.FeeToken == "" || !common.IsHexAddress(l.FeeToken) {
		laneConfigError = multierr.Append(laneConfigError, errors.New("must set proper address for fee_token"))
	}
	if len(l.BridgeTokens) < 1 {
		laneConfigError = multierr.Append(laneConfigError, errors.New("must set at least 1 bridge_tokens"))
	}
	for _, token := range l.BridgeTokens {
		if token == "" || !common.IsHexAddress(token) {
			laneConfigError = multierr.Append(laneConfigError, errors.New("must set proper address for bridge_tokens"))
		}
	}
	if len(l.BridgeTokenPools) < 1 {
		laneConfigError = multierr.Append(laneConfigError, errors.New("must set at least 1 bridge_tokens_pools"))
	}
	for _, pool := range l.BridgeTokenPools {
		if pool == "" || !common.IsHexAddress(pool) {
			laneConfigError = multierr.Append(laneConfigError, errors.New("must set proper address for bridge_tokens_pools"))
		}
	}
	if l.Router == "" || !common.IsHexAddress(l.Router) {
		laneConfigError = multierr.Append(laneConfigError, errors.New("must set proper address for router"))
	}
	if l.PriceRegistry == "" || !common.IsHexAddress(l.PriceRegistry) {
		laneConfigError = multierr.Append(laneConfigError, errors.New("must set proper address for price_registry"))
	}
	if l.WrappedNative == "" || !common.IsHexAddress(l.WrappedNative) {
		laneConfigError = multierr.Append(laneConfigError, errors.New("must set proper address for wrapped_native"))
	}
	return laneConfigError
}

type Lanes struct {
	LaneConfigs map[string]*LaneConfig `json:"lane_configs"`
}

func (l *Lanes) ReadLaneConfig(networkA string) (*LaneConfig, error) {
	laneMu.Lock()
	defer laneMu.Unlock()
	_, ok := l.LaneConfigs[networkA]
	if !ok {
		l.LaneConfigs[networkA] = &LaneConfig{}

	}
	return l.LaneConfigs[networkA], nil
}

func (l *Lanes) WriteLaneConfig(networkA string, cfg *LaneConfig) error {
	laneMu.Lock()
	defer laneMu.Unlock()
	if l.LaneConfigs == nil {
		l.LaneConfigs = make(map[string]*LaneConfig)
	}
	err := cfg.Validate()
	if err != nil {
		return err
	}
	l.LaneConfigs[networkA] = cfg
	return nil
}

func ReadLanesFromExistingDeployment() (*Lanes, error) {
	var lanes Lanes
	if err := json.Unmarshal(ExistingContracts, &lanes); err != nil {
		return nil, err
	}
	return &lanes, nil
}

func CreateDeploymentJSON(path string) (*Lanes, error) {
	existingLanes := Lanes{
		LaneConfigs: make(map[string]*LaneConfig),
	}
	err := WriteLanesToJSON(path, &existingLanes)
	return &existingLanes, err
}

func WriteLanesToJSON(path string, lanes *Lanes) error {
	b, err := json.MarshalIndent(lanes, "", "  ")
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(b)
	return err
}
