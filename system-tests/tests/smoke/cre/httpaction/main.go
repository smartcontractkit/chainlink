//go:build wasip1

package main

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/httpaction/config"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
	"google.golang.org/protobuf/types/known/durationpb"
	"gopkg.in/yaml.v3"
)

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

func onCronTrigger(wfCfg config.Config, runtime cre.Runtime, payload *cron.Payload) (_ any, _ error) {
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

// Expected Set-Cookie values from the fake server (v2_http_action_test.go).
const (
	setCookieSession = "sessionid=multi-e2e-1"
	setCookieCSRF    = "csrf=multi-e2e-2"
	setCookiePref    = "pref=multi-e2e-3"
)

// bothSetRegressionExpectedSubstrings are substrings that must appear in the capability's user error
// when both Headers and MultiHeaders are set on a request (validation rejects).
const (
	bothSetRegressionExpected1 = "Headers or MultiHeaders"
	bothSetRegressionExpected2 = "not both"
	bothSetRegressionSuccess   = "HTTP Action multi-headers regression completed"
)

// runBothSetRegressionTest sends a request with both Headers and MultiHeaders set; the capability
// must reject it with a user error. This is a regression test to ensure validation is enforced.
func runBothSetRegressionTest(cfg config.Config, nodeRuntime cre.NodeRuntime, client *http.Client, log interface{ Info(string, ...any) }) (string, error) {
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
	errStr := err.Error()
	if !strings.Contains(errStr, bothSetRegressionExpected1) || !strings.Contains(errStr, bothSetRegressionExpected2) {
		return "", fmt.Errorf("multi-headers regression: expected user error containing %q and %q, got: %w", bothSetRegressionExpected1, bothSetRegressionExpected2, err)
	}
	log.Info("HTTP Action multi-headers regression passed: both Headers and MultiHeaders set correctly rejected")
	return bothSetRegressionSuccess, nil
}

// runMultiHeadersTest sends two requests (Headers-only and MultiHeaders-only), then asserts
// backwards compatibility (response Headers) and the new feature (response MultiHeaders).
func runMultiHeadersTest(cfg config.Config, nodeRuntime cre.NodeRuntime, client *http.Client, log interface{ Info(string, ...any) }) (string, error) {
	timeout := &durationpb.Duration{Seconds: 10}

	// 1) Request using Headers field only; assert all sent headers in response Headers (backwards compatibility).
	// Set-Cookie in response Headers must be a comma-joined string.
	sentHeaders := map[string]string{
		"Content-Type":    "application/json",
		"Accept-Language": "en,fr",
	}
	reqHeaders := &http.Request{
		Url:     cfg.URL,
		Method:  cfg.Method,
		Headers: sentHeaders,
		Body:    []byte(cfg.Body),
		Timeout: timeout,
	}
	resp1, err := client.SendRequest(nodeRuntime, reqHeaders).Await()
	if err != nil {
		return "", fmt.Errorf("HTTP Action multi-headers (Headers request) failed: %w", err)
	}
	if resp1.StatusCode < 200 || resp1.StatusCode >= 300 {
		return "", fmt.Errorf("HTTP Action multi-headers (Headers request) status: %d", resp1.StatusCode)
	}
	h := resp1.GetHeaders() //nolint:staticcheck
	if h == nil {
		return "", fmt.Errorf("HTTP Action multi-headers test failed: response Headers is nil (backwards compat)")
	}
	for name, want := range sentHeaders {
		got, ok := h[name]
		if !ok {
			return "", fmt.Errorf("HTTP Action multi-headers test failed: sent header %q not in response Headers (backwards compat)", name)
		}
		if got != want {
			return "", fmt.Errorf("HTTP Action multi-headers test failed: response Headers[%q] = %q, want %q", name, got, want)
		}
	}

	// Note: this is a malformatted Header, but kept for backwards compatibility
	setCookieJoined, ok := h["Set-Cookie"]
	if !ok {
		return "", fmt.Errorf("HTTP Action multi-headers test failed: Set-Cookie not in response Headers (backwards compat)")
	}

	if !strings.Contains(setCookieJoined, setCookieSession) || !strings.Contains(setCookieJoined, setCookieCSRF) || !strings.Contains(setCookieJoined, setCookiePref) {
		return "", fmt.Errorf("HTTP Action multi-headers test failed: Set-Cookie in Headers should be comma-joined with all three values, got %q", setCookieJoined)
	}
	log.Info("HTTP Action multi-headers test: Headers (backwards compat) OK")

	// 2) Request using MultiHeaders field; assert all sent headers in response, and Set-Cookie has three distinct values in MultiHeaders.
	sentMultiHeaders := map[string]*http.HeaderValues{
		"Content-Type":    {Values: []string{"application/json"}},
		"Accept-Language": {Values: []string{"en", "fr"}},
	}
	reqMulti := &http.Request{
		Url:          cfg.URL,
		Method:       cfg.Method,
		MultiHeaders: sentMultiHeaders,
		Body:         []byte(cfg.Body),
		Timeout:      timeout,
	}
	resp2, err := client.SendRequest(nodeRuntime, reqMulti).Await()
	if err != nil {
		return "", fmt.Errorf("HTTP Action multi-headers (MultiHeaders request) failed: %w", err)
	}
	if resp2.StatusCode < 200 || resp2.StatusCode >= 300 {
		return "", fmt.Errorf("HTTP Action multi-headers (MultiHeaders request) status: %d", resp2.StatusCode)
	}
	mh := resp2.GetMultiHeaders()
	if mh == nil {
		return "", fmt.Errorf("response MultiHeaders is nil")
	}
	for name, wantHV := range sentMultiHeaders {
		gotHV, ok := mh[name]
		if !ok || gotHV == nil {
			return "", fmt.Errorf("sent header %q not in response MultiHeaders", name)
		}
		gotVals := gotHV.GetValues()
		wantVals := wantHV.GetValues()
		if len(gotVals) != len(wantVals) {
			return "", fmt.Errorf("response MultiHeaders[%q] has %d values, want %d", name, len(gotVals), len(wantVals))
		}
		for i, w := range wantVals {
			if i >= len(gotVals) || gotVals[i] != w {
				return "", fmt.Errorf("response MultiHeaders[%q] = %v, want %v", name, gotVals, wantVals)
			}
		}
	}
	setCookieHV, ok := mh["Set-Cookie"]
	if !ok || setCookieHV == nil {
		return "", fmt.Errorf("Set-Cookie not in MultiHeaders")
	}
	vals := setCookieHV.GetValues()
	if len(vals) != 3 {
		return "", fmt.Errorf("Set-Cookie in MultiHeaders should have 3 distinct values, got %d: %v", len(vals), vals)
	}
	seen := map[string]bool{}
	for _, v := range vals {
		seen[v] = true
	}
	for _, sub := range []string{setCookieSession, setCookieCSRF, setCookiePref} {
		found := false
		for v := range seen {
			if strings.Contains(v, sub) {
				found = true
				break
			}
		}
		if !found {
			return "", fmt.Errorf("Set-Cookie MultiHeaders missing expected value containing %q, got %v", sub, vals)
		}
	}
	log.Info("HTTP Action multi-headers test passed", "setCookieCount", len(vals))
	return "HTTP Action multi-headers test completed", nil
}

func runCRUDSuccessTest(wfCfg config.Config, runtime cre.Runtime) (string, error) {
	logger := runtime.Logger()
	logger.Info(
		"Running HTTP Action capability",
		"testCase",
		wfCfg.TestCase,
		"method",
		wfCfg.Method,
		"url",
		wfCfg.URL,
	)

	crudPromise := cre.RunInNodeMode(wfCfg, runtime,
		func(cfg config.Config, nodeRuntime cre.NodeRuntime) (string, error) {
			client := &http.Client{}

			if cfg.TestCase == "multi-headers" {
				return runMultiHeadersTest(cfg, nodeRuntime, client, logger)
			}
			if cfg.TestCase == "mh-regression-both" {
				return runBothSetRegressionTest(cfg, nodeRuntime, client, logger)
			}

			req := &http.Request{
				Url:     cfg.URL,
				Method:  cfg.Method,
				Headers: map[string]string{"Content-Type": "application/json"},
				Body:    []byte(cfg.Body),
				Timeout: &durationpb.Duration{Seconds: 10},
			}

			logger.Info("Testing HTTP Action capability with configuration",
				"url", req.Url,
				"method", req.Method,
				"hasBody", len(cfg.Body) > 0)

			resp, err := client.SendRequest(nodeRuntime, req).Await()
			if err != nil {
				logger.Error(
					"Failed to complete HTTP Action request",
					"error",
					err,
					"url",
					req.Url,
					"method",
					req.Method,
				)
				return "", fmt.Errorf("HTTP Action %s request failed: %w", req.Method, err)
			}

			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				logger.Error(
					"Failed response to HTTP Action request",
					"status",
					resp.StatusCode,
					"url",
					req.Url,
					"method",
					req.Method,
				)
				return "", fmt.Errorf("HTTP Action %s request failed with status: %d", req.Method, resp.StatusCode)
			}

			logger.Info(
				"HTTP Action completed",
				"url",
				req.Url,
				"method",
				req.Method,
				"status",
				resp.StatusCode,
				"body",
				string(resp.Body),
			)

			return fmt.Sprintf("HTTP Action CRUD success test completed: %s", cfg.TestCase), nil
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
