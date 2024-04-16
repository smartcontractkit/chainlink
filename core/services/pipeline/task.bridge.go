package pipeline

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline/internal/eautils"
)

// NOTE: These metrics generate a new label per bridge, this should be safe
// since the number of bridges is almost always relatively small (<< 1000)
//
// We already have promHTTPFetchTime but the bridge-specific gauges allow for
// more granular metrics
var (
	promBridgeLatency = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "bridge_latency_seconds",
		Help: "Bridge latency in seconds scoped by name",
	},
		[]string{"name"},
	)
	promBridgeErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "bridge_errors_total",
		Help: "Bridge error count scoped by name",
	},
		[]string{"name"},
	)
	promBridgeCacheHits = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "bridge_cache_hits_total",
		Help: "Bridge cache hits count scoped by name",
	},
		[]string{"name"},
	)
	promBridgeCacheErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "bridge_cache_errors_total",
		Help: "Bridge cache errors count scoped by name",
	},
		[]string{"name"},
	)
)

// Return types:
//
//	string
type BridgeTask struct {
	BaseTask `mapstructure:",squash"`

	Name              string `json:"name"`
	RequestData       string `json:"requestData"`
	IncludeInputAtKey string `json:"includeInputAtKey"`
	Async             string `json:"async"`
	CacheTTL          string `json:"cacheTTL"`
	Headers           string `json:"headers"`

	specId       int32
	orm          bridges.ORM
	config       Config
	bridgeConfig BridgeConfig
	httpClient   *http.Client
}

var _ Task = (*BridgeTask)(nil)

var zeroURL = new(url.URL)

const stalenessCap = 30 * time.Minute

func (t *BridgeTask) Type() TaskType {
	return TaskTypeBridge
}

