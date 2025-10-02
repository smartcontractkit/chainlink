package config

import "github.com/ethereum/go-ethereum/common"

type Config struct {
	FeedID         string
	ChainSelector  uint64
	FunctionToTest string
	InvalidInput   string
	DataFeedsCache
}

type DataFeedsCache struct {
	DataFeedsCacheAddress common.Address
}
