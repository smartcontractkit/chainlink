package monitoring

import (
	"context"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

const testEntrypointDurationSec = 15

func TestEntrypoint(t *testing.T) {
	if _, isPresent := os.LookupEnv("FEATURE_TEST_ONLY_ENV_RUNNING"); !isPresent {
		t.Skip()
	}
	defer goleak.VerifyNone(t)

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

	ctx, cancel := context.WithTimeout(context.Background(), testEntrypointDurationSec*time.Second)
	defer cancel()

	chainConfig := fakeChainConfig{
		ReadTimeout:  100 * time.Millisecond,
		PollInterval: 100 * time.Millisecond,
	}

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

	entrypointWg := &sync.WaitGroup{}
	entrypointWg.Add(1)
	go func() {
		defer entrypointWg.Done()
		entrypoint.Run()
	}()
	// Wait for the entrypoint to start.
	<-time.After(5 * time.Second)

	kafkaConfig := &kafka.ConfigMap{
		"bootstrap.servers": "localhost:29092",
		"client.id":         "test-entrypoint",
		"group.id":          "test-entrypoint",
		"security.protocol": "PLAINTEXT",
		"sasl.mechanisms":   "PLAIN",
		"auto.offset.reset": "earliest",
	}
	consumerConfig, err := kafka.NewConsumer(kafkaConfig)
	require.NoError(t, err)
	err = consumerConfig.Subscribe("config_set_simplified", nil)
	require.NoError(t, err)

	consumerTransmission, err := kafka.NewConsumer(kafkaConfig)
	require.NoError(t, err)
	err = consumerTransmission.Subscribe("transmission", nil)
	require.NoError(t, err)

	// Wait for the subscriptions to start.
	<-time.After(2 * time.Second)

	var transmissionsCounter uint64 = 0
	var configsCounter uint64 = 0
	countersWg := &sync.WaitGroup{}
	countersWg.Add(2)
	go func() {
		defer countersWg.Done()
		for i := 0; i < 10; {
			select {
			case <-ctx.Done():
				return
			default:
				event := consumerTransmission.Poll(500)
				if event != nil {
					atomic.AddUint64(&transmissionsCounter, 1)
					i++
				}
			}
		}
	}()
	go func() {
		defer countersWg.Done()
		for i := 0; i < 10; {
			select {
			case <-ctx.Done():
				return
			default:
				event := consumerConfig.Poll(500)
				if event != nil {
					atomic.AddUint64(&configsCounter, 1)
					i++
				}
			}
		}
	}()
	countersWg.Wait()

	cancel()
	entrypointWg.Wait()

	err = consumerConfig.Unsubscribe()
	require.NoError(t, err)
	err = consumerTransmission.Unsubscribe()
	require.NoError(t, err)

	require.Equal(t, uint64(10), configsCounter)
	require.Equal(t, uint64(10), transmissionsCounter)
}
