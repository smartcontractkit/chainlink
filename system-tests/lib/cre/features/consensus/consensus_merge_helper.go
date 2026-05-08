package consensus

import (
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

func WaitForLogPollerToBeHealthy(don *cre.Don) error {
	eg := &errgroup.Group{}
	for _, node := range don.Nodes {
		eg.Go(func() error {
			return node.Clients.RestClient.WaitHealthy(".*ConfigWatcher", "passing", 100)
		})
	}
	if waitErr := eg.Wait(); waitErr != nil {
		return errors.Wrap(waitErr, "failed to wait for ConfigWatcher health check")
	}

	return nil
}
