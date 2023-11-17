package monitoring

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
)

// Pipeline represents a succession of transformations on the data coming from a source:
// source output -> adapt to a map -> encode to AVRO -> send to Kafka
type Pipeline struct {
	Topic  string
	Mapper Mapper
	Schema Schema
}

// NewKafkaExporterFactory produces Kafka exporters which consume, format and publish source outputs to kafka.
func NewKafkaExporterFactory(
	log Logger,
	producer Producer,
	pipelines []Pipeline,
) (ExporterFactory, error) {
	// Check pipeline topics match schema subjects.
	for _, pipeline := range pipelines {
		if SubjectFromTopic(pipeline.Topic) != pipeline.Schema.Subject() {
			return nil, fmt.Errorf("topic '%s' does not match schema subject '%s'", pipeline.Topic, pipeline.Schema.Subject())
		}
	}
	return &kafkaExporterFactory{
		log,
		producer,
		pipelines,
	}, nil
}

type kafkaExporterFactory struct {
	log       Logger
	producer  Producer
	pipelines []Pipeline
}

func (k *kafkaExporterFactory) NewExporter(
	params ExporterParams,
) (Exporter, error) {
	return &kafkaExporter{
		params.ChainConfig,
		params.FeedConfig,

		logger.With(k.log, "feed", params.FeedConfig.GetName()),
		k.producer,

		k.pipelines,
	}, nil
}

type kafkaExporter struct {
	chainConfig ChainConfig
	feedConfig  FeedConfig

	log      Logger
	producer Producer

	pipelines []Pipeline
}

func (k *kafkaExporter) Export(_ context.Context, data interface{}) {
	envelope, isEnvelope := data.(Envelope)
	if !isEnvelope {
		return
	}
	key := k.feedConfig.GetContractAddressBytes()

	var subs utils.Subprocesses
	defer subs.Wait()
	for _, pipeline := range k.pipelines {
		pipeline := pipeline
		subs.Go(func() {
			envelopeMapping, err := pipeline.Mapper(envelope, k.chainConfig, k.feedConfig)
			if err != nil {
				k.log.Errorw("failed to map envelope", "error", err, "topic", pipeline.Topic)
				return
			}
			encoded, err := pipeline.Schema.Encode(envelopeMapping)
			if err != nil {
				k.log.Errorw("failed to encode envelope to Avro", "payload", envelopeMapping, "error", err, "topic", pipeline.Topic)
				return
			}
			if err := k.producer.Produce(key, encoded, pipeline.Topic); err != nil {
				k.log.Errorw("failed to publish encoded payload to Kafka", "payload", envelopeMapping, "error", err)
				return
			}
		})
	}
}

func (k *kafkaExporter) Cleanup(_ context.Context) {} // noop
