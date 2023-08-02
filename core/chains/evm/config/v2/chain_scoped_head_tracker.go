package v2

import "time"

type headTrackerConfig struct {
	c HeadTracker
}

func (h *headTrackerConfig) HistoryDepth() uint32 {
	return *h.c.HistoryDepth
}

func (h *headTrackerConfig) MaxBufferSize() uint32 {
	return *h.c.MaxBufferSize
}

func (h *headTrackerConfig) SamplingInterval() time.Duration {
	return h.c.SamplingInterval.Duration()
}
