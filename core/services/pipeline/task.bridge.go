package pipeline

import (
	"context"
	"database/sql"
	stderrors "errors"
	"maps"
	"math"
	"math/big"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline/eautils"
)

// NOTE: These metrics generate a new label per bridge, this should be safe
// since the number of bridges is almost always relatively small (<< 1000)
//
// We already have promHTTPFetchTime but the bridge-specific gauges allow for
// more granular metrics
var (
	promBridgeLatency = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "bridge_latency_seconds",
		Help: "Bridge latency in seconds scoped by name and response status code",
	},
		[]string{"name", "status_code_group"},
	)
	promBridgeLatencyHist = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "bridge_latency_histogram_ms",
		Help: "Bridge latency histogram in milliseconds scoped by name and response status code",
		Buckets: []float64{
			25, 50, 100, 250, 500,
		},
	},
		[]string{"name", "status_code_group"},
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
	// CheckRequired when "true" enables validation that the HTTP response JSON
	// contains paths required by strict downstream jsonparse tasks (see
	// requiredJSONPaths). When empty or "false", that check is skipped.
	CheckRequired string `json:"checkRequired"`

	specId       int32
	orm          bridges.ORM
	config       Config
	bridgeConfig BridgeConfig
	httpClient   *http.Client

	// requiredJSONPaths is populated in runner.InitializePipeline from strict
	// downstream jsonparse tasks. When CheckRequired is true and cacheTTL is set,
	// validation uses these paths to fall back to cache when the live response
	// omits required keys.
	requiredJSONPaths [][]string
}

// bridgeHTTPOutcome holds response state after the HTTP round-trip through EA JSON status,
// optional required-path validation, and optional cache fallback.
type bridgeHTTPOutcome struct {
	body           []byte
	statusCode     int
	err            error
	cachedResponse bool
}

type BridgeTelemetry struct {
	RequestStartTimestamp  time.Time `json:"requestStartTimestamp"`
	RequestFinishTimestamp time.Time `json:"requestFinishTimestamp"`
	RequestData            []byte    `json:"requestData"`
	ResponseData           []byte    `json:"responseData"`
	Name                   string    `json:"name"`
	DotID                  string    `json:"dotID"`
	ResponseError          *string   `json:"responseError"`
	StreamID               *uint32   `json:"streamID"`
	SpecID                 int32     `json:"specID"`
	ResponseStatusCode     int       `json:"responseStatusCode"`
	LocalCacheHit          bool      `json:"localCacheHit"`
}

var _ Task = (*BridgeTask)(nil)

