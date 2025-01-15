package ccip

type LoadConfig struct {
	LoadDuration              *string
	NoOfNodes                 *int
	LokiEndpoint              *string
	MessageTypeWeights        *[]int
	RequestFrequency          *string
	EnabledDestionationChains *[]uint64
	CribEnvDirectory          *string
}
