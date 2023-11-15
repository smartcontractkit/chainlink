package loop

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	envPromPort               = "CL_PROMETHEUS_PORT"
	envTracingEnabled         = "CL_TRACING_ENABLED"
	envTracingCollectorTarget = "CL_TRACING_COLLECTOR_TARGET"
	envTracingSamplingRatio   = "CL_TRACING_SAMPLING_RATIO"
	envTracingAttribute       = "CL_TRACING_ATTRIBUTE_"
	envTracingTLSCertPath     = "CL_TRACING_TLS_CERT_PATH"
)

// EnvConfig is the configuration between the application and the LOOP executable. The values
// are fully resolved and static and passed via the environment.
type EnvConfig struct {
	PrometheusPort int

	TracingEnabled         bool
	TracingCollectorTarget string
	TracingSamplingRatio   float64
	TracingTLSCertPath     string
	TracingAttributes      map[string]string
}

// AsCmdEnv returns a slice of environment variable key/value pairs for an exec.Cmd.
func (e *EnvConfig) AsCmdEnv() (env []string) {
	injectEnv := map[string]string{
		envPromPort:               strconv.Itoa(e.PrometheusPort),
		envTracingEnabled:         strconv.FormatBool(e.TracingEnabled),
		envTracingCollectorTarget: e.TracingCollectorTarget,
		envTracingSamplingRatio:   strconv.FormatFloat(e.TracingSamplingRatio, 'f', -1, 64),
		envTracingTLSCertPath:     e.TracingTLSCertPath,
	}

	for k, v := range e.TracingAttributes {
		injectEnv[envTracingAttribute+k] = v
	}

	for k, v := range injectEnv {
		env = append(env, k+"="+v)
	}
	return
}

// parse deserializes environment variables
func (e *EnvConfig) parse() error {
	promPortStr := os.Getenv(envPromPort)
	var err error
	e.PrometheusPort, err = strconv.Atoi(promPortStr)
	if err != nil {
		return fmt.Errorf("failed to parse %s = %q: %w", envPromPort, promPortStr, err)
	}

	e.TracingEnabled, err = getTracingEnabled()
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", envTracingEnabled, err)
	}

	if e.TracingEnabled {
		e.TracingCollectorTarget, err = getValidCollectorTarget()
		if err != nil {
			return err
		}
		e.TracingAttributes = getTracingAttributes()
		e.TracingSamplingRatio = getTracingSamplingRatio()
		e.TracingTLSCertPath = getTLSCertPath()
	}
	return nil
}

// isTracingEnabled parses and validates the TRACING_ENABLED environment variable.
func getTracingEnabled() (bool, error) {
	tracingEnabledString := os.Getenv(envTracingEnabled)
	if tracingEnabledString == "" {
		return false, nil
	}
	return strconv.ParseBool(tracingEnabledString)
}

// getValidCollectorTarget validates TRACING_COLLECTOR_TARGET as a URL.
func getValidCollectorTarget() (string, error) {
	tracingCollectorTarget := os.Getenv(envTracingCollectorTarget)
	_, err := url.ParseRequestURI(tracingCollectorTarget)
	if err != nil {
		return "", fmt.Errorf("invalid %s: %w", envTracingCollectorTarget, err)
	}
	return tracingCollectorTarget, nil
}

// getTracingAttributes collects attributes prefixed with TRACING_ATTRIBUTE_.
func getTracingAttributes() map[string]string {
	tracingAttributes := make(map[string]string)
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, envTracingAttribute) {
			tracingAttributes[strings.TrimPrefix(env, envTracingAttribute)] = os.Getenv(env)
		}
	}
	return tracingAttributes
}

// getTracingSamplingRatio parses the TRACING_SAMPLING_RATIO environment variable.
// Any errors in parsing result in a sampling ratio of 0.0.
func getTracingSamplingRatio() float64 {
	tracingSamplingRatio := os.Getenv(envTracingSamplingRatio)
	if tracingSamplingRatio == "" {
		return 0.0
	}
	samplingRatio, err := strconv.ParseFloat(tracingSamplingRatio, 64)
	if err != nil {
		return 0.0
	}
	return samplingRatio
}

// getTLSCertPath parses the CL_TRACING_TLS_CERT_PATH environment variable.
func getTLSCertPath() string {
	return os.Getenv(envTracingTLSCertPath)
}
