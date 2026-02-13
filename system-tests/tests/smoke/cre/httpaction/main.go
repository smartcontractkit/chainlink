//go:build wasip1

package main

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/httpaction/config"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
	"google.golang.org/protobuf/types/known/durationpb"
	"gopkg.in/yaml.v3"
)

const bothSetRegressionSuccess = "HTTP Action multi-headers regression completed"

// Expected Set-Cookie values from the fake server (v2_http_action_test.go).
var expectedSetCookieSubstrings = []string{
	"sessionid=multi-e2e-1",
	"csrf=multi-e2e-2",
	"pref=multi-e2e-3",
}

// bothSetRegressionExpectedSubstrings must all appear in the capability's user error
// when both Headers and MultiHeaders are set on a request (validation rejects).
var bothSetRegressionExpectedSubstrings = []string{
	"Headers or MultiHeaders",
	"not both",
}

// nodeTestFunc is the type of a test run inside node mode (cfg, runtime, client, logger) -> (result, err).
// Logger is the CRE runtime logger; see cre.Runtime.Logger() which returns *slog.Logger.
type nodeTestFunc func(config.Config, cre.NodeRuntime, *http.Client, *slog.Logger) (string, error)

// testCaseHandlers dispatch by cfg.TestCase; special cases run first, default runs generic CRUD.
var testCaseHandlers = map[string]nodeTestFunc{
	"multi-headers":      runMultiHeadersTest,
	"mh-regression-both": runBothSetRegressionTest,
}

// --- Entry and workflow ---

func main() {
	wasm.NewRunner(func(b []byte) (config.Config, error) {
		wfCfg := config.Config{}
		if err := yaml.Unmarshal(b, &wfCfg); err != nil {
			return config.Config{}, fmt.Errorf("error unmarshalling config: %w", err)
		}
		return wfCfg, nil
	}).Run(RunHTTPActionSuccessWorkflow)
}

func RunHTTPActionSuccessWorkflow(wfCfg config.Config, _ *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[config.Config], error) {
	return cre.Workflow[config.Config]{
		cre.Handler(
			cron.Trigger(&cron.Config{Schedule: "*/30 * * * * *"}),
			onCronTrigger,
		),
	}, nil
}

func onCronTrigger(wfCfg config.Config, runtime cre.Runtime, payload *cron.Payload) (any, error) {
	logger := runtime.Logger()
	logger.Info(
		"HTTP Action workflow triggered",
		"testCase",
		wfCfg.TestCase,
		"method",
		wfCfg.Method,
		"url",
		wfCfg.URL,
	)

	return runCRUDSuccessTest(wfCfg, runtime)
}

// --- Test orchestration and test functions ---

func runCRUDSuccessTest(wfCfg config.Config, runtime cre.Runtime) (string, error) {
	logger := runtime.Logger()
	logger.Info("Running HTTP Action capability", "testCase", wfCfg.TestCase, "method", wfCfg.Method, "url", wfCfg.URL)

	crudPromise := cre.RunInNodeMode(wfCfg, runtime,
		func(cfg config.Config, nodeRuntime cre.NodeRuntime) (string, error) {
			client := &http.Client{}
			if fn, ok := testCaseHandlers[cfg.TestCase]; ok {
				return fn(cfg, nodeRuntime, client, logger)
			}
			return runDefaultCRUDTest(cfg, nodeRuntime, client, logger)
		},
		cre.ConsensusIdenticalAggregation[string](),
	)

	result, err := crudPromise.Await()
	if err != nil {
		logger.Error("Failed to complete HTTP Action capability", "error", err)
		return "", fmt.Errorf("HTTP Action test failed: %w", err)
	}
	logger.Info("HTTP Action test completed", "result", result)
	return result, nil
}

