package v0_5

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/reward_manager_v0_5_0"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/contracts/evm"
	dsutil "github.com/smartcontractkit/chainlink/deployment/data-streams/utils"
)

// RecipientInfo represents the data for each recipient
type RecipientInfo struct {
	PayeeAddress string `json:"payeeAddress"`
	Weight       string `json:"weight"`
}

// RewardManagerView represents a processed view of reward manager data
type RewardManagerView struct {
	Owner            string                              `json:"owner,omitempty"`
	TypeAndVersion   string                              `json:"typeAndVersion,omitempty"`
	RecipientWeights map[string]map[string]RecipientInfo `json:"recipientWeights"` // poolId > -> payeeAddress -> recipient info
}

// SerializeView serializes view to JSON
func (r RewardManagerView) SerializeView() (string, error) {
	bytes, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal contract view: %w", err)
	}
	return string(bytes), nil
}

// RewardManagerViewParams defines parameters for generating the view
type RewardManagerViewParams struct {
	FromBlock uint64
	ToBlock   *uint64
}

// RewardManagerContract defines a minimal interface
type RewardManagerContract interface {
	TypeAndVersion(opts *bind.CallOpts) (string, error)
	Owner(opts *bind.CallOpts) (common.Address, error)
	// Methods to get pool IDs
	SRegisteredPoolIds(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error)

	// Methods to get pool and recipient data
	STotalRewardRecipientFees(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error)
	SRewardRecipientWeights(opts *bind.CallOpts, arg0 [32]byte, arg1 common.Address) (*big.Int, error)
	SRewardRecipientWeightsSet(opts *bind.CallOpts, arg0 [32]byte) (bool, error)

	// Event filtering methods
	FilterRewardRecipientsUpdated(opts *bind.FilterOpts, poolID [][32]byte) (evm.LogIterator[reward_manager_v0_5_0.RewardManagerRewardRecipientsUpdated], error)
}

// RewardManagerViewGenerator generates views for reward manager contracts
type RewardManagerViewGenerator struct {
	contract RewardManagerContract
}

// NewRewardManagerViewGenerator creates a new view generator
func NewRewardManagerViewGenerator(contract RewardManagerContract) *RewardManagerViewGenerator {
	return &RewardManagerViewGenerator{
		contract: contract,
	}
}

// Generate creates a view based on the provided parameters
func (g *RewardManagerViewGenerator) Generate(ctx context.Context, params RewardManagerViewParams) (RewardManagerView, error) {
	view := RewardManagerView{}

	if err := g.fetchContractState(ctx, &view); err != nil {
		return view, fmt.Errorf("failed to fetch contract state: %w", err)
	}

	// We need the pool IDs to get the recipient weights
	poolIDs, err := g.getPoolIDs(ctx)
	if err != nil {
		return view, fmt.Errorf("failed to get registered pool IDs: %w", err)
	}

	recipientWeights, err := g.getRecipientWeights(ctx, poolIDs, params)
	if err != nil {
		return view, fmt.Errorf("failed to get recipient weights: %w", err)
	}

	view.RecipientWeights = recipientWeights

	return view, nil
}

func (g *RewardManagerViewGenerator) fetchContractState(ctx context.Context, view *RewardManagerView) error {
	owner, err := g.contract.Owner(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("failed to get contract owner: %w", err)
	}
	view.Owner = owner.Hex()

	typeAndVersion, err := g.contract.TypeAndVersion(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("failed to get contract type and version: %w", err)
	}
	view.TypeAndVersion = typeAndVersion

	return nil
}

func (g *RewardManagerViewGenerator) getRecipientWeights(ctx context.Context, poolIDs [][32]byte, params RewardManagerViewParams) (map[string]map[string]RecipientInfo, error) {
	filterOpts := &bind.FilterOpts{
		Context: ctx,
		Start:   params.FromBlock,
		End:     params.ToBlock,
	}

	recipientWeights := make(map[string]map[string]RecipientInfo)

	updateIterator, err := g.contract.FilterRewardRecipientsUpdated(filterOpts, poolIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to filter reward recipients updated events: %w", err)
	}
	defer updateIterator.Close()

	// Collect events to collect all recipients and their weights
	for updateIterator.Next() {
		event := updateIterator.GetEvent()

		poolIDHex := dsutil.HexEncodeBytes32(event.PoolId)

		// Initialize the recipient map for this pool
		if _, exists := recipientWeights[poolIDHex]; !exists {
			recipientWeights[poolIDHex] = make(map[string]RecipientInfo)
		}

		for _, r := range event.NewRewardRecipients {
			recipient := r.Addr.Hex()
			weight := r.Weight

			recipientInfo := RecipientInfo{
				PayeeAddress: recipient,
				Weight:       strconv.FormatUint(weight, 10),
			}

			recipientWeights[poolIDHex][recipient] = recipientInfo
		}
	}

	if err := updateIterator.Error(); err != nil {
		return nil, fmt.Errorf("error iterating through reward recipients updated events: %w", err)
	}

	return recipientWeights, nil
}

func (g *RewardManagerViewGenerator) getPoolIDs(ctx context.Context) ([][32]byte, error) {
	var poolIDs [][32]byte
	for i := int64(0); ; i++ {
		poolID, err := g.contract.SRegisteredPoolIds(&bind.CallOpts{Context: ctx}, big.NewInt(i))
		if err != nil {
			// We'll see a revert when we reach the end of registered pool IDs.
			// Note: this is a simplification. We may want to check for a specific error type instead.
			break
		}

		weightsSet, err := g.contract.SRewardRecipientWeightsSet(&bind.CallOpts{Context: ctx}, poolID)
		if err != nil {
			return nil, fmt.Errorf("failed to check if weights are set for pool %x: %w", poolID, err)
		}

		if weightsSet {
			poolIDs = append(poolIDs, poolID)
		}
	}
	return poolIDs, nil
}
