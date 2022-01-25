package monitoring

import (
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/smartcontractkit/chainlink-relay/pkg/monitoring/config"
)

// Producer is an abstraction on top of Kafka to aid with tests.
type Producer interface {
	Produce(key, value []byte, topic string) error
}

type producer struct {
	log          Logger
	backend      *kafka.Producer
	deliveryChan chan kafka.Event
	cfg          config.Kafka
}

func NewProducer(ctx context.Context, log Logger, cfg config.Kafka) (Producer, error) {
	backend, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Brokers,
		"client.id":         cfg.ClientID,
		"security.protocol": cfg.SecurityProtocol,
		"sasl.mechanisms":   cfg.SaslMechanism,
		"sasl.username":     cfg.SaslUsername,
		"sasl.password":     cfg.SaslPassword,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}
	p := &producer{
		log,
		backend,
		make(chan kafka.Event),
		cfg,
	}
	go p.drainDeliveryChan(ctx)
	return p, nil
}

// drainDeliveryChan should be executed as a goroutine.
func (p *producer) drainDeliveryChan(ctx context.Context) {
	for {
		select {
		case event := <-p.deliveryChan:
			p.log.Debugw("received delivery event", "event", event.String())
		case <-ctx.Done():
			p.backend.Close()
			return
		}
	}
}

func (p *producer) Produce(key, value []byte, topic string) error {
	return p.backend.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:   key,
		Value: value,
	}, p.deliveryChan)
}
