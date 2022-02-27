package monitoring

import (
	"context"
	"fmt"
	"net/http"
)

// rddSource produces a list of feeds to monitor.
// Any feed with the "status" field set to "dead" will be ignored and not returned by this source.
type rddSource struct {
	rddURL     string
	httpClient *http.Client
	feedParser FeedParser
}

func NewRDDSource(
	rddURL string,
	feedParser FeedParser,
) Source {
	return &rddSource{
		rddURL,
		&http.Client{},
		feedParser,
	}
}

func (r *rddSource) Fetch(ctx context.Context) (interface{}, error) {
	readFeedsReq, err := http.NewRequestWithContext(ctx, http.MethodGet, r.rddURL, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to build a request to the RDD: %w", err)
	}
	res, err := r.httpClient.Do(readFeedsReq)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch RDD data: %w", err)
	}
	defer res.Body.Close()
	feeds, err := r.feedParser(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to parse RDD data into an array of feed configurations: %w", err)
	}
	feeds = removeDeadFeeds(feeds)
	return feeds, nil
}

func removeDeadFeeds(feeds []FeedConfig) []FeedConfig {
	out := []FeedConfig{}
	for _, feed := range feeds {
		if feed.GetContractStatus() != "dead" {
			out = append(out, feed)
		}
	}
	return out
}