func (t *BridgeTask) Run(ctx context.Context, lggr logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	inputValues, err := CheckInputs(inputs, -1, -1, 0)
	if err != nil {
		return Result{Error: errors.Wrap(err, "task inputs")}, runInfo
	}

	var (
		name              StringParam
		requestData       MapParam
		includeInputAtKey StringParam
		cacheTTL          Uint64Param
		reqHeaders        StringSliceParam
	)
	err = multierr.Combine(
		errors.Wrap(ResolveParam(&name, From(NonemptyString(t.Name))), "name"),
		errors.Wrap(ResolveParam(&requestData, From(VarExpr(t.RequestData, vars), JSONWithVarExprs(t.RequestData, vars, false), nil)), "requestData"),
		errors.Wrap(ResolveParam(&includeInputAtKey, From(t.IncludeInputAtKey)), "includeInputAtKey"),
		errors.Wrap(ResolveParam(&cacheTTL, From(ValidDurationInSeconds(t.CacheTTL), t.bridgeConfig.BridgeCacheTTL().Seconds())), "cacheTTL"),
		errors.Wrap(ResolveParam(&reqHeaders, From(NonemptyString(t.Headers), "[]")), "reqHeaders"),
	)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	if len(reqHeaders)%2 != 0 {
		return Result{Error: errors.Errorf("headers must have an even number of elements")}, runInfo
	}

	url, err := t.getBridgeURLFromName(ctx, name)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	var metaMap MapParam

	meta, _ := vars.Get("jobRun.meta")
	switch v := meta.(type) {
	case map[string]interface{}:
		metaMap = MapParam(v)
	case nil:
	default:
		lggr.Warnw(`"meta" field on task run is malformed, discarding`,
			"task", t.DotID(),
			"meta", meta,
		)
	}

	requestData = withRunInfo(requestData, metaMap)
	if t.IncludeInputAtKey != "" {
		if len(inputValues) > 0 {
			requestData[string(includeInputAtKey)] = inputValues[0]
		}
	}

	if t.Async == "true" {
		responseURL := t.bridgeConfig.BridgeResponseURL()
		if responseURL != nil && *responseURL != *zeroURL {
			responseURL.Path = path.Join(responseURL.Path, "/v2/resume/", t.uuid.String())
		}
		var s string
		if responseURL != nil {
			s = responseURL.String()
		}
		requestData["responseURL"] = s
	}

	requestDataJSON, err := json.Marshal(requestData)
	if err != nil {
		return Result{Error: err}, runInfo
	}
	lggr.Tracew("Bridge task: sending request",
		"requestData", string(requestDataJSON),
		"url", url.String(),
	)

	requestCtx, cancel := httpRequestCtx(ctx, t, t.config)
	defer cancel()

	// cacheTTL should not exceed stalenessCap.
	cacheDuration := time.Duration(cacheTTL) * time.Second
	if cacheDuration > stalenessCap {
		lggr.Warnf("bridge task cacheTTL exceeds stalenessCap %s, overriding value to stalenessCap", stalenessCap)
		cacheDuration = stalenessCap
	}

	var cachedResponse bool
	responseBytes, statusCode, headers, elapsed, err := makeHTTPRequest(requestCtx, lggr, "POST", url, reqHeaders, requestData, t.httpClient, t.config.DefaultHTTPLimit())

	// check for external adapter response object status
	if code, ok := eautils.BestEffortExtractEAStatus(responseBytes); ok {
		statusCode = code
	}

	if err != nil || statusCode != http.StatusOK {
		promBridgeErrors.WithLabelValues(t.Name).Inc()
		if cacheTTL == 0 {
			return Result{Error: err}, RunInfo{IsRetryable: isRetryableHTTPError(statusCode, err)}
		}

		var cacheErr error
		responseBytes, cacheErr = t.orm.GetCachedResponse(ctx, t.dotID, t.specId, cacheDuration)
		if cacheErr != nil {
			promBridgeCacheErrors.WithLabelValues(t.Name).Inc()
			if !errors.Is(cacheErr, sql.ErrNoRows) {
				lggr.Warnw("Bridge task: cache fallback failed",
					"err", cacheErr.Error(),
					"url", url.String(),
				)
			}
			return Result{Error: err}, RunInfo{IsRetryable: isRetryableHTTPError(statusCode, err)}
		}
		promBridgeCacheHits.WithLabelValues(t.Name).Inc()
		lggr.Debugw("Bridge task: request failed, falling back to cache",
			"response", string(responseBytes),
			"url", url.String(),
		)
		cachedResponse = true
	} else {
		promBridgeLatency.WithLabelValues(t.Name).Set(elapsed.Seconds())
	}

	if t.Async == "true" {
		// Look for a `pending` flag. This check is case-insensitive because http.Header normalizes header names
		if _, ok := headers["X-Chainlink-Pending"]; ok {
			return result, pendingRunInfo()
		}

		var response struct {
			Pending bool `json:"pending"`
		}
		if err := json.Unmarshal(responseBytes, &response); err == nil && response.Pending {
			return Result{}, pendingRunInfo()
		}
	}

	if !cachedResponse && cacheTTL > 0 {
		err := t.orm.UpsertBridgeResponse(ctx, t.dotID, t.specId, responseBytes)
		if err != nil {
			lggr.Errorw("Bridge task: failed to upsert response in bridge cache", "err", err)
		}
	}

	// NOTE: We always stringify the response since this is required for all current jobs.
	// If a binary response is required we might consider adding an adapter
	// flag such as  "BinaryMode: true" which passes through raw binary as the
	// value instead.
	result = Result{Value: string(responseBytes)}

	promHTTPFetchTime.WithLabelValues(t.DotID()).Set(float64(elapsed))
	promHTTPResponseBodySize.WithLabelValues(t.DotID()).Set(float64(len(responseBytes)))

	lggr.Tracew("Bridge task: fetched answer",
		"answer", result.Value,
		"url", url.String(),
		"dotID", t.DotID(),
		"cached", cachedResponse,
	)
	return result, runInfo
}

func (t BridgeTask) getBridgeURLFromName(ctx context.Context, name StringParam) (URLParam, error) {
	bt, err := t.orm.FindBridge(ctx, bridges.BridgeName(name))
	if err != nil {
		return URLParam{}, errors.Wrapf(err, "could not find bridge with name '%s'", name)
	}
	return URLParam(bt.URL), nil
}

func withRunInfo(request MapParam, meta MapParam) MapParam {
	output := make(MapParam)
	for k, v := range request {
		output[k] = v
	}
	if meta != nil {
		output["meta"] = meta
	}
	return output
}
