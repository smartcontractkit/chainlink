package automation

type Config struct {
	Common *Common `toml:"Common"`
}

// Common is a common configuration for all automation performance tests
// TODO maybe we should put it under "Performance" as an umbrella term, as opposed to "Smoke"
type Common struct {
	NumberOfNodes         *int      `toml:"number_of_nodes"`
	NumberOfUpkeeps       *int      `toml:"number_of_upkeeps"`
	Duration              *int      `toml:"duration"`
	BlockTime             *int      `toml:"block_time"`
	NumberOfEvents        *int      `toml:"number_of_events"`
	SpecType              *string   `toml:"spec_type"`
	ChainlinkNodeLogLevel *string   `toml:"chainlink_node_log_level"`
	TestInputs            *[]string `toml:"test_inputs"`
}

func (c *Config) ApplyOverrides(from *Config) error {
	if from == nil {
		return nil
	}
	if from.Common == nil {
		return nil
	}

	if c.Common == nil {
		c.Common = from.Common
		return nil
	}

	if from.Common.NumberOfNodes != nil {
		c.Common.NumberOfNodes = from.Common.NumberOfNodes
	}
	if from.Common.NumberOfUpkeeps != nil {
		c.Common.NumberOfUpkeeps = from.Common.NumberOfUpkeeps
	}
	if from.Common.Duration != nil {
		c.Common.Duration = from.Common.Duration
	}
	if from.Common.BlockTime != nil {
		c.Common.BlockTime = from.Common.BlockTime
	}
	if from.Common.NumberOfEvents != nil {
		c.Common.NumberOfEvents = from.Common.NumberOfEvents
	}
	if from.Common.SpecType != nil {
		c.Common.SpecType = from.Common.SpecType
	}
	if from.Common.ChainlinkNodeLogLevel != nil {
		c.Common.ChainlinkNodeLogLevel = from.Common.ChainlinkNodeLogLevel
	}
	if from.Common.TestInputs != nil {
		c.Common.TestInputs = from.Common.TestInputs
	}

	return nil
}

func (c *Config) Validate() error {
	//TODO implement me
	return nil
}
