package metering

import (
	"errors"
	"sync"

	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInvalidAmount       = errors.New("amount must be greater than 0")
)

type balanceStore struct {
	// Whether negative balances should return an error
	allowNegative bool
	// A balance of credits
	balance int64
	// Conversion rates of resource type to number of credits
	conversions map[string]decimal.Decimal
	lggr        logger.Logger
	mu          sync.RWMutex
}

type BalanceStore interface {
	Get() (balance int64)
	GetAs(unit string) (balance int64)
	Minus(amount int64) error
	MinusAs(unit string, amount int64) error
	Add(amount int64) error
	AddAs(unit string, amount int64) error
	AllowNegative()
}

var _ BalanceStore = (BalanceStore)(nil)

func NewBalanceStore(startingBalance int64, conversions map[string]decimal.Decimal, lggr logger.Logger) *balanceStore {
	// validations
	for resource, rate := range conversions {
		if rate.IsNegative() {
			// fail open
			lggr.Errorw("conversion rates must be a positive number, not using conversion", "resource", resource, "rate", rate)
			delete(conversions, resource)
		}
	}

	return &balanceStore{
		allowNegative: false,
		balance:       startingBalance,
		conversions:   conversions,
		lggr:          lggr,
	}
}

// convertToBalance converts a resource type amount to a credit amount
// This method should only be used under a read lock
func (bs *balanceStore) convertToBalance(fromUnit string, amount int64) (credits int64) {
	rate, ok := bs.conversions[fromUnit]
	if !ok {
		// Fail open, continue optimistically
		bs.lggr.Errorw("could not find conversion rate, continuing as 1:1", "unit", fromUnit)
		rate = decimal.NewFromInt(1)
	}
	return decimal.NewFromInt(amount).Mul(rate).RoundUp(0).IntPart()
}

// convertFromBalance converts a credit amount to a resource type amount
// This method should only be used under a read lock
func (bs *balanceStore) convertFromBalance(toUnit string, amount int64) (resources int64) {
	rate, ok := bs.conversions[toUnit]
	if !ok {
		// Fail open, continue optimistically
		bs.lggr.Errorw("could not find conversion rate, continuing as 1:1", "unit", toUnit)
		rate = decimal.NewFromInt(1)
	}
	return decimal.NewFromInt(amount).Div(rate).RoundUp(0).IntPart()
}

// Get returns the current credit balance
func (bs *balanceStore) Get() (balance int64) {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	return bs.balance
}

// GetAs returns the current universal credit balance expressed as a resource
func (bs *balanceStore) GetAs(unit string) (balance int64) {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	if bs.balance <= 0 {
		return 0
	}
	return bs.convertFromBalance(unit, bs.balance)
}

// Minus lowers the current credit balance
func (bs *balanceStore) Minus(amount int64) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if amount > bs.balance && !bs.allowNegative {
		return ErrInsufficientBalance
	}
	bs.balance -= amount
	return nil
}

// MinusAs lowers the current credit balance based on an amount of a type of resource
func (bs *balanceStore) MinusAs(unit string, amount int64) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	if amount <= 0 {
		return ErrInvalidAmount
	}
	balToMinus := bs.convertToBalance(unit, amount)
	if balToMinus > bs.balance && !bs.allowNegative {
		return ErrInsufficientBalance
	}
	bs.balance -= balToMinus
	return nil
}

// Add increases the current credit balance
func (bs *balanceStore) Add(amount int64) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	if amount <= 0 {
		return ErrInvalidAmount
	}
	bs.balance += amount
	return nil
}

// AddAs increases the current credit balance based on an amount of a type of resource
func (bs *balanceStore) AddAs(unit string, amount int64) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	if amount <= 0 {
		return ErrInvalidAmount
	}
	balToAdd := bs.convertToBalance(unit, amount)
	bs.balance += balToAdd
	return nil
}

// AllowNegative turns on the flag to allow negative balances
func (bs *balanceStore) AllowNegative() {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	bs.allowNegative = true
}
