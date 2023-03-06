package networking

type BinaryNetworkEndpointLimits struct {
	MaxMessageLength          int
	MessagesRatePerOracle     float64
	MessagesCapacityPerOracle int
	BytesRatePerOracle        float64
	BytesCapacityPerOracle    int
}
