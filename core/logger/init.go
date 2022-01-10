package logger

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"

	"github.com/smartcontractkit/chainlink/core/static"
)

func init() {
	initMemorySink()
	initConsoleSink()
	initSentry()
}

func initMemorySink() {
	err := zap.RegisterSink("memory", func(*url.URL) (zap.Sink, error) {
		return PrettyConsoleSink{Sink: &testMemoryLog}, nil
	})
	if err != nil {
		panic(err)
	}
}

func initConsoleSink() {
	err := zap.RegisterSink("pretty", func(*url.URL) (zap.Sink, error) {
		return PrettyConsoleSink{os.Stderr}, nil
	})
	if err != nil {
		panic(fmt.Sprintf("failed to register pretty printer %+v", err))
	}
}

func initSentry() {
	// If SENTRY_DSN is set at runtime, sentry will be enabled and send metrics to this URL
	sentrydsn := os.Getenv("SENTRY_DSN")
	if sentrydsn == "" {
		// Do not initialize sentry at all if the DSN is missing
		return
	}

	// If SENTRY_ENVIRONMENT is set, it will override everything. Otherwise infers from CHAINLINK_DEV.
	var sentryenv string
	if env := os.Getenv("SENTRY_ENVIRONMENT"); env != "" {
		sentryenv = env
	} else if os.Getenv("CHAINLINK_DEV") == "true" {
		sentryenv = "dev"
	} else {
		sentryenv = "prod"
	}

	// If SENTRY_RELEASE is set, it will override everything. Otherwise, static.Version will be used.
	var sentryrelease string
	if release := os.Getenv("SENTRY_RELEASE"); release != "" {
		sentryrelease = release
	} else {
		sentryrelease = static.Version
	}

	// Set SENTRY_DEBUG=true to enable printing of SDK debug messages
	sentrydebug := os.Getenv("SENTRY_DEBUG") == "true"

	err := sentry.Init(sentry.ClientOptions{
		// AttachStacktrace is needed to send stacktrace alongside panics
		AttachStacktrace: true,
		Dsn:              sentrydsn,
		Environment:      sentryenv,
		Release:          sentryrelease,
		Debug:            sentrydebug,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
}
