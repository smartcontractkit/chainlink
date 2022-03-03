package monitoring

import (
	"context"
	"time"
)

func NewInstrumentedSourceFactory(sourceFactory SourceFactory, chainMetrics ChainMetrics) SourceFactory {
	return &instrumentedSourceFactory{sourceFactory, chainMetrics}
}

type instrumentedSourceFactory struct {
	sourceFactory SourceFactory
	chainMetrics  ChainMetrics
}

func (i *instrumentedSourceFactory) NewSource(chainConfig ChainConfig, feedConfig FeedConfig) (Source, error) {
	source, err := i.sourceFactory.NewSource(chainConfig, feedConfig)
	if err != nil {
		return nil, err
	}
	return &instrumentedSource{
		i.sourceFactory.GetType(),
		source,
		NewFeedMetrics(chainConfig, feedConfig),
	}, nil
}

func (i *instrumentedSourceFactory) GetType() string {
	return i.sourceFactory.GetType()
}

type instrumentedSource struct {
	sourceType  string
	source      Source
	feedMetrics FeedMetrics
}

func (i *instrumentedSource) Fetch(ctx context.Context) (interface{}, error) {
	fetchStart := time.Now()
	data, err := i.source.Fetch(ctx)
	i.feedMetrics.ObserveFetchFromSourceDuraction(time.Since(fetchStart), i.sourceType)
	if err != nil {
		i.feedMetrics.IncFetchFromSourceFailed(i.sourceType)
	} else {
		i.feedMetrics.IncFetchFromSourceSucceeded(i.sourceType)
	}
	return data, err
}
