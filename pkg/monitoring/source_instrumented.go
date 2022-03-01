package monitoring

import (
	"context"
	"time"
)

func NewInstrumentedSourceFactory(name string, sourceFactory SourceFactory, chainMetrics ChainMetrics) SourceFactory {
	return &instrumentedSourceFactory{name, sourceFactory, chainMetrics}
}

type instrumentedSourceFactory struct {
	name          string
	sourceFactory SourceFactory
	chainMetrics  ChainMetrics
}

func (i *instrumentedSourceFactory) NewSource(chainConfig ChainConfig, feedConfig FeedConfig) (Source, error) {
	source, err := i.sourceFactory.NewSource(chainConfig, feedConfig)
	if err != nil {
		return nil, err
	}
	return &instrumentedSource{
		i.name,
		source,
		NewFeedMetrics(chainConfig, feedConfig),
	}, nil
}

type instrumentedSource struct {
	name        string
	source      Source
	feedMetrics FeedMetrics
}

func (i *instrumentedSource) Fetch(ctx context.Context) (interface{}, error) {
	fetchStart := time.Now()
	data, err := i.source.Fetch(ctx)
	i.feedMetrics.ObserveFetchFromSourceDuraction(time.Since(fetchStart), i.name)
	if err != nil {
		i.feedMetrics.IncFetchFromSourceFailed(i.name)
	} else {
		i.feedMetrics.IncFetchFromSourceSucceeded(i.name)
	}
	return data, err
}
