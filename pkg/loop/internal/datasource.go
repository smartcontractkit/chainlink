package internal

import (
	"context"
	"math/big"
	"os"
	"time"

	"google.golang.org/grpc"

	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
)

// github.com/smartcontractkit/libocr/offchainreporting2plus/internal/protocol.ReportingPluginTimeoutWarningGracePeriod
var datasourceOvertime = 100 * time.Millisecond

func init() {
	// undocumented escape hatch
	// TODO: remove with https://smartcontract-it.atlassian.net/browse/BCF-2209
	if v := os.Getenv("CL_DATASOURCE_OVERTIME"); v != "" {
		d, err := time.ParseDuration(v)
		if err == nil {
			datasourceOvertime = d
		}
	}
}

var _ median.DataSource = (*dataSourceClient)(nil)

type dataSourceClient struct {
	grpc pb.DataSourceClient
}

func newDataSourceClient(cc grpc.ClientConnInterface) *dataSourceClient {
	return &dataSourceClient{grpc: pb.NewDataSourceClient(cc)}
}

func (d *dataSourceClient) Observe(ctx context.Context, timestamp types.ReportTimestamp) (*big.Int, error) {
	reply, err := d.grpc.Observe(ctx, &pb.ObserveRequest{
		ReportTimestamp: pbReportTimestamp(timestamp),
	})
	if err != nil {
		return nil, err
	}
	return reply.Value.Int(), nil
}

var _ pb.DataSourceServer = (*dataSourceServer)(nil)

type dataSourceServer struct {
	pb.UnimplementedDataSourceServer

	impl median.DataSource
}

func (d *dataSourceServer) Observe(ctx context.Context, request *pb.ObserveRequest) (*pb.ObserveReply, error) {
	// Pipeline observations may return results after the context is cancelled, so we modify the
	// deadline to give them time to return before the parent context deadline.
	// TODO: remove with https://smartcontract-it.atlassian.net/browse/BCF-2209
	var cancel func()
	ctx, cancel = utils.ContextWithDeadlineFn(ctx, func(orig time.Time) time.Time {
		if tenPct := time.Until(orig) / 10; datasourceOvertime > tenPct {
			return orig.Add(-tenPct)
		}
		return orig.Add(-datasourceOvertime)
	})
	defer cancel()
	timestamp, err := reportTimestamp(request.ReportTimestamp)
	if err != nil {
		return nil, err
	}
	val, err := d.impl.Observe(ctx, timestamp)
	if err != nil {
		return nil, err
	}
	return &pb.ObserveReply{Value: pb.NewBigIntFromInt(val)}, nil
}