var zeroURL = new(url.URL)

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
		checkRequired     BoolParam
	)
	err = stderrors.Join(
		errors.Wrap(ResolveParam(&name, From(NonemptyString(t.Name))), "name"),
		errors.Wrap(ResolveParam(&requestData, From(VarExpr(t.RequestData, vars), JSONWithVarExprs(t.RequestData, vars, false), nil)), "requestData"),
		errors.Wrap(ResolveParam(&includeInputAtKey, From(t.IncludeInputAtKey)), "includeInputAtKey"),
		errors.Wrap(ResolveParam(&cacheTTL, From(ValidDurationInSeconds(t.CacheTTL), t.bridgeConfig.BridgeCacheTTL().Seconds())), "cacheTTL"),
		errors.Wrap(ResolveParam(&reqHeaders, From(NonemptyString(t.Headers), "[]")), "reqHeaders"),
		errors.Wrap(ResolveParam(&checkRequired, From(NonemptyString(t.CheckRequired), false)), "checkRequired"),
	)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	if len(reqHeaders)%2 != 0 {
		return Result{Error: errors.Errorf("headers must have an even number of elements")}, runInfo
	}

	overtimeCtx, cancel := overtimeContext(ctx)
	defer cancel()

	url, err := t.getBridgeURLFromName(overtimeCtx, name)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	requestDataJSON, err := t.finalizeAndMarshalBridgeRequestData(lggr, vars, inputValues, &requestData, includeInputAtKey)
	if err != nil {
		return Result{Error: err}, runInfo
	}
	logger.Sugared(lggr).Tracew("Bridge task: sending request",
		"requestData", string(requestDataJSON),
		"url", url.String(),
	)

	requestCtx, cancel := httpRequestCtx(ctx, t, t.config)
	defer cancel()

	var cachedResponse bool
	responseBytes, statusCode, headers, start, finish, err := makeHTTPRequest(requestCtx, lggr, "POST", url, reqHeaders, requestData, t.httpClient, t.config.DefaultHTTPLimit())
	elapsed := finish.Sub(start)
	promBridgeLatency.WithLabelValues(t.Name, statusCodeGroup(statusCode)).Set(elapsed.Seconds())
	promBridgeLatencyHist.WithLabelValues(t.Name, statusCodeGroup(statusCode)).Observe(float64(elapsed.Milliseconds()))

	out := bridgeHTTPOutcome{
		body:           responseBytes,
		statusCode:     statusCode,
		err:            err,
		cachedResponse: false,
	}

	defer func() {
		// Runs when Run returns; reads the final err, responseBytes, statusCode, and cachedResponse
		// after EA JSON handling, required-path validation, and optional cache fallback.
		telemetryCh := GetTelemetryCh(ctx)
		if telemetryCh != nil {
			bt := &BridgeTelemetry{
				Name:                   t.Name,
				RequestData:            requestDataJSON,
				ResponseData:           responseBytes,
				ResponseStatusCode:     statusCode,
				RequestStartTimestamp:  start,
				RequestFinishTimestamp: finish,
				LocalCacheHit:          cachedResponse,
				SpecID:                 t.specId,
				DotID:                  t.DotID(),
			}
			if err != nil {
				bt.ResponseError = new(string)
				*bt.ResponseError = err.Error()
			}

			bt.resolveStreamID(t, vars, lggr)

			select {
			case telemetryCh <- bt:
			default:
				lggr.Warn("bridge task: telemetry channel is full, dropping telemetry")
			}
		}
	}()

	out.statusCode = eaJSONResponseStatus(out.body, out.statusCode)
	liveOK := out.err == nil && out.statusCode == http.StatusOK
	out.err = t.maybeValidateRequiredJSONPaths(out.err, liveOK, cacheTTL, checkRequired, out.body)

	out, earlyResult, earlyRunInfo := t.resolveFailureOrCache(overtimeCtx, lggr, url, out, cacheTTL)
	if earlyResult != nil {
		return *earlyResult, *earlyRunInfo
	}

	responseBytes = out.body
	statusCode = out.statusCode
	err = out.err
	cachedResponse = out.cachedResponse

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
		err := t.orm.UpsertBridgeResponse(overtimeCtx, t.dotID, t.specId, responseBytes)
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

	logger.Sugared(lggr).Tracew("Bridge task: fetched answer",
		"answer", result.Value,
		"url", url.String(),
		"dotID", t.DotID(),
		"cached", cachedResponse,
	)
	return result, runInfo
}

