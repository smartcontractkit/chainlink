package metering

import (
	"errors"
	"fmt"
	"sync"

	"github.com/shopspring/decimal"
)

var (
	ErrInsufficientBalance  = errors.New("insufficient balance")
	ErrInvalidAmount        = errors.New("amount must be greater than 0")
	ErrResourceTypeNotFound = errors.New("could not find conversion rate, continuing as 1:1")
)

// balanceStore is a locked down interface to the in-execution credit balance.
// no state change details (like switching to metering mode) should be handled in it;
// rather consumers should consider errors core to business logic of metering/billing.
type balanceStore struct {
	// A balance of credits
	balance decimal.Decimal
	// Conversion rates of resource dimensions to number of units per credit
	conversions map[string]decimal.Decimal // TODO flip this
	// Total credits spent during execution
	spent decimal.Decimal
	mu    sync.RWMutex
}

func NewBalanceStore(
	startingBalance decimal.Decimal,
	conversions map[string]decimal.Decimal,
) (*balanceStore, error) {
	// validations
	for resource, rate := range conversions {
		if rate.IsNegative() {
			return nil, fmt.Errorf("conversion rate %s must be a positive number for resource %s", rate, resource)
		}
	}

	return &balanceStore{
		balance:     startingBalance,
		conversions: conversions,
	}, nil
}

// convertToBalance converts a resource dimension amount to a credit amount.
// This method should only be used under a read lock.
func (bs *balanceStore) convertToBalance(fromResourceType string, amount decimal.Decimal) (decimal.Decimal, error) {
	rate, ok := bs.conversions[fromResourceType]
	if !ok {
		return amount, ErrResourceTypeNotFound
	}

	// Special case for gas as gas token conversions are provided in amount per credit.
	// Other rates are provided as the inverse.
	if isGasSpendType(fromResourceType) {
		return amount.Div(rate).Round(defaultDecimalPrecision), nil
	}

	return amount.Mul(rate), nil
}

// ConvertToBalance converts a resource dimensions amount to a credit amount.
func (bs *balanceStore) ConvertToBalance(fromResourceType string, amount decimal.Decimal) (decimal.Decimal, error) {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	return bs.convertToBalance(fromResourceType, amount)
}

// convertFromBalance converts a credit amount to a resource dimensions amount.
// This method should only be used under a read lock.
func (bs *balanceStore) convertFromBalance(toResourceType string, amount decimal.Decimal) (decimal.Decimal, error) {
	rate, ok := bs.conversions[toResourceType]
	if !ok {
		return amount, ErrResourceTypeNotFound
	}

	// Special case for gas as gas token conversions are provided in amount per credit.
	// Other rates are provided as the inverse.
	if isGasSpendType(toResourceType) {
		return amount.Mul(rate).Round(0), nil
	}

	return amount.Div(rate), nil
}

// ConvertFromBalance converts a credit amount to a resource dimensions amount.
func (bs *balanceStore) ConvertFromBalance(toResourceType string, amount decimal.Decimal) (decimal.Decimal, error) {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	return bs.convertFromBalance(toResourceType, amount)
}

// Set sets the current balance to the provided amount and resets spend.
func (bs *balanceStore) Set(amount decimal.Decimal) {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	bs.balance = amount
	bs.spent = decimal.Zero
}

// Get returns the current credit balance
func (bs *balanceStore) Get() decimal.Decimal {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	return bs.balance
}

// GetAs returns the current universal credit balance expressed as a resource dimensions.
func (bs *balanceStore) GetAs(unit string) (decimal.Decimal, error) {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	return bs.convertFromBalance(unit, bs.balance)
}

// Minus lowers the current credit balance.
func (bs *balanceStore) Minus(amount decimal.Decimal) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if amount.LessThan(decimal.Zero) {
		return ErrInvalidAmount
	}

	if amount.GreaterThan(bs.balance) {
		return ErrInsufficientBalance
	}

	bs.balance = bs.balance.Sub(amount)

	return nil
}

// MinusAs lowers the current credit balance based on an amount of resource dimensions.
func (bs *balanceStore) MinusAs(resourceType string, amount decimal.Decimal) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if amount.LessThan(decimal.Zero) {
		return ErrInvalidAmount
	}

	balToMinus, err := bs.convertToBalance(resourceType, amount)
	if err != nil {
		return err
	}

	if balToMinus.GreaterThan(bs.balance) {
		return ErrInsufficientBalance
	}

	bs.balance = bs.balance.Sub(balToMinus)

	return nil
}

// Add increases the current credit balance.
func (bs *balanceStore) Add(amount decimal.Decimal) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if amount.LessThan(decimal.Zero) {
		return ErrInvalidAmount
	}

	bs.balance = bs.balance.Add(amount)

	return nil
}

// AddAs increases the current credit balance based on an amount of resource dimensions.
func (bs *balanceStore) AddAs(resourceType string, amount decimal.Decimal) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if amount.LessThan(decimal.Zero) {
		return ErrInvalidAmount
	}

	bal, err := bs.convertToBalance(resourceType, amount)
	if err != nil {
		return err
	}

	bs.balance = bs.balance.Add(bal)

	return nil
}

func (bs *balanceStore) AddSpent(amount decimal.Decimal) {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	bs.spent = bs.spent.Add(amount)
}

// GetSpent returns the total credits spent during execution.
// TODO: This should eventually be removed in favor of computing the spent amount from the metering report.
func (bs *balanceStore) GetSpent() decimal.Decimal {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	return bs.spent
}
