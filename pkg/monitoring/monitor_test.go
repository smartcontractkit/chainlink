package monitoring

import (
	"context"
	"io"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"github.com/smartcontractkit/chainlink-common/pkg/utils"
)

const testMonitorDurationSec = 15

func TestMonitor(t *testing.T) {
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
	os.Setenv("NODES_URL", "http://some-nodes.com")
	os.Setenv("HTTP_ADDRESS", "http://localhost:3000")

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

	ctx, cancel := context.WithTimeout(context.Background(), testMonitorDurationSec*time.Second)
	defer cancel()

	chainConfig := fakeChainConfig{
		ReadTimeout:  100 * time.Millisecond,
		PollInterval: 100 * time.Millisecond,
	}

	stopCh := make(chan struct{})
	context.AfterFunc(ctx, func() { close(stopCh) })
	monitor, err := NewMonitor(
		stopCh,
		newNullLogger(),
		chainConfig,
		&fakeRandomDataSourceFactory{make(chan interface{})},
		&fakeRandomDataSourceFactory{make(chan interface{})},
		func(buf io.ReadCloser) ([]FeedConfig, error) { return []FeedConfig{}, nil },
		func(buf io.ReadCloser) ([]NodeConfig, error) { return []NodeConfig{}, nil },
	)
	require.NoError(t, err)

	var monitorSubs utils.Subprocesses
	monitorSubs.Go(func() {
		monitor.Run()
	})

	// Wait for the monitor to start.
	<-time.After(5 * time.Second)

	kafkaConfig := &kafka.ConfigMap{
		"bootstrap.servers": "localhost:29092",
		"client.id":         "test-monitor",
		"group.id":          "test-monitor",
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

	var transmissionsCounter uint64
	var configsCounter uint64
	var countersSubs utils.Subprocesses
	countersSubs.Go(func() {
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
	})
	countersSubs.Go(func() {
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
	})
	countersSubs.Wait()

	cancel()
	monitorSubs.Wait()

	err = consumerConfig.Unsubscribe()
	require.NoError(t, err)
	err = consumerTransmission.Unsubscribe()
	require.NoError(t, err)

	require.Equal(t, uint64(10), configsCounter)
	require.Equal(t, uint64(10), transmissionsCounter)
}
