package types

type WorkflowConfig struct {
	WriteTargetName       string `yaml:"write_target_name"`
	ChainFamily           string `yaml:"chain_family,omitempty"`
	ChainID               string `yaml:"chain_id,omitempty"`
	DataFeedsCacheAddress string `yaml:"data_feeds_cache_address"`
	AllowedTriggerSender  string `yaml:"allowed_trigger_sender"`
	AllowedTriggerTopic   string `yaml:"allowed_trigger_topic"`
	FeedID                string `yaml:"feed_id"`
	BalanceReaderConfig
}

type BalanceReaderConfig struct {
	BalanceReaderAddress string `yaml:"balance_reader_address"`
}
