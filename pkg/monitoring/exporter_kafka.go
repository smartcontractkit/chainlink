package monitoring

import (
	"context"
	"sync"
)

func NewKafkaExporterFactory(
	log Logger,
	producer Producer,

	transmissionSchema Schema,
	configSetSimplifiedSchema Schema,

	transmissionTopic string,
	configSetSimplifiedTopic string,
) ExporterFactory {
	return &kafkaExporterFactory{
		log,
		producer,

		transmissionSchema,
		configSetSimplifiedSchema,

		transmissionTopic,
		configSetSimplifiedTopic,
	}
}

type kafkaExporterFactory struct {
	log      Logger
	producer Producer

	transmissionSchema        Schema
	configSetSimplifiedSchema Schema

	transmissionTopic        string
	configSetSimplifiedTopic string
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

		k.transmissionSchema,
		k.configSetSimplifiedSchema,

		k.transmissionTopic,
		k.configSetSimplifiedTopic,
	}, nil
}

type kafkaExporter struct {
	chainConfig ChainConfig
	feedConfig  FeedConfig

	log      Logger
	producer Producer

	transmissionSchema        Schema
	configSetSimplifiedSchema Schema

	transmissionTopic        string
	configSetSimplifiedTopic string
}

func (k *kafkaExporter) Export(ctx context.Context, data interface{}) {
	key := k.feedConfig.GetContractAddressBytes()
	envelope, ok := data.(Envelope)
	if !ok {
		k.log.Errorw("expected payload of type Envelope but got %#v", data)
		return
	}
	wg := &sync.WaitGroup{}
	defer wg.Wait()
	wg.Add(2)
	go func(key []byte, envelope Envelope) {
		defer wg.Done()
		transmissionMapping, err := MakeTransmissionMapping(envelope, k.chainConfig, k.feedConfig)
		if err != nil {
			k.log.Errorw("failed to map transmission", "error", err)
			return
		}
		transmissionEncoded, err := k.transmissionSchema.Encode(transmissionMapping)
		if err != nil {
			k.log.Errorw("failed to encode transmission to Avro", "payload", transmissionMapping, "error", err)
			return
		}
		if err := k.producer.Produce(key, transmissionEncoded, k.transmissionTopic); err != nil {
			k.log.Errorw("failed to publish transmission", "payload", transmissionMapping, "error", err)
			return
		}
	}(key, envelope)
	go func(key []byte, envelope Envelope) {
		defer wg.Done()
		configSetSimplifiedMapping, err := MakeConfigSetSimplifiedMapping(envelope, k.feedConfig)
		if err != nil {
			k.log.Errorw("failed to map config_set_simplified", "error", err)
			return
		}
		configSetSimplifiedEncoded, err := k.configSetSimplifiedSchema.Encode(configSetSimplifiedMapping)
		if err != nil {
			k.log.Errorw("failed to encode config_set_simplified to Avro", "payload", configSetSimplifiedMapping, "error", err)
			return
		}
		if err := k.producer.Produce(key, configSetSimplifiedEncoded, k.configSetSimplifiedTopic); err != nil {
			k.log.Errorw("failed to publish config_set_simplified", "payload", configSetSimplifiedMapping, "error", err)
			return
		}
	}(key, envelope)
}

func (k *kafkaExporter) Cleanup(_ context.Context) {} // noop
