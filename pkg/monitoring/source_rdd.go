package monitoring

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/smartcontractkit/chainlink-relay/pkg/utils"
	"go.uber.org/multierr"
)

type RDDData struct {
	Feeds []FeedConfig `json:"feeds,omitempty"`
	Nodes []NodeConfig `json:"nodes,omitempty"`
}

// rddSource produces a list of feeds to monitor.
// Any feed with the "status" field set to "dead" will be ignored and not returned by this source.
type rddSource struct {
	feedsURL       string
	feedsParser    FeedsParser
	feedsIgnoreIDs map[string]struct{}
	nodesURL       string
	nodesParser    NodesParser
	httpClient     *http.Client
	log            Logger
}

func NewRDDSource(
	feedsURL string,
	feedsParser FeedsParser,
	feedsIgnoreIDs []string,
	nodesURL string,
	nodesParser NodesParser,
	log Logger,
) Source {
	return &rddSource{
		feedsURL,
		feedsParser,
		makeSet(feedsIgnoreIDs),
		nodesURL,
		nodesParser,
		&http.Client{},
		log,
	}
}

func (r *rddSource) Fetch(ctx context.Context) (interface{}, error) {
	var subs utils.Subprocesses
	data := RDDData{}
	dataMu := &sync.Mutex{}
	var dataErr error
	subs.Go(func() {
		feeds, feedsErr := r.fetchFeeds(ctx)
		dataMu.Lock()
		defer dataMu.Unlock()
		if feedsErr != nil {
			dataErr = multierr.Combine(dataErr, feedsErr)
		} else {
			data.Feeds = feeds
		}
	})
	subs.Go(func() {
		nodes, nodesErr := r.fetchNodes(ctx)
		dataMu.Lock()
		defer dataMu.Unlock()
		if nodesErr != nil {
			dataErr = multierr.Combine(dataErr, nodesErr)
		} else {
			data.Nodes = nodes
		}
	})
	subs.Wait()
	return data, dataErr
}

func (r *rddSource) fetchFeeds(ctx context.Context) ([]FeedConfig, error) {
	readFeedsReq, err := http.NewRequestWithContext(ctx, http.MethodGet, r.feedsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to build a request to get feeds from the RDD: %w", err)
	}
	res, err := r.httpClient.Do(readFeedsReq)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch feeds RDD data: %w", err)
	}
	defer res.Body.Close()
	feeds, err := r.feedsParser(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to parse RDD feeds data: %w", err)
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

func (r *rddSource) fetchNodes(ctx context.Context) ([]NodeConfig, error) {
	readFeedsReq, err := http.NewRequestWithContext(ctx, http.MethodGet, r.nodesURL, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to build a request to get nodes from the RDD: %w", err)
	}
	res, err := r.httpClient.Do(readFeedsReq)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch nodes RDD data: %w", err)
	}
	defer res.Body.Close()
	nodes, err := r.nodesParser(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to parse RDD nodes data: %w", err)
	}
	return nodes, nil
}

// Helpers

func makeSet(ids []string) map[string]struct{} {
	out := make(map[string]struct{})
	for _, id := range ids {
		out[id] = struct{}{}
	}
	return out
}
