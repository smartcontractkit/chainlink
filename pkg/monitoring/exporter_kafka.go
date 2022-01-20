package monitoring

import (
	"context"
	"sync"
)

func NewKafkaExporter(
	chainConfig ChainConfig,
	feedConfig FeedConfig,

	log Logger,
	producer Producer,

	transmissionSchema Schema,
	configSetSimplifiedSchema Schema,

	transmissionTopic string,
	configSetSimplifiedTopic string,
) Exporter {
	return &kafkaExporter{
		chainConfig,
		feedConfig,

		log,
		producer,

		transmissionSchema,
		configSetSimplifiedSchema,

		transmissionTopic,
		configSetSimplifiedTopic,
	}
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
	wg.Add(2)
	defer wg.Wait()
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
}

func (k *kafkaExporter) Cleanup() {} // noop
