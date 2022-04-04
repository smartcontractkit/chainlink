package monitoring

import (
	"context"
	"fmt"
	"net/http"
)

// rddSource produces a list of feeds to monitor.
// Any feed with the "status" field set to "dead" will be ignored and not returned by this source.
type rddSource struct {
	rddURL         string
	httpClient     *http.Client
	feedParser     FeedParser
	log            Logger
	feedsIgnoreIDs map[string]struct{}
}

func NewRDDSource(
	rddURL string,
	feedParser FeedParser,
	log Logger,
	feedsIgnoreIDs []string,
) Source {
	return &rddSource{
		rddURL,
		&http.Client{},
		feedParser,
		log,
		makeSet(feedsIgnoreIDs),
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
	feeds = r.filterFeeds(feeds)
	return feeds, nil
}

// filterFeeds removes feeds that:
// - have status=="dead"
// - have their ID specified in FEEDS_IGNORE_IDS env var.
func (r *rddSource) filterFeeds(feeds []FeedConfig) []FeedConfig {
	out := []FeedConfig{}
	for _, feed := range feeds {
		if feed.GetContractStatus() == "dead" {
			r.log.Infow("ignoring feed because of contract_status=dead", "feed_id", feed.GetID())
			continue
		}
		if _, isIgnored := r.feedsIgnoreIDs[feed.GetID()]; isIgnored {
			r.log.Debugw("skipping feed because of it is marked as ignored in the FEEDS_IGNORE_IDS env var", "feed_id", feed.GetID())
			continue
		}
		out = append(out, feed)
	}
	return out
}

// Helpers

func makeSet(ids []string) map[string]struct{} {
	out := make(map[string]struct{})
	for _, id := range ids {
		out[id] = struct{}{}
	}
	return out
}
