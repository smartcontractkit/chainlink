package automation

import "errors"

type Config struct {
	Performance *Performance `toml:"Performance"`
}

// Performance is a common configuration for all automation performance tests
type Performance struct {
	NumberOfNodes         *int      `toml:"number_of_nodes"`
	NumberOfUpkeeps       *int      `toml:"number_of_upkeeps"`
	Duration              *int      `toml:"duration"`
	BlockTime             *int      `toml:"block_time"`
	NumberOfEvents        *int      `toml:"number_of_events"`
	SpecType              *string   `toml:"spec_type"`
	ChainlinkNodeLogLevel *string   `toml:"chainlink_node_log_level"`
	TestInputs            *[]string `toml:"test_inputs"` //is this still needed?
}

func (c *Config) ApplyOverrides(from *Config) error {
	if from == nil {
		return nil
	}
	if from.Performance == nil {
		return nil
	}
	if c.Performance == nil {
		c.Performance = from.Performance
		return nil
	}
	if from.Performance.NumberOfNodes != nil {
		c.Performance.NumberOfNodes = from.Performance.NumberOfNodes
	}
	if from.Performance.NumberOfUpkeeps != nil {
		c.Performance.NumberOfUpkeeps = from.Performance.NumberOfUpkeeps
	}
	if from.Performance.Duration != nil {
		c.Performance.Duration = from.Performance.Duration
	}
	if from.Performance.BlockTime != nil {
		c.Performance.BlockTime = from.Performance.BlockTime
	}
	if from.Performance.NumberOfEvents != nil {
		c.Performance.NumberOfEvents = from.Performance.NumberOfEvents
	}
	if from.Performance.SpecType != nil {
		c.Performance.SpecType = from.Performance.SpecType
	}
	if from.Performance.ChainlinkNodeLogLevel != nil {
		c.Performance.ChainlinkNodeLogLevel = from.Performance.ChainlinkNodeLogLevel
	}
	if from.Performance.TestInputs != nil {
		c.Performance.TestInputs = from.Performance.TestInputs
	}

	return nil
}

func (c *Config) Validate() error {
	if c.Performance == nil {
		return nil
	}
	if c.Performance.NumberOfNodes == nil || *c.Performance.NumberOfNodes < 1 {
		return errors.New("number_of_nodes must be set to a positive integer")
	}
	if c.Performance.NumberOfUpkeeps == nil || *c.Performance.NumberOfUpkeeps < 1 {
		return errors.New("number_of_upkeeps must be set to a positive integer")
	}
	if c.Performance.Duration == nil || *c.Performance.Duration < 1 {
		return errors.New("duration must be set to a positive integer")
	}
	if c.Performance.BlockTime == nil || *c.Performance.BlockTime < 1 {
		return errors.New("block_time must be set to a positive integer")
	}
	if c.Performance.NumberOfEvents == nil || *c.Performance.NumberOfEvents < 1 {
		return errors.New("number_of_events must be set to a positive integer")
	}
	if c.Performance.SpecType == nil {
		return errors.New("spec_type must be set")
	}
	if c.Performance.ChainlinkNodeLogLevel == nil {
		return errors.New("chainlink_node_log_level must be set")
	}

	return nil
}