func runDefaultCRUDTest(cfg config.Config, nodeRuntime cre.NodeRuntime, client *http.Client, logger *slog.Logger) (string, error) {
	req := &http.Request{
		Url:     cfg.URL,
		Method:  cfg.Method,
		Headers: map[string]string{"Content-Type": "application/json"},
		Body:    []byte(cfg.Body),
		Timeout: &durationpb.Duration{Seconds: 10},
	}
	logger.Info("Testing HTTP Action capability with configuration", "url", req.Url, "method", req.Method, "hasBody", len(cfg.Body) > 0)

	resp, err := client.SendRequest(nodeRuntime, req).Await()
	if err != nil {
		logger.Error("Failed to complete HTTP Action request", "error", err, "url", req.Url, "method", req.Method)
		return "", fmt.Errorf("HTTP Action %s request failed: %w", req.Method, err)
	}
	if !statusOK(resp.StatusCode) {
		logger.Error("Failed response to HTTP Action request", "status", resp.StatusCode, "url", req.Url, "method", req.Method)
		return "", fmt.Errorf("HTTP Action %s request failed with status: %d", req.Method, resp.StatusCode)
	}
	logger.Info("HTTP Action completed", "url", req.Url, "method", req.Method, "status", resp.StatusCode, "body", string(resp.Body))
	return fmt.Sprintf("HTTP Action CRUD success test completed: %s", cfg.TestCase), nil
}

// runMultiHeadersTest sends two requests (Headers-only and MultiHeaders-only), then asserts
// backwards compatibility (response Headers) and the new feature (response MultiHeaders).
func runMultiHeadersTest(cfg config.Config, nodeRuntime cre.NodeRuntime, client *http.Client, log *slog.Logger) (string, error) {
	timeout := &durationpb.Duration{Seconds: 10}

	// 1) Headers-only request: assert response Headers match sent and Set-Cookie is comma-joined (backwards compat).
	sentHeaders := map[string]string{
		"Content-Type":    "application/json",
		"Accept-Language": "en,fr",
	}
	resp1, err := client.SendRequest(nodeRuntime, &http.Request{
		Url:     cfg.URL,
		Method:  cfg.Method,
		Headers: sentHeaders,
		Body:    []byte(cfg.Body),
		Timeout: timeout,
	}).Await()
	if err != nil {
		return "", fmt.Errorf("HTTP Action multi-headers (Headers request) failed: %w", err)
	}
	if !statusOK(resp1.StatusCode) {
		return "", fmt.Errorf("HTTP Action multi-headers (Headers request) status: %d", resp1.StatusCode)
	}
	h := resp1.GetHeaders() //nolint:staticcheck
	if err := assertMapContains(h, sentHeaders, "HTTP Action multi-headers response Headers (backwards compat)"); err != nil {
		return "", err
	}
	setCookieJoined, ok := h["Set-Cookie"]
	if !ok {
		return "", fmt.Errorf("HTTP Action multi-headers test failed: Set-Cookie not in response Headers (backwards compat)")
	}
	if slices.IndexFunc(expectedSetCookieSubstrings, func(sub string) bool { return !strings.Contains(setCookieJoined, sub) }) != -1 {
		return "", fmt.Errorf("HTTP Action multi-headers test failed: Set-Cookie in Headers should be comma-joined with all three values, got %q", setCookieJoined)
	}
	log.Info("HTTP Action multi-headers test: Headers (backwards compat) OK")

	// 2) MultiHeaders-only request: assert response MultiHeaders match sent and Set-Cookie has three distinct values.
	sentMultiHeaders := map[string]*http.HeaderValues{
		"Content-Type":    {Values: []string{"application/json"}},
		"Accept-Language": {Values: []string{"en", "fr"}},
	}
	resp2, err := client.SendRequest(nodeRuntime, &http.Request{
		Url:          cfg.URL,
		Method:       cfg.Method,
		MultiHeaders: sentMultiHeaders,
		Body:         []byte(cfg.Body),
		Timeout:      timeout,
	}).Await()
	if err != nil {
		return "", fmt.Errorf("HTTP Action multi-headers (MultiHeaders request) failed: %w", err)
	}
	if !statusOK(resp2.StatusCode) {
		return "", fmt.Errorf("HTTP Action multi-headers (MultiHeaders request) status: %d", resp2.StatusCode)
	}
	mh := resp2.GetMultiHeaders()
	if err := assertMultiHeadersContains(mh, sentMultiHeaders, "HTTP Action multi-headers response MultiHeaders"); err != nil {
		return "", err
	}
	setCookieHV, ok := mh["Set-Cookie"]
	if !ok || setCookieHV == nil {
		return "", fmt.Errorf("Set-Cookie not in MultiHeaders")
	}
	vals := setCookieHV.GetValues()
	if len(vals) != 3 {
		return "", fmt.Errorf("Set-Cookie in MultiHeaders should have 3 distinct values, got %d: %v", len(vals), vals)
	}
	for _, sub := range expectedSetCookieSubstrings {
		if !slices.ContainsFunc(vals, func(v string) bool { return strings.Contains(v, sub) }) {
			return "", fmt.Errorf("Set-Cookie MultiHeaders missing expected value containing %q, got %v", sub, vals)
		}
	}
	log.Info("HTTP Action multi-headers test passed", "setCookieCount", len(vals))
	return "HTTP Action multi-headers test completed", nil
}

