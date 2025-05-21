package cre

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
)

func setupFakeDataProvider(testLogger zerolog.Logger, input *fake.Input, authKey string, expectedPrices map[string][]float64, priceIndexes map[string]*int) (string, error) {
	_, err := fake.NewFakeDataProvider(input)
	if err != nil {
		return "", errors.Wrap(err, "failed to set up fake data provider")
	}
	fakeAPIPath := "/fake/api/price"
	host := framework.HostDockerInternal()
	fakeFinalURL := fmt.Sprintf("%s:%d%s", host, input.Port, fakeAPIPath)

	getPriceResponseFn := func(feedID string) (map[string]interface{}, error) {
		testLogger.Info().Msgf("Preparing response for feedID: %s", feedID)
		priceIndex, ok := priceIndexes[feedID]
		if !ok {
			return nil, fmt.Errorf("no price index not found for feedID: %s", feedID)
		}

		expectedPrices, ok := expectedPrices[feedID]
		if !ok {
			return nil, fmt.Errorf("no expected prices not found for feedID: %s", feedID)
		}

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

		return response, nil
	}

	err = fake.Func("GET", fakeAPIPath, func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader != authKey {
			testLogger.Info().Msgf("Unauthorized request, expected auth key: %s actual auth key: %s", authKey, authHeader)
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}

		feedID := c.Query("feedID")
		if feedID == "" {
			testLogger.Info().Msgf("No feedID provided, returning error")
			c.JSON(400, gin.H{"error": "no feedID provided"})
			return
		}

		reponseBody, responseErr := getPriceResponseFn(feedID)
		if responseErr != nil {
			testLogger.Info().Msgf("Failed to get price response for feedID: %s, error: %s", feedID, responseErr)
			c.JSON(400, gin.H{"error": responseErr.Error()})
			return
		}

		c.JSON(200, reponseBody)
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
	NextPrice(feedID string, price *big.Int, elapsed time.Duration) bool
	ExpectedPrices(feedID string) []*big.Int
	ActualPrices(feedID string) []*big.Int
	AuthKey() string
}

// TrueUSDPriceProvider is a PriceProvider implementation that uses a live feed to get the price
type TrueUSDPriceProvider struct {
	testLogger   zerolog.Logger
	url          string
	actualPrices map[string][]*big.Int
}

func NewTrueUSDPriceProvider(testLogger zerolog.Logger, feedIDs []string) PriceProvider {
	pr := &TrueUSDPriceProvider{
		testLogger:   testLogger,
		url:          "https://api.real-time-reserves.verinumus.io/v1/chainlink/proof-of-reserves/TrueUSD",
		actualPrices: make(map[string][]*big.Int),
	}

	for _, feedID := range feedIDs {
		pr.actualPrices[feedID] = make([]*big.Int, 0)
	}

	return pr
}

func (l *TrueUSDPriceProvider) NextPrice(feedID string, price *big.Int, elapsed time.Duration) bool {
	cleanFeedID := cleanFeedID(feedID)
	// if price is nil or 0 it means that the feed hasn't been updated yet
	if price == nil || price.Cmp(big.NewInt(0)) == 0 {
		l.testLogger.Info().Msgf("Feed %s not updated yet, waiting for %s", cleanFeedID, elapsed)
		return true
	}

	l.testLogger.Info().Msgf("Feed %s updated after %s - price set, price=%s", cleanFeedID, elapsed, price)
	l.actualPrices[cleanFeedID] = append(l.actualPrices[cleanFeedID], price)

	// no other price to return, we are done
	return false
}

func (l *TrueUSDPriceProvider) URL() string {
	return l.url
}

func (l *TrueUSDPriceProvider) ExpectedPrices(feedID string) []*big.Int {
	// we don't have a way to check the price in the live feed, so we always assume it's correct
	// as long as it's != 0. And we only wait for the first price to be set.
	return l.actualPrices[cleanFeedID(feedID)]
}

func (l *TrueUSDPriceProvider) ActualPrices(feedID string) []*big.Int {
	// we don't have a way to check the price in the live feed, so we always assume it's correct
	// as long as it's != 0. And we only wait for the first price to be set.
	return l.actualPrices[cleanFeedID(feedID)]
}

func (l *TrueUSDPriceProvider) AuthKey() string {
	return ""
}

// FakePriceProvider is a PriceProvider implementation that uses a mocked feed to get the price
// It returns a configured price sequence and makes sure that the feed has been correctly updated
type FakePriceProvider struct {
	testLogger     zerolog.Logger
	priceIndex     map[string]*int
	url            string
	expectedPrices map[string][]*big.Int
	actualPrices   map[string][]*big.Int
	authKey        string
}

func cleanFeedID(feedID string) string {
	cleanFeedID := strings.TrimPrefix(feedID, "0x")
	if len(cleanFeedID) > 32 {
		cleanFeedID = cleanFeedID[:32]
	}
	return "0x" + cleanFeedID
}

func NewFakePriceProvider(testLogger zerolog.Logger, input *fake.Input, authKey string, feedIDs []string) (PriceProvider, error) {
	cleanFeedIDs := make([]string, 0, len(feedIDs))
	// workflow is sending feedIDs with 0x prefix and 32 bytes
	for _, feedID := range feedIDs {
		cleanFeedIDs = append(cleanFeedIDs, cleanFeedID(feedID))
	}
	priceIndexes := make(map[string]*int)
	for _, feedID := range cleanFeedIDs {
		priceIndexes[feedID] = ptr.Ptr(0)
	}

	expectedPrices := make(map[string][]*big.Int)
	pricesToServe := make(map[string][]float64)
	for _, feedID := range cleanFeedIDs {
		// Add more prices here as needed
		pricesFloat64 := []float64{math.Round((rand.Float64()*199+1)*100) / 100, math.Round((rand.Float64()*199+1)*100) / 100}
		pricesToServe[feedID] = pricesFloat64

		expectedPrices[feedID] = make([]*big.Int, len(pricesFloat64))
		for i, p := range pricesFloat64 {
			// convert float64 to big.Int by multiplying by 100
			// just like the PoR workflow does
			expectedPrices[feedID][i] = libc.Float64ToBigInt(p * 100)
		}
	}

	actualPrices := make(map[string][]*big.Int)
	for _, feedID := range cleanFeedIDs {
		actualPrices[feedID] = make([]*big.Int, 0)
	}

	url, err := setupFakeDataProvider(testLogger, input, authKey, pricesToServe, priceIndexes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set up fake data provider")
	}

	return &FakePriceProvider{
		testLogger:     testLogger,
		expectedPrices: expectedPrices,
		actualPrices:   actualPrices,
		priceIndex:     priceIndexes,
		url:            url,
		authKey:        authKey,
	}, nil
}

func (f *FakePriceProvider) priceAlreadyFound(feedID string, price *big.Int) bool {
	for _, p := range f.actualPrices[feedID] {
		if p.Cmp(price) == 0 {
			return true
		}
	}

	return false
}

func (f *FakePriceProvider) NextPrice(feedID string, price *big.Int, elapsed time.Duration) bool {
	cleanFeedID := cleanFeedID(feedID)
	// if price is nil or 0 it means that the feed hasn't been updated yet
	if price == nil || price.Cmp(big.NewInt(0)) == 0 {
		f.testLogger.Info().Msgf("Feed %s not updated yet, waiting for %s", cleanFeedID, elapsed)
		return true
	}

	if !f.priceAlreadyFound(cleanFeedID, price) {
		f.testLogger.Info().Msgf("Feed %s updated after %s - price set, price=%s", cleanFeedID, elapsed, price)
		f.actualPrices[cleanFeedID] = append(f.actualPrices[cleanFeedID], price)

		if len(f.actualPrices[cleanFeedID]) == len(f.expectedPrices[cleanFeedID]) {
			// all prices found, nothing more to check
			return false
		}

		if len(f.actualPrices[cleanFeedID]) > len(f.expectedPrices[cleanFeedID]) {
			panic("more prices found than expected")
		}
		f.testLogger.Info().Msgf("Changing price provider price for feed %s to %s", cleanFeedID, f.expectedPrices[cleanFeedID][len(f.actualPrices[cleanFeedID])].String())
		f.priceIndex[cleanFeedID] = ptr.Ptr(len(f.actualPrices[cleanFeedID]))

		// set new price and continue checking
		return true
	}

	// continue checking, price not updated yet
	return true
}

func (f *FakePriceProvider) ActualPrices(feedID string) []*big.Int {
	return f.actualPrices[cleanFeedID(feedID)]
}

func (f *FakePriceProvider) ExpectedPrices(feedID string) []*big.Int {
	return f.expectedPrices[cleanFeedID(feedID)]
}

func (f *FakePriceProvider) URL() string {
	return f.url
}

func (f *FakePriceProvider) AuthKey() string {
	return f.authKey
}
