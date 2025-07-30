package types

type WorkflowConfig struct {
	WriteTargetName       string `yaml:"writeTargetName"`
	DataFeedsCacheAddress string `yaml:"dataFeedsCacheAddress"`
	AllowedTriggerSender  string `yaml:"allowedTriggerSender"`
	AllowedTriggerTopic   string `yaml:"allowedTriggerTopic"`
	FeedID                string `yaml:"feedID"`
}
