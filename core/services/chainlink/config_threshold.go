package chainlink

import v2 "github.com/smartcontractkit/chainlink/v2/core/config/v2"

type thresholdConfig struct {
	s v2.ThresholdKeyShareSecrets
}

func (t *thresholdConfig) ThresholdKeyShare() string {
	if t.s.ThresholdKeyShare == nil {
		return ""
	}
	return string(*t.s.ThresholdKeyShare)
}
