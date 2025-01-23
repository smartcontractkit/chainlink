package pg

import (
	commonpg "github.com/smartcontractkit/chainlink-common/pkg/sqlutil/pg"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type (
	StatFn   = commonpg.StatFn
	ReportFn = commonpg.ReportFn
)

type StatsReporter = commonpg.StatsReporter

func NewStatsReporter(fn StatFn, lggr logger.Logger, opts ...commonpg.StatsReporterOpt) *StatsReporter {
	return commonpg.NewStatsReporter(fn, lggr, opts...)
}
