package customendpoint

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/shopspring/decimal"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/core/bridges"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var _ types.ContractTransmitter = (*contractTracker)(nil)

var (
	promHTTPFetchTime = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "custom_endpoint_http_fetch_time",
		Help: "Time taken to fully execute the HTTP request",
	},
		[]string{"endpoint_target"},
	)
	promHTTPResponseBodySize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "custom_endpoint_http_response_body_size",
		Help: "Size (in bytes) of the HTTP response body",
	},
		[]string{"endpoint_target"},
	)
)

// Transmit uploads the result to customendpoint API endpoint, by making an HTTP call to
// the customendpoint ExternalAdapter. The EA does the actual signing and uploading work.
func (c *contractTracker) Transmit(
	ctx context.Context,
	reportCtx types.ReportContext,
	report types.Report,
	sigs []types.AttributedOnchainSignature,
) error {

	result, medianBigInt, err := c.getMedianFromReport(report)
	if err != nil {
		return err
	}
	c.transmittersWg.Add(1)

	// Don't block the current thread for transmitting results to the endpoint
	go c.doTransmit(ctx, result, medianBigInt, reportCtx.Epoch, reportCtx.Round)

	return nil
}

// Returns the latest epoch from the last stored transmission.
func (c *contractTracker) LatestConfigDigestAndEpoch(
	ctx context.Context,
) (
	configDigest types.ConfigDigest,
	epoch uint32,
	err error,
) {
	digest, err := c.digester.configDigest()
	c.ansLock.RLock()
	defer c.ansLock.RUnlock()
	return digest, c.storedAnswer.epoch, err
}

// TODO: Check if returning an item from StaticTransmitters value is good enough
func (c *contractTracker) FromAccount() types.Account {
	return StaticTransmitters[0]
}

func (c *contractTracker) doTransmit(
	ctx context.Context,
	result decimal.Decimal,
	medianBigInt *big.Int,
	epoch uint32,
	round uint8) {
	err := c.uploadResults(ctx, result)
	if err != nil {
		c.lggr.Warnw("Customendpoint Transmitter: Transmission failed",
			"EndpointName: ", c.digester.EndpointName,
			"EndpointTarget: ", c.digester.EndpointTarget,
			"PayloadType: ", c.digester.PayloadType,
			"Epoch: ", epoch,
			"Round: ", round,
			"Result: ", result.String(),
			"Error:", err)
	} else {
		c.ansLock.RLock()
		defer c.ansLock.RUnlock()

		// Skip saving the storedAnswer if a more recent saved storedAnswer already exists
		if epoch > c.storedAnswer.epoch || (epoch == c.storedAnswer.epoch && round > c.storedAnswer.round) {
			c.ansLock.RLock()
			defer c.ansLock.RUnlock()
			c.storedAnswer = answer{
				Data:      medianBigInt,
				Timestamp: c.clock.Now(),
				epoch:     epoch,
				round:     round,
			}
		}
	}
	c.transmittersWg.Done()
}

// Sends the result to the customendpoint External Adapter.
func (c *contractTracker) uploadResults(ctx context.Context, result decimal.Decimal) error {

	url, err := c.getBridgeURLFromName(c.digester.EndpointTarget)
	if err != nil {
		return err
	}
	priceFeed := fmt.Sprintf(c.bridgeRequestData, c.digester.PayloadType)

	var requestData = make(map[string]interface{})

	requestData[c.bridgeInputAtKey] = result
	mapParams, err := utils.UnmarshalToMap(priceFeed)
	if err != nil {
		return err
	}
	for k, v := range mapParams {
		requestData[k] = v
	}

	// URL is "safe" because it comes from the node's own database
	// Some node operators may run external adapters on their own hardware
	allowUnrestrictedNetworkAccess := true

	requestDataJSON, err := json.Marshal(requestData)
	if err != nil {
		return err
	}
	c.lggr.Debugw("Transmitter: sending request",
		"requestData", string(requestDataJSON),
		"url", url.String(),
	)

	requestCtx, cancel := context.WithTimeout(ctx, c.config.DefaultHTTPTimeout().Duration())
	defer cancel()

	responseBytes, _, _, elapsed, err := pipeline.MakeHTTPRequest(requestCtx, c.lggr, "POST", url, requestData, allowUnrestrictedNetworkAccess, c.config.DefaultHTTPLimit())
	if err != nil {
		return err
	}

	promHTTPFetchTime.WithLabelValues(url.String()).Set(float64(elapsed))
	promHTTPResponseBodySize.WithLabelValues(url.String()).Set(float64(len(responseBytes)))

	c.lggr.Debugw("Custom Endpoint Transmitter: fetched transmission response",
		"response", string(responseBytes),
		"url", url.String())
	return nil
}

func (c *contractTracker) getBridgeURLFromName(name string) (url.URL, error) {
	var bt bridges.BridgeType
	err := c.pipelineORM.GetQ().Get(&bt, "SELECT * FROM bridge_types WHERE name = $1", name)
	if err != nil {
		return url.URL{}, errors.Wrapf(err, "could not find bridge with name '%s'", name)
	}
	return url.URL(bt.URL), nil
}

// Gets the median from Report, and returns {median/multiplierUsed}, median, err.
// The multiplierUsed is the multiplier used in the OCR2 Job.
func (c *contractTracker) getMedianFromReport(report types.Report) (decimal.Decimal, *big.Int, error) {
	zero := big.NewInt(int64(0))
	median, err := c.reportCodec.MedianFromReport(report)
	if err != nil {
		return decimal.NewFromInt(0), zero, err
	}
	divideBy := big.NewInt(int64(c.multiplierUsed))

	if divideBy.Cmp(zero) == 0 {
		return decimal.NewFromInt(0), zero, errors.New("multiplierUsed cannot be 0")
	}
	var powerOf10 int32
	one := big.NewInt(int64(1))
	for divideBy.Cmp(one) != 0 {
		mod := big.NewInt(int64(0))
		divideBy, mod = divideBy.DivMod(divideBy, big.NewInt(int64(10)), mod)
		if mod.Cmp(zero) != 0 {
			return decimal.NewFromInt(0), zero,
				errors.New("multiplierUsed should only be a power of 10")
		}
		powerOf10--
	}
	return decimal.NewFromBigInt(median, powerOf10), median, err
}
