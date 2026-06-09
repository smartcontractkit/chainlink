package pkg

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
)

var (
	ConfigStore MCMSConfigStore
)

type MCMSConfigStore struct {
	m map[string]*cldfproposalutils.MCMSWithTimelockConfig
}

func (r *MCMSConfigStore) Get(profileID string) (*cldfproposalutils.MCMSWithTimelockConfig, error) {
	profile, exists := r.m[profileID]
	if !exists {
		return nil, errors.New("mcms profile not found: " + profileID)
	}

	return profile, nil
}
func (r *MCMSConfigStore) List() []string {
	profiles := make([]string, 0, len(r.m))
	for profileID := range r.m {
		profiles = append(profiles, profileID)
	}

	return profiles
}

func (r *MCMSConfigStore) Put(profileID string, config cldfproposalutils.MCMSWithTimelockConfig) {
	r.m[profileID] = &config
}

func init() {
	ConfigStore = MCMSConfigStore{m: make(map[string]*cldfproposalutils.MCMSWithTimelockConfig)}
}

func MustGetMCMSConfig(quorum uint8, signers []common.Address, groupSigners []mcmstypes.Config) mcmstypes.Config {
	cfg, err := mcmstypes.NewConfig(quorum, signers, groupSigners)
	if err != nil {
		panic(fmt.Errorf("failed to create MCMS config: %w", err))
	}

	return cfg
}
