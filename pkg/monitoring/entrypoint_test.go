package monitoring

import (
	"context"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/require"
)

const testEntrypointE2eDurationSec = 10

func TestEntrypointE2e(t *testing.T) {
	if _, isPresent := os.LookupEnv("FEATURE_TEST_ONLY_ENV_RUNNING"); !isPresent {
		t.Skip()
	}

	// Configuration for the process via env vars.
	os.Setenv("KAFKA_BROKERS", "localhost:29092")
	os.Setenv("KAFKA_CLIENT_ID", "client")
	os.Setenv("KAFKA_SECURITY_PROTOCOL", "PLAINTEXT")
	os.Setenv("KAFKA_SASL_MECHANISM", "PLAIN")
	os.Setenv("KAFKA_SASL_USERNAME", "")
	os.Setenv("KAFKA_SASL_PASSWORD", "")
	os.Setenv("KAFKA_TRANSMISSION_TOPIC", "transmission")
	os.Setenv("KAFKA_CONFIG_SET_SIMPLIFIED_TOPIC", "config_set_simplified")
	os.Setenv("SCHEMA_REGISTRY_URL", "http://localhost:8989")
	os.Setenv("SCHEMA_REGISTRY_USERNAME", "")
	os.Setenv("SCHEMA_REGISTRY_PASSWORD", "")
	os.Setenv("FEEDS_URL", "http://some-feeds.com")
	os.Setenv("HTTP_ADDRESS", "http://localhost:3000")
	os.Setenv("FEATURE_TEST_ONLY_FAKE_READERS", "true")
	os.Setenv("FEATURE_TEST_ONLY_FAKE_RDD", "true")

	defer os.Unsetenv("KAFKA_BROKERS")
	defer os.Unsetenv("KAFKA_CLIENT_ID")
	defer os.Unsetenv("KAFKA_SECURITY_PROTOCOL")
	defer os.Unsetenv("KAFKA_SASL_MECHANISM")
	defer os.Unsetenv("KAFKA_SASL_USERNAME")
	defer os.Unsetenv("KAFKA_SASL_PASSWORD")
	defer os.Unsetenv("KAFKA_TRANSMISSION_TOPIC")
	defer os.Unsetenv("KAFKA_CONFIG_SET_SIMPLIFIED_TOPIC")
	defer os.Unsetenv("SCHEMA_REGISTRY_URL")
	defer os.Unsetenv("SCHEMA_REGISTRY_USERNAME")
	defer os.Unsetenv("SCHEMA_REGISTRY_PASSWORD")
	defer os.Unsetenv("FEEDS_URL")
	defer os.Unsetenv("HTTP_ADDRESS")
	defer os.Unsetenv("FEATURE_TEST_ONLY_FAKE_READERS")
	defer os.Unsetenv("FEATURE_TEST_ONLY_FAKE_RDD")

	ctx, cancel := context.WithTimeout(context.Background(), testEntrypointE2eDurationSec*time.Second)
	defer cancel()

	chainConfig := generateChainConfig()

	entrypoint, err := NewEntrypoint(
		ctx,
		newNullLogger(),
		chainConfig,
		// These are not needed as all the sources are faked. See config.Feature
		nil,
		nil,
		func(buf io.ReadCloser) ([]FeedConfig, error) { return []FeedConfig{}, nil },
	)
	require.NoError(t, err)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		entrypoint.Run()
	}()

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:29092",
		"client.id":         "test-entrypoint",
		"group.id":          "test-entrypoint",
		"security.protocol": "PLAINTEXT",
		"sasl.mechanisms":   "PLAIN",
		"auto.offset.reset": "earliest",
	})
	require.NoError(t, err)
	err = consumer.SubscribeTopics(
		[]string{"config_set_simplified", "transmission"},
		func(*kafka.Consumer, kafka.Event) error { return nil }, // kafka.RebalanceCb
	)
	require.NoError(t, err)

	events := []string{}
	i := 0
	for i < 10 {
		select {
		case <-ctx.Done():
			break
		default:
			if event := consumer.Poll(1000); event != nil {
				events = append(events, event.String())
				i += 1
			}
		}
	}

	err = consumer.Unsubscribe()
	require.NoError(t, err)
	cancel()
	wg.Wait()

	require.NotEqual(t, 0, len(events))
	foundConfig, foundTransmission := false, false
	for _, event := range events {
		if strings.HasPrefix("transmission", event) {
			foundConfig = true
		} else if strings.HasPrefix("config_set_simplified", event) {
			foundTransmission = true
		}
	}
	require.True(t, foundConfig && foundTransmission, "both transmission and config messages")
}
