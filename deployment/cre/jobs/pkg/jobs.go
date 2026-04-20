package pkg

import (
	"context"
	"errors"
	"fmt"
	"sync"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	nodeapiv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"

	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
)

type ProposeJobRequest struct {
	DONName   string
	Domain    string
	Spec      string // toml
	Env       string // staging, testnet, mainnet, etc...
	JobLabels map[string]string
	DONFilter *nodeapiv1.ListNodesRequest_Filter
}

// ProposeJob sends a single job spec to all the nodes in the DON indicated by `req.DONFilter`.
func ProposeJob(ctx context.Context, e cldf.Environment, req ProposeJobRequest) (map[string][]string, error) {
	nodes, err := offchain.FetchNodesFromJD(ctx, e.Offchain, req.DONFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to get don nodes: %w", err)
	}
	if len(nodes) == 0 {
		return nil, fmt.Errorf("no nodes found for DON `%s` with filters %+v", req.DONName, req.DONFilter)
	}

	jobSpecs := map[string][]string{}
	domain := offchain.ProductLabel
	if req.Domain != "" {
		domain = req.Domain
	}

	var (
		mu   sync.Mutex
		merr error
		wg   sync.WaitGroup
	)
	for _, node := range nodes {
		wg.Go(func() {
			e.Logger.Debugw("Proposing job", logLabels(req, node)...)
			offchainReq := offchain.ProposeJobRequest{
				Job:            req.Spec,
				Domain:         domain,
				Environment:    req.Env,
				PublicKeys:     []string{node.GetPublicKey()},
				JobLabels:      req.JobLabels,
				OffchainClient: e.Offchain,
				Lggr:           e.Logger,
				ExtraSelectors: req.DONFilter.GetSelectors(),
			}
			propErr := offchain.ProposeJob(ctx, offchainReq)

			mu.Lock()
			defer mu.Unlock()
			if propErr != nil {
				merr = errors.Join(merr, fmt.Errorf("failed to propose job: %w", propErr))
				return
			}
			jobSpecs[node.Id] = append(jobSpecs[node.Id], offchainReq.Job)
		})
	}
	wg.Wait()

	if merr != nil {
		return jobSpecs, merr
	}
	if len(jobSpecs) == 0 {
		return nil, errors.New("no jobs were proposed")
	}

	return jobSpecs, nil
}

func logLabels(req ProposeJobRequest, node *nodeapiv1.Node) []any {
	labels := []any{
		"nodeName",
		node.Name,
		"nodeID",
		node.Id,
		"publicKey",
		node.PublicKey,
		"target DON",
		req.DONName,
	}
	for k, v := range req.JobLabels {
		labels = append(labels, k, v)
	}

	return labels
}