// runBothSetRegressionTest sends a request with both Headers and MultiHeaders set; the capability
// must reject it with a user error. This is a regression test to ensure validation is enforced.
func runBothSetRegressionTest(cfg config.Config, nodeRuntime cre.NodeRuntime, client *http.Client, log *slog.Logger) (string, error) {
	timeout := &durationpb.Duration{Seconds: 10}
	req := &http.Request{
		Url:     cfg.URL,
		Method:  cfg.Method,
		Headers: map[string]string{"X-Test": "value"},
		MultiHeaders: map[string]*http.HeaderValues{
			"Accept": {Values: []string{"application/json"}},
		},
		Body:    []byte(cfg.Body),
		Timeout: timeout,
	}
	_, err := client.SendRequest(nodeRuntime, req).Await()
	if err == nil {
		return "", fmt.Errorf("multi-headers regression: expected user error when both Headers and MultiHeaders are set, but request succeeded")
	}
	if !errorContainsAll(err, bothSetRegressionExpectedSubstrings) {
		return "", fmt.Errorf("multi-headers regression: expected user error containing %v, got: %w", bothSetRegressionExpectedSubstrings, err)
	}
	log.Info("HTTP Action multi-headers regression passed: both Headers and MultiHeaders set correctly rejected")
	return bothSetRegressionSuccess, nil
}

// --- Helper functions ---

// statusOK reports whether code is in the 2xx range.
func statusOK(code uint32) bool { return code >= 200 && code < 300 }

// errorContainsAll reports whether err is non-nil and its message contains every substring.
func errorContainsAll(err error, substrings []string) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return slices.IndexFunc(substrings, func(sub string) bool { return !strings.Contains(s, sub) }) == -1
}

// assertMapContains checks that got contains every key from want with the same value.
func assertMapContains(got map[string]string, want map[string]string, context string) error {
	if got == nil {
		return fmt.Errorf("%s: response map is nil", context)
	}
	for name, wantVal := range want {
		gotVal, ok := got[name]
		if !ok {
			return fmt.Errorf("%s: missing key %q", context, name)
		}
		if gotVal != wantVal {
			return fmt.Errorf("%s: %q = %q, want %q", context, name, gotVal, wantVal)
		}
	}
	return nil
}

// assertMultiHeadersContains checks that got contains every key from want with the same Values slice.
func assertMultiHeadersContains(got map[string]*http.HeaderValues, want map[string]*http.HeaderValues, context string) error {
	if got == nil {
		return fmt.Errorf("%s: response MultiHeaders is nil", context)
	}
	for name, wantHV := range want {
		gotHV, ok := got[name]
		if !ok || gotHV == nil {
			return fmt.Errorf("%s: missing or nil key %q", context, name)
		}
		gotVals := gotHV.GetValues()
		wantVals := wantHV.GetValues()
		if !slices.Equal(gotVals, wantVals) {
			return fmt.Errorf("%s: %q = %v, want %v", context, name, gotVals, wantVals)
		}
	}
	return nil
}
