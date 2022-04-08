package customendpoint

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/shopspring/decimal"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/core/bridges"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var _ types.ContractTransmitter = (*ContractTracker)(nil)

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
func (c *ContractTracker) Transmit(
	ctx context.Context,
	reportCtx types.ReportContext,
	report types.Report,
	sigs []types.AttributedOnchainSignature,
) error {

	var result decimal.Decimal = decimal.NewFromInt(0)

	c.transmittersWg.Add(1)

	// Don't block the current thread for transmitting results to the endpoint
	go c.doTransmit(ctx, result, reportCtx.Epoch, reportCtx.Round)

	return nil
}

// Returns the latest epoch from the last stored transmission.
func (c *ContractTracker) LatestConfigDigestAndEpoch(
	ctx context.Context,
) (
	configDigest types.ConfigDigest,
	epoch uint32,
	err error,
) {
	digester, err := c.digester.configDigest()
	c.ansLock.RLock()
	defer c.ansLock.RUnlock()
	return digester, c.answer.epoch, err
}

// TODO: Check if returning an item from StaticTransmitters value is good enough
func (c *ContractTracker) FromAccount() types.Account {
	return StaticTransmitters[0]
}

func (c *ContractTracker) doTransmit(
	ctx context.Context, result decimal.Decimal, epoch uint32, round uint8) {
	err := c.uploadResults(ctx, result)
	if err != nil {
		c.lggr.Warnw("Customendpoint Transmitter: Transmission failed",
			"EndpointName: ", c.digester.EndpointName,
			"EndpointTarget: ", c.digester.EndpointTarget,
			"PayloadType: ", c.digester.PayloadType,
			"Epoch: ", epoch,
			"Round: ", round,
			"Error:", err)
	}

	c.ansLock.RLock()
	defer c.ansLock.RUnlock()

	// Skip saving the answer if a more recent saved answer already exists
	if epoch > c.answer.epoch || (epoch == c.answer.epoch && round > c.answer.round) {
		c.ansLock.RLock()
		defer c.ansLock.RUnlock()
		c.answer = Answer{
			Data:      result.BigInt(),
			Timestamp: time.Now(),
			epoch:     epoch,
			round:     round,
		}
	}
	c.transmittersWg.Done()
}

// Sends the result to the customendpoint External Adapter.
func (c *ContractTracker) uploadResults(ctx context.Context, result decimal.Decimal) error {

	url, err := c.getBridgeURLFromName(c.digester.EndpointTarget)
	if err != nil {
		return err
	}

	priceFeed := fmt.Sprintf(c.bridgeRequestData, c.digester.PayloadType)
	var requestData map[string]interface{}
	requestData[c.bridgeRequestData] = utils.MustUnmarshalToMap(priceFeed)
	requestData[c.bridgeInputAtKey] = result

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

	c.lggr.Debugw("Custom Endpoint Transmitter: fetched answer",
		"answer", string(responseBytes),
		"url", url.String())
	return nil
}

func (c *ContractTracker) getBridgeURLFromName(name string) (url.URL, error) {
	var bt bridges.BridgeType
	err := c.pipelineORM.GetQ().Get(&bt, "SELECT * FROM bridge_types WHERE name = $1", name)
	if err != nil {
		return url.URL{}, errors.Wrapf(err, "could not find bridge with name '%s'", name)
	}
	return url.URL(bt.URL), nil
}
