package types

import "github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk"

type WorkflowConfig struct {
	// name of the secret that stores authentication key
	AuthKeySecretName string `yaml:"authKeySecretName"`
	ComputeConfig
}

type ComputeConfig struct {
	FeedID                string          `yaml:"feedID"`
	URL                   string          `yaml:"url"`
	DataFeedsCacheAddress string          `yaml:"dataFeedsCacheAddress"`
	WriteTargetName       string          `yaml:"writeTargetName"`
	AuthKey               sdk.SecretValue `yaml:",omitempty"`
}
