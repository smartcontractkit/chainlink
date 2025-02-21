package capabilities_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
)

func setupFakeDataProvider(testLogger zerolog.Logger, input *fake.Input, expectedPrices []float64, priceIndex *int) (string, error) {
	_, err := fake.NewFakeDataProvider(input)
	if err != nil {
		return "", errors.Wrap(err, "failed to set up fake data provider")
	}
	fakeAPIPath := "/fake/api/price"
	host := framework.HostDockerInternal()
	fakeFinalURL := fmt.Sprintf("%s:%d%s", host, input.Port, fakeAPIPath)

	getPriceResponseFn := func() map[string]interface{} {
		response := map[string]interface{}{
			"accountName": "TrueUSD",
			"totalTrust":  expectedPrices[*priceIndex],
			"ripcord":     false,
			"updatedAt":   time.Now().Format(time.RFC3339),
		}

		marshalled, mErr := json.Marshal(response)
		if mErr == nil {
			testLogger.Info().Msgf("Returning response: %s", string(marshalled))
		} else {
			testLogger.Info().Msgf("Returning response: %v", response)
		}

		return response
	}

	err = fake.Func("GET", fakeAPIPath, func(c *gin.Context) {
		c.JSON(200, getPriceResponseFn())
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to set up fake data provider")
	}

	return fakeFinalURL, nil
}

// PriceProvider abstracts away the logic of checking whether the feed has been correctly updated
// and it also returns port and URL of the price provider. This is so, because when using a mocked
// price provider we need start a separate service and whitelist its port and IP with the gateway job.
// Also, since it's a mocked price provider we can now check whether the feed has been correctly updated
// instead of only checking whether it has some price that's != 0.
type PriceProvider interface {
	URL() string
	NextPrice(price *big.Int, elapsed time.Duration) bool
	ExpectedPrices() []*big.Int
	ActualPrices() []*big.Int
}

// TrueUSDPriceProvider is a PriceProvider implementation that uses a live feed to get the price
type TrueUSDPriceProvider struct {
	testLogger   zerolog.Logger
	url          string
	actualPrices []*big.Int
}

func NewTrueUSDPriceProvider(testLogger zerolog.Logger) PriceProvider {
	return &TrueUSDPriceProvider{
		testLogger: testLogger,
		url:        "https://api.real-time-reserves.verinumus.io/v1/chainlink/proof-of-reserves/TrueUSD",
	}
}

func (l *TrueUSDPriceProvider) NextPrice(price *big.Int, elapsed time.Duration) bool {
	// if price is nil or 0 it means that the feed hasn't been updated yet
	if price == nil || price.Cmp(big.NewInt(0)) == 0 {
		return true
	}

	l.testLogger.Info().Msgf("Feed updated after %s - price set, price=%s", elapsed, price)
	l.actualPrices = append(l.actualPrices, price)

	// no other price to return, we are done
	return false
}

func (l *TrueUSDPriceProvider) URL() string {
	return l.url
}

func (l *TrueUSDPriceProvider) ExpectedPrices() []*big.Int {
	// we don't have a way to check the price in the live feed, so we always assume it's correct
	// as long as it's != 0. And we only wait for the first price to be set.
	return l.actualPrices
}

func (l *TrueUSDPriceProvider) ActualPrices() []*big.Int {
	// we don't have a way to check the price in the live feed, so we always assume it's correct
	// as long as it's != 0. And we only wait for the first price to be set.
	return l.actualPrices
}

// FakePriceProvider is a PriceProvider implementation that uses a mocked feed to get the price
// It returns a configured price sequence and makes sure that the feed has been correctly updated
type FakePriceProvider struct {
	testLogger     zerolog.Logger
	priceIndex     *int
	url            string
	expectedPrices []*big.Int
	actualPrices   []*big.Int
}

func NewFakePriceProvider(testLogger zerolog.Logger, input *fake.Input) (PriceProvider, error) {
	priceIndex := ptr.Ptr(0)
	// Add more prices here as needed
	expectedPricesFloat64 := []float64{182.9}
	expectedPrices := make([]*big.Int, len(expectedPricesFloat64))
	for i, p := range expectedPricesFloat64 {
		// convert float64 to big.Int by multiplying by 100
		// just like the PoR workflow does
		expectedPrices[i] = libc.Float64ToBigInt(p)
	}

	url, err := setupFakeDataProvider(testLogger, input, expectedPricesFloat64, priceIndex)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set up fake data provider")
	}

	return &FakePriceProvider{
		testLogger:     testLogger,
		expectedPrices: expectedPrices,
		priceIndex:     priceIndex,
		url:            url,
	}, nil
}

func (f *FakePriceProvider) priceAlreadyFound(price *big.Int) bool {
	for _, p := range f.actualPrices {
		if p.Cmp(price) == 0 {
			return true
		}
	}

	return false
}

func (f *FakePriceProvider) NextPrice(price *big.Int, elapsed time.Duration) bool {
	// if price is nil or 0 it means that the feed hasn't been updated yet
	if price == nil || price.Cmp(big.NewInt(0)) == 0 {
		return true
	}

	if !f.priceAlreadyFound(price) {
		f.testLogger.Info().Msgf("Feed updated after %s - price set, price=%s", elapsed, price)
		f.actualPrices = append(f.actualPrices, price)

		if len(f.actualPrices) == len(f.expectedPrices) {
			// all prices found, nothing more to check
			return false
		}

		if len(f.actualPrices) > len(f.expectedPrices) {
			panic("more prices found than expected")
		}
		f.testLogger.Info().Msgf("Changing price provider price to %s", f.expectedPrices[len(f.actualPrices)].String())
		*f.priceIndex = len(f.actualPrices)

		// set new price and continue checking
		return true
	}

	// continue checking, price not updated yet
	return true
}

func (f *FakePriceProvider) ActualPrices() []*big.Int {
	return f.actualPrices
}

func (f *FakePriceProvider) ExpectedPrices() []*big.Int {
	return f.expectedPrices
}

func (f *FakePriceProvider) URL() string {
	return f.url
}
