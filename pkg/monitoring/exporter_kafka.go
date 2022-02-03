package monitoring

import (
	"context"
	"fmt"
	"sync"
)

type Pipeline struct {
	Topic  string
	Mapper Mapper
	Schema Schema
}

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
	chainConfig ChainConfig,
	feedConfig FeedConfig,
) (Exporter, error) {
	return &kafkaExporter{
		chainConfig,
		feedConfig,

		k.log.With("feed", feedConfig.GetName()),
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

	wg := &sync.WaitGroup{}
	defer wg.Wait()
	wg.Add(len(k.pipelines))
	for _, pipeline := range k.pipelines {
		go func(pipeline Pipeline) {
			defer wg.Done()
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
		}(pipeline)
	}
}

func (k *kafkaExporter) Cleanup(_ context.Context) {} // noop
