package v0_5

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	fee_manager "github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/fee_manager_v0_5_0"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/contracts/evm"
	dsutil "github.com/smartcontractkit/chainlink/deployment/data-streams/utils"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/view/interfaces"
)

// FeeManagerView represents a view of the FeeManager contract state
type FeeManagerView struct {
	LinkAddress         string                               `json:"linkAddress"`
	NativeAddress       string                               `json:"nativeAddress"`
	ProxyAddress        string                               `json:"proxyAddress"`
	RewardManager       string                               `json:"rewardManager"`
	NativeSurcharge     string                               `json:"nativeSurcharge"`
	LinkAvailable       string                               `json:"linkAvailable"`
	TypeAndVersion      string                               `json:"typeAndVersion,omitempty"`
	Owner               string                               `json:"owner,omitempty"`
	SubscriberDiscounts map[string]map[string]TokenDiscounts `json:"subscriberDiscounts"` // Map[subscriberAddress][feedId]TokenDiscounts
}

type TokenDiscounts struct {
	Link     string `json:"link"`
	Native   string `json:"native"`
	IsGlobal bool   `json:"isGlobal"`
}

// FeeManagerView implements the ContractView interface
var _ interfaces.ContractView = (*FeeManagerView)(nil)

// SerializeView serializes view to JSON
func (v FeeManagerView) SerializeView() (string, error) {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal contract view: %w", err)
	}
	return string(bytes), nil
}

// FeeManagerViewParams represents parameters for generating a FeeManager view
// In this simple case, we might not need parameters, but including as an example
type FeeManagerViewParams struct {
	FromBlock uint64
	ToBlock   *uint64
}

// FeeManagerViewGenerator implements ContractViewGenerator for FeeManager
type FeeManagerViewGenerator struct {
	contract FeeManagerReader
}

// FeeManagerViewGenerator implements ContractViewGenerator
var _ interfaces.ContractViewGenerator[FeeManagerViewParams, FeeManagerView] = (*FeeManagerViewGenerator)(nil)

// FeeManagerReader defines the minimal interface needed for FeeManagerViewGenerator
type FeeManagerReader interface {
	// Call methods
	TypeAndVersion(opts *bind.CallOpts) (string, error)
	Owner(opts *bind.CallOpts) (common.Address, error)
	ILinkAddress(opts *bind.CallOpts) (common.Address, error)
	INativeAddress(opts *bind.CallOpts) (common.Address, error)
	IProxyAddress(opts *bind.CallOpts) (common.Address, error)
	IRewardManager(opts *bind.CallOpts) (common.Address, error)
	SNativeSurcharge(opts *bind.CallOpts) (*big.Int, error)
	LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error)
	SGlobalDiscounts(opts *bind.CallOpts, subscriber common.Address, token common.Address) (*big.Int, error)
	SSubscriberDiscounts(opts *bind.CallOpts, subscriber common.Address, feedID [32]byte, token common.Address) (*big.Int, error)

	// Filter methods
	FilterSubscriberDiscountUpdated(opts *bind.FilterOpts, subscriber []common.Address, feedID [][32]byte) (evm.LogIterator[fee_manager.FeeManagerSubscriberDiscountUpdated], error)
}

func NewFeeManagerViewGenerator(contract FeeManagerReader) *FeeManagerViewGenerator {
	return &FeeManagerViewGenerator{
		contract: contract,
	}
}

func (f *FeeManagerViewGenerator) Generate(ctx context.Context, params FeeManagerViewParams) (FeeManagerView, error) {
	// Create and return the view
	view := &FeeManagerView{}
	if err := f.fetchContractState(ctx, view); err != nil {
		return *view, fmt.Errorf("failed to fetch contract state: %w", err)
	}

	discounts, err := f.gatherOrganizedDiscounts(ctx, params, view)
	if err != nil {
		return *view, fmt.Errorf("failed to gather organized discounts: %w", err)
	}

	// Create and return the view
	view.SubscriberDiscounts = discounts

	return *view, nil
}