// finalizeAndMarshalBridgeRequestData merges job meta, upstream inputs, and async resume URL into requestData,
// writes the merged map back through requestData for use by makeHTTPRequest, and returns the JSON body for logging
// and telemetry.
func (t *BridgeTask) finalizeAndMarshalBridgeRequestData(lggr logger.Logger, vars Vars, inputValues []any, requestData *MapParam, includeInputAtKey StringParam) ([]byte, error) {
	var metaMap MapParam

	meta, _ := vars.Get("jobRun.meta")
	switch v := meta.(type) {
	case map[string]any:
		metaMap = MapParam(v)
	case nil:
	default:
		lggr.Warnw(`"meta" field on task run is malformed, discarding`,
			"task", t.DotID(),
			"meta", meta,
		)
	}

	merged := withRunInfo(*requestData, metaMap)
	if t.IncludeInputAtKey != "" {
		if len(inputValues) > 0 {
			merged[string(includeInputAtKey)] = inputValues[0]
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
		merged["responseURL"] = s
	}

	*requestData = merged
	return json.Marshal(merged)
}

// eaJSONResponseStatus returns the external-adapter status from the response body when present, otherwise the
// HTTP status from the round-trip.
func eaJSONResponseStatus(body []byte, httpStatus int) int {
	if code, ok := eautils.BestEffortExtractEAStatus(body); ok {
		return code
	}
	return httpStatus
}

// maybeValidateRequiredJSONPaths runs jsonDecodeValidateRequiredPaths when the live call succeeded (HTTP 200, no
// transport error), checkRequired is enabled, cache is enabled, and the pipeline registered required paths. A
// failed check sets err so the caller can fall back to the bridge cache when configured.
func (t *BridgeTask) maybeValidateRequiredJSONPaths(err error, liveOK bool, cacheTTL Uint64Param, checkRequired BoolParam, responseBytes []byte) error {
	if err != nil || !liveOK || cacheTTL == 0 || !bool(checkRequired) || len(t.requiredJSONPaths) == 0 {
		return err
	}
	if verr := jsonDecodeValidateRequiredPaths(responseBytes, t.requiredJSONPaths); verr != nil {
		return errors.Wrap(verr, "bridge response failed required JSON path check for downstream jsonparse")
	}
	return err
}

// resolveFailureOrCache handles a non-success HTTP outcome: it prefers the EA error from the body over the
// transport error, increments error metrics, then either returns immediately (no cache TTL or cache miss) or
// replaces the response body from the bridge cache. Non-nil early Result and RunInfo mean the caller must return
// without continuing the success path.
func (t *BridgeTask) resolveFailureOrCache(
	ctx context.Context,
	lggr logger.Logger,
	url URLParam,
	out bridgeHTTPOutcome,
	cacheTTL Uint64Param,
) (bridgeHTTPOutcome, *Result, *RunInfo) {
	if out.err == nil && out.statusCode == http.StatusOK {
		return out, nil, nil
	}
	if adapterErr := eautils.BestEffortExtractEAError(out.body); adapterErr != nil {
		out.err = adapterErr
	}

	promBridgeErrors.WithLabelValues(t.Name).Inc()
	if cacheTTL == 0 {
		lggr.Debugw("Bridge task: request failed",
			"response", string(out.body),
			"url", url.String(),
			"status_code", out.statusCode,
			"error", out.err,
		)
		retry := RunInfo{IsRetryable: isRetryableHTTPError(out.statusCode, out.err)}
		return out, &Result{Error: out.err}, &retry
	}

	//nolint:gosec // disable G115
	cachedBytes, cacheErr := t.orm.GetCachedResponse(ctx, t.dotID, t.specId, time.Duration(cacheTTL)*time.Second)
	if cacheErr != nil {
		promBridgeCacheErrors.WithLabelValues(t.Name).Inc()
		if !errors.Is(cacheErr, sql.ErrNoRows) {
			lggr.Warnw("Bridge task: cache fallback failed",
				"err", cacheErr.Error(),
				"url", url.String(),
			)
		}
		retry := RunInfo{IsRetryable: isRetryableHTTPError(out.statusCode, out.err)}
		return out, &Result{Error: out.err}, &retry
	}
	promBridgeCacheHits.WithLabelValues(t.Name).Inc()
	lggr.Debugw("Bridge task: request failed, falling back to cache",
		"response", string(cachedBytes),
		"url", url.String(),
	)
	out.body = cachedBytes
	out.cachedResponse = true
	return out, nil, nil
}

func (bt *BridgeTelemetry) resolveStreamID(t *BridgeTask, vars Vars, lggr logger.Logger) {
	if t.StreamID.Valid {
		bt.StreamID = &t.StreamID.Uint32
	} else {
		if streamID, sErr := vars.Get("jb.streamID"); sErr == nil {
			if streamIDptr, ok := streamID.(*uint32); !ok {
				lggr.Debugw("Bridge task: streamID from vars is not a *uint32", "streamID", streamID)
			} else {
				bt.StreamID = streamIDptr
			}
		} else {
			lggr.Debugw("Bridge task: failed to get streamID from vars", "err", sErr)
		}
	}
}

func (t *BridgeTask) getBridgeURLFromName(ctx context.Context, name StringParam) (URLParam, error) {
	bt, err := t.orm.FindBridge(ctx, bridges.BridgeName(name))
	if err != nil {
		return URLParam{}, errors.Wrapf(err, "could not find bridge with name '%s'", name)
	}
	return URLParam(bt.URL), nil
}

func withRunInfo(request MapParam, meta MapParam) MapParam {
	output := make(MapParam)
	maps.Copy(output, request)
	if meta != nil {
		output["meta"] = meta
	}
	return output
}

// getRequiredJSONPaths returns JSON path segments (split the same way as jsonparse) that strict
// downstream jsonparse tasks require when they read this bridge task's string output directly.
//
// Limitations: does not follow merge/median/etc.; skips lax jsonparse, dynamic path/data
// containing "$(", data="$(bridge.field)", or jsonparse whose data input is not this bridge
// (by lowest output index or explicit $(bridgeID)).
func (t *BridgeTask) getRequiredJSONPaths() [][]string {
	if t == nil {
		return nil
	}
	bridgeID := t.ID()
	bridgeDot := t.DotID()

	seen := make(map[string]struct{})
	var out [][]string

	for _, d := range t.GetDescendantTasks() {
		jp, ok := d.(*JSONParseTask)
		if !ok {
			continue
		}
		if jsonParseLaxStatic(jp) {
			continue
		}
		if strings.Contains(jp.Path, "$(") {
			continue
		}
		dataTrim := strings.TrimSpace(jp.Data)
		if strings.Contains(dataTrim, "$(") {
			m := bridgeDataRootVarRegexp.FindStringSubmatch(dataTrim)
			if len(m) != 2 || m[1] != bridgeDot {
				continue
			}
		}
		if !jsonParseDataReferencesBridge(jp, bridgeID, bridgeDot) {
			continue
		}
		segs, ok := jsonParseStaticPathSegments(jp)
		if !ok {
			continue
		}
		key := strings.Join(segs, "\x00")
		if _, dup := seen[key]; dup {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, segs)
	}
	return out
}

// bridgeDataRootVarRegexp matches a data field that is only a single $(dotID)
// reference with no nested keypath (dotID: [a-zA-Z0-9_]+).
var bridgeDataRootVarRegexp = regexp.MustCompile(`^\$\(\s*([a-zA-Z0-9_]+)\s*\)$`)

func jsonParseLaxStatic(j *JSONParseTask) bool {
	trimmed := strings.TrimSpace(j.Lax)
	if trimmed == "" {
		return false
	}
	b, err := strconv.ParseBool(trimmed)
	return err == nil && b
}

func jsonParseDataReferencesBridge(j *JSONParseTask, bridgeID int, bridgeDot string) bool {
	dataTrim := strings.TrimSpace(j.Data)
	if dataTrim == "" {
		src, ok := jsonParseLowestIndexPropagatingInput(j)
		return ok && src.ID() == bridgeID
	}
	m := bridgeDataRootVarRegexp.FindStringSubmatch(dataTrim)
	return len(m) == 2 && m[1] == bridgeDot
}

func jsonParseLowestIndexPropagatingInput(j *JSONParseTask) (Task, bool) {
	var (
		found        bool
		minOutputIdx int32 = math.MaxInt32
		src          Task
	)
	for _, dep := range j.Inputs() {
		if !dep.PropagateResult {
			continue
		}
		idx := dep.InputTask.OutputIndex()
		if !found || idx < minOutputIdx {
			found = true
			minOutputIdx = idx
			src = dep.InputTask
		}
	}
	return src, found
}

func jsonParseStaticPathSegments(j *JSONParseTask) ([]string, bool) {
	if strings.TrimSpace(j.Path) == "" {
		return nil, false
	}
	sep := strings.TrimSpace(j.Separator)
	if sep == "" {
		sep = ","
	}
	parts := strings.Split(j.Path, sep)
	if len(parts) == 0 {
		return nil, false
	}
	return parts, true
}

// jsonDecodeValidateRequiredPaths checks each path in the raw JSON body using
// streaming lookup (no full value unmarshal). Missing keys, JSON null, JSON numbers
// that are not valid for shopspring/decimal (including exponent/size limits), and the
// string "NaN" are treated as failure so bridge cache fallback can apply.
//
// For each present value, parseRequiredValidValue must succeed: decimals and floats
// via shopspring/decimal, integer literals (decimal or hex with optional sign) via
// math/big so huge 0x… values validate and scientific notation is not mistaken for
// hexadecimal.
func jsonDecodeValidateRequiredPaths(body []byte, paths [][]string) error {
	if len(paths) == 0 {
		return nil
	}
	for _, path := range paths {
		if len(path) == 0 {
			continue
		}

		value, dataType, _, err := jsonparser.Get(body, path...)
		if err != nil {
			return errors.Wrapf(err, "required path %q", strings.Join(path, ","))
		}

		if dataType == jsonparser.Null {
			return errors.Errorf("required path %q is null", strings.Join(path, ","))
		}

		if err := parseRequiredValidValue(value); err != nil {
			return errors.Wrapf(err, "required path %q", strings.Join(path, ","))
		}
	}
	return nil
}

// parseRequiredValidValue accepts values that can be parsed by shopspring/decimal or math/big.
func parseRequiredValidValue(value []byte) error {
	strValue := string(value)
	if _, err := decimal.NewFromString(strValue); err == nil {
		return nil
	}
	if _, ok := new(big.Int).SetString(strValue, 0); ok {
		return nil
	}
	return errors.Errorf("invalid value: %s", string(value))
}
