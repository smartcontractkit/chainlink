package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"
)

func Parse() (Config, error) {
	cfg := Config{}

	if err := parseEnvVars(&cfg); err != nil {
		return cfg, err
	}

	applyDefaults(&cfg)

	err := validateConfig(cfg)
	return cfg, err
}

func parseEnvVars(cfg *Config) error {
	if value, isPresent := os.LookupEnv("KAFKA_BROKERS"); isPresent {
		cfg.Kafka.Brokers = value
	}
	if value, isPresent := os.LookupEnv("KAFKA_CLIENT_ID"); isPresent {
		cfg.Kafka.ClientID = value
	}
	if value, isPresent := os.LookupEnv("KAFKA_SECURITY_PROTOCOL"); isPresent {
		cfg.Kafka.SecurityProtocol = value
	}

	if value, isPresent := os.LookupEnv("KAFKA_SASL_MECHANISM"); isPresent {
		cfg.Kafka.SaslMechanism = value
	}
	if value, isPresent := os.LookupEnv("KAFKA_SASL_USERNAME"); isPresent {
		cfg.Kafka.SaslUsername = value
	}
	if value, isPresent := os.LookupEnv("KAFKA_SASL_PASSWORD"); isPresent {
		cfg.Kafka.SaslPassword = value
	}

	if value, isPresent := os.LookupEnv("KAFKA_TRANSMISSION_TOPIC"); isPresent {
		cfg.Kafka.TransmissionTopic = value
	}
	if value, isPresent := os.LookupEnv("KAFKA_CONFIG_SET_SIMPLIFIED_TOPIC"); isPresent {
		cfg.Kafka.ConfigSetSimplifiedTopic = value
	}

	if value, isPresent := os.LookupEnv("SCHEMA_REGISTRY_URL"); isPresent {
		cfg.SchemaRegistry.URL = value
	}
	if value, isPresent := os.LookupEnv("SCHEMA_REGISTRY_USERNAME"); isPresent {
		cfg.SchemaRegistry.Username = value
	}
	if value, isPresent := os.LookupEnv("SCHEMA_REGISTRY_PASSWORD"); isPresent {
		cfg.SchemaRegistry.Password = value
	}

	if value, isPresent := os.LookupEnv("FEEDS_URL"); isPresent {
		cfg.Feeds.URL = value
	}
	if value, isPresent := os.LookupEnv("FEEDS_RDD_READ_TIMEOUT"); isPresent {
		readTimeout, err := time.ParseDuration(value)
		if err != nil {
			return fmt.Errorf("failed to parse env var FEEDS_RDD_READ_TIMEOUT, see https://pkg.go.dev/time#ParseDuration: %w", err)
		}
		cfg.Feeds.RDDReadTimeout = readTimeout
	}
	if value, isPresent := os.LookupEnv("FEEDS_RDD_POLL_INTERVAL"); isPresent {
		pollInterval, err := time.ParseDuration(value)
		if err != nil {
			return fmt.Errorf("failed to parse env var FEEDS_RDD_POLL_INTERVAL, see https://pkg.go.dev/time#ParseDuration: %w", err)
		}
		cfg.Feeds.RDDPollInterval = pollInterval
	}

	if value, isPresent := os.LookupEnv("HTTP_ADDRESS"); isPresent {
		cfg.HTTP.Address = value
	}

	if value, isPresent := os.LookupEnv("FEATURE_TEST_ONLY_FAKE_READERS"); isPresent {
		isTestMode, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("failed to parse boolean env var '%s'. See https://pkg.go.dev/strconv#ParseBool", "FEATURE_TEST_ONLY_FAKE_READERS")
		}
		cfg.Feature.TestOnlyFakeReaders = isTestMode
	}
	if value, isPresent := os.LookupEnv("FEATURE_TEST_ONLY_FAKE_RDD"); isPresent {
		isTestMode, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("failed to parse boolean env var '%s'. See https://pkg.go.dev/strconv#ParseBool", "FEATURE_TEST_ONLY_FAKE_RDD")
		}
		cfg.Feature.TestOnlyFakeRdd = isTestMode
	}

	return nil
}

func applyDefaults(cfg *Config) {
	if cfg.Feeds.RDDReadTimeout == 0 {
		cfg.Feeds.RDDReadTimeout = 1 * time.Second
	}
	if cfg.Feeds.RDDPollInterval == 0 {
		cfg.Feeds.RDDPollInterval = 10 * time.Second
	}
}

func validateConfig(cfg Config) error {
	// Required config
	for envVarName, currentValue := range map[string]string{
		"KAFKA_BROKERS":           cfg.Kafka.Brokers,
		"KAFKA_CLIENT_ID":         cfg.Kafka.ClientID,
		"KAFKA_SECURITY_PROTOCOL": cfg.Kafka.SecurityProtocol,
		"KAFKA_SASL_MECHANISM":    cfg.Kafka.SaslMechanism,

		"KAFKA_TRANSMISSION_TOPIC":          cfg.Kafka.TransmissionTopic,
		"KAFKA_CONFIG_SET_SIMPLIFIED_TOPIC": cfg.Kafka.ConfigSetSimplifiedTopic,

		"SCHEMA_REGISTRY_URL": cfg.SchemaRegistry.URL,

		"FEEDS_URL": cfg.Feeds.URL,

		"HTTP_ADDRESS": cfg.HTTP.Address,
	} {
		if currentValue == "" {
			return fmt.Errorf("'%s' env var is required", envVarName)
		}
	}
	// Validate URLs.
	for envVarName, currentValue := range map[string]string{
		"SCHEMA_REGISTRY_URL": cfg.SchemaRegistry.URL,
		"FEEDS_URL":           cfg.Feeds.URL,
	} {
		if _, err := url.ParseRequestURI(currentValue); err != nil {
			return fmt.Errorf("%s='%s' is not a valid URL: %w", envVarName, currentValue, err)
		}
	}
	return nil
}
