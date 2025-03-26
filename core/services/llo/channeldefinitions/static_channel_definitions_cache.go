package channeldefinitions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// A CDC that loads a static JSON of channel definitions; useful for
// benchmarking and testing

var _ llotypes.ChannelDefinitionCache = &staticCDC{}

type staticCDC struct {
	services.StateMachine
	lggr logger.Logger

	definitions llotypes.ChannelDefinitions
}

func NewStaticChannelDefinitionCache(lggr logger.Logger, dfnstr string) (llotypes.ChannelDefinitionCache, error) {
	var definitions llotypes.ChannelDefinitions
	if err := json.Unmarshal([]byte(dfnstr), &definitions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal static channel definitions: %w", err)
	}
	return &staticCDC{services.StateMachine{}, logger.Named(lggr, "StaticChannelDefinitionCache"), definitions}, nil
}

func (s *staticCDC) Start(context.Context) error {
	return s.StartOnce("StaticChannelDefinitionCache", func() error {
		s.lggr.Infow("Initialized static channel definition cache", "definitions", s.definitions)
		return nil
	})
}

func (s *staticCDC) Close() error {
	return s.StopOnce("StaticChannelDefinitionCache", func() error {
		return nil
	})
}

func (s *staticCDC) Definitions() llotypes.ChannelDefinitions {
	return s.definitions
}

func (s *staticCDC) HealthReport() map[string]error {
	report := map[string]error{s.Name(): s.Healthy()}
	return report
}

func (s *staticCDC) Name() string {
	return s.lggr.Name()
}
