package monitoring

func NewInstrumentedProducer(producer Producer, chainMetrics ChainMetrics) Producer {
	return &instrumentedProducer{producer, chainMetrics}
}

type instrumentedProducer struct {
	producer     Producer
	chainMetrics ChainMetrics
}

func (i *instrumentedProducer) Produce(key, value []byte, topic string) error {
	err := i.producer.Produce(key, value, topic)
	if err != nil {
		i.chainMetrics.IncSendMessageToKafkaFailed(topic)
	} else {
		i.chainMetrics.IncSendMessageToKafkaSucceeded(topic)
	}
	return err
}