func (f *FeeManagerViewGenerator) fetchContractState(ctx context.Context, view *FeeManagerView) error {
	callOpts := &bind.CallOpts{Context: ctx}

	owner, err := f.contract.Owner(callOpts)
	if err != nil {
		return fmt.Errorf("failed to get owner: %w", err)
	}
	view.Owner = owner.Hex()

	linkAddress, err := f.contract.ILinkAddress(callOpts)
	if err != nil {
		return fmt.Errorf("failed to get link address: %w", err)
	}
	view.LinkAddress = linkAddress.Hex()

	nativeAddress, err := f.contract.INativeAddress(callOpts)
	if err != nil {
		return fmt.Errorf("failed to get native address: %w", err)
	}
	view.NativeAddress = nativeAddress.Hex()

	proxyAddress, err := f.contract.IProxyAddress(callOpts)
	if err != nil {
		return fmt.Errorf("failed to get proxy address: %w", err)
	}
	view.ProxyAddress = proxyAddress.Hex()

	rewardManager, err := f.contract.IRewardManager(callOpts)
	if err != nil {
		return fmt.Errorf("failed to get reward manager: %w", err)
	}
	view.RewardManager = rewardManager.Hex()

	nativeSurcharge, err := f.contract.SNativeSurcharge(callOpts)
	if err != nil {
		return fmt.Errorf("failed to get native surcharge: %w", err)
	}
	view.NativeSurcharge = nativeSurcharge.String()

	linkAvailable, err := f.contract.LinkAvailableForPayment(callOpts)
	if err != nil {
		return fmt.Errorf("failed to get link available: %w", err)
	}
	view.LinkAvailable = linkAvailable.String()

	typeAndVersion, err := f.contract.TypeAndVersion(callOpts)
	if err != nil {
		return fmt.Errorf("failed to get type and version: %w", err)
	}
	view.TypeAndVersion = typeAndVersion

	return nil
}

// Function to gather all discounts and organize them by subscriber and feedId
func (f *FeeManagerViewGenerator) gatherOrganizedDiscounts(ctx context.Context,
	params FeeManagerViewParams,
	view *FeeManagerView) (map[string]map[string]TokenDiscounts, error) {
	// Create filter options
	filterOpts := &bind.FilterOpts{
		Start:   params.FromBlock,
		End:     params.ToBlock,
		Context: ctx,
	}

	// Get all subscriber discount events
	iterator, err := f.contract.FilterSubscriberDiscountUpdated(filterOpts, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to filter subscriber discount events: %w", err)
	}
	defer iterator.Close()

	// Get references to token addresses for comparison
	callOpts := &bind.CallOpts{Context: ctx}

	type discountKey struct {
		subscriber common.Address
		feedID     [32]byte
		token      common.Address
	}

	// Find all combinations of subscriber, feedId, and token
	discountMap := make(map[string]discountKey)
	for iterator.Next() {
		event := iterator.GetEvent()

		feedIDStr := dsutil.HexEncodeBytes32(event.FeedId)

		// Create a unique key for this combination
		key := fmt.Sprintf("%s-%s-%s", event.Subscriber.Hex(), feedIDStr, event.Token.Hex())

		// Store the combination
		discountMap[key] = discountKey{
			subscriber: event.Subscriber,
			feedID:     event.FeedId,
			token:      event.Token,
		}
	}

	if err := iterator.Error(); err != nil {
		return nil, fmt.Errorf("error iterating through events: %w", err)
	}

	// Map[subscriberAddress][feedId]TokenDiscounts
	result := make(map[string]map[string]TokenDiscounts)

	for _, combo := range discountMap {
		subscriberAddr := combo.subscriber.Hex()
		feedIDHex := dsutil.HexEncodeBytes32(combo.feedID)

		// global discount is set using feedId of all zeros
		isGlobalDiscount := false
		var zeroBytes [32]byte
		if combo.feedID == zeroBytes {
			isGlobalDiscount = true
			feedIDHex = "global"
		}

		if result[subscriberAddr] == nil {
			result[subscriberAddr] = make(map[string]TokenDiscounts)
		}

		tokenDiscounts := result[subscriberAddr][feedIDHex]
		tokenDiscounts.IsGlobal = isGlobalDiscount

		var discount *big.Int
		var err error

		if isGlobalDiscount {
			discount, err = f.contract.SGlobalDiscounts(callOpts, combo.subscriber, combo.token)
		} else {
			discount, err = f.contract.SSubscriberDiscounts(callOpts, combo.subscriber, combo.feedID, combo.token)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to query discount: %w", err)
		}

		if combo.token.String() == view.LinkAddress {
			tokenDiscounts.Link = discount.String()
		} else if combo.token.String() == view.NativeAddress {
			tokenDiscounts.Native = discount.String()
		}

		result[subscriberAddr][feedIDHex] = tokenDiscounts
	}

	return result, nil
}
