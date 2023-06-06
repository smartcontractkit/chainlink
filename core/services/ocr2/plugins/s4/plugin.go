package s4

import (
	"context"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/s4"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type plugin struct {
	logger       logger.Logger
	config       *PluginConfig
	orm          s4.ORM
	addressRange *s4.AddressRange
}

var _ types.ReportingPlugin = (*plugin)(nil)

func NewReportingPlugin(logger logger.Logger, config *PluginConfig, orm s4.ORM) (types.ReportingPlugin, error) {
	if config.MaxObservationEntries == 0 {
		return nil, errors.New("max number of observation entries cannot be zero")
	}
	if config.MaxReportEntries == 0 {
		return nil, errors.New("max number of report entries cannot be zero")
	}
	if config.MaxDeleteExpiredEntries == 0 {
		return nil, errors.New("max number of delete expired entries cannot be zero")
	}

	addressRange, err := s4.NewInitialAddressRangeForIntervals(config.NSnapshotShards)
	if err != nil {
		return nil, err
	}

	return &plugin{
		logger:       logger.Named("OCR2-S4").With("product", config.ProductName),
		config:       config,
		orm:          orm,
		addressRange: addressRange,
	}, nil
}

func (c *plugin) Query(ctx context.Context, _ types.ReportTimestamp) (types.Query, error) {
	promReportingPluginQuery.WithLabelValues(c.config.ProductName).Inc()

	snapshot, err := c.orm.GetSnapshot(c.addressRange, pg.WithParentCtx(ctx))
	if err != nil {
		return nil, errors.Wrap(err, "failed to GetVersions in Query()")
	}

	rows := make([]*SnapshotRow, len(snapshot))
	for i, v := range snapshot {
		rows[i] = &SnapshotRow{
			Address: v.Address.Bytes(),
			Slotid:  uint32(v.SlotId),
			Version: v.Version,
		}
	}

	queryBytes, err := MarshalQuery(rows)
	if err != nil {
		return nil, err
	}

	promReportingPluginsQueryRowsCount.WithLabelValues(c.config.ProductName).Set(float64(len(rows)))
	promReportingPluginsQueryByteSize.WithLabelValues(c.config.ProductName).Set(float64(len(queryBytes)))

	c.addressRange.Advance()

	return queryBytes, err
}

func (c *plugin) Observation(ctx context.Context, _ types.ReportTimestamp, query types.Query) (types.Observation, error) {
	promReportingPluginObservation.WithLabelValues(c.config.ProductName).Inc()

	now := time.Now().UTC()
	if err := c.orm.DeleteExpired(c.config.MaxDeleteExpiredEntries, now, pg.WithParentCtx(ctx)); err != nil {
		return nil, errors.Wrap(err, "failed to DeleteExpired in Observation()")
	}

	queryRows := make([]*s4.Row, 0)
	unconfirmedRows, err := c.orm.GetUnconfirmedRows(c.config.MaxObservationEntries, pg.WithParentCtx(ctx))
	if err != nil {
		return nil, errors.Wrap(err, "failed to GetUnconfirmedRows in Observation()")
	}

	if uint(len(unconfirmedRows)) < c.config.MaxObservationEntries {
		versionRows, err := UnmarshalQuery(query)
		if err != nil {
			c.logger.Errorw("Failed to UnmarshalQuery, likely data is malformed", "err", err)
		} else {
			maxObservationRows := int(c.config.MaxObservationEntries) - len(unconfirmedRows)
			for _, vr := range versionRows {
				address := UnmarshalAddress(vr.Address)
				row, err := c.orm.Get(address, uint(vr.Slotid), pg.WithParentCtx(ctx))
				if err == nil && row.Version > vr.Version {
					queryRows = append(queryRows, row)
				} else if err != nil && !errors.Is(err, s4.ErrNotFound) {
					c.logger.Errorw("ORM Get error", "err", err)
				}
				if len(queryRows) >= maxObservationRows {
					break
				}
			}
		}
	}

	rows := convertRows(append(unconfirmedRows, queryRows...))

	promReportingPluginsObservationRowsCount.WithLabelValues(c.config.ProductName).Set(float64(len(rows)))

	return MarshalRows(rows)
}

func (c *plugin) Report(_ context.Context, _ types.ReportTimestamp, _ types.Query, aos []types.AttributedObservation) (bool, types.Report, error) {
	promReportingPluginReport.WithLabelValues(c.config.ProductName).Inc()

	type key struct {
		address string
		slotID  uint
	}

	reportMap := make(map[key]*Row)

	for _, ao := range aos {
		observationRows, err := UnmarshalRows(ao.Observation)
		if err != nil {
			return false, nil, errors.Wrap(err, "failed to UnmarshalRows in Report()")
		}

		for _, row := range observationRows {
			if err := row.VerifySignature(); err != nil {
				promReportingPluginWrongSigCount.WithLabelValues(c.config.ProductName).Inc()
				c.logger.Errorw("Report detected invalid signature", "err", err, "oracleID", ao.Observer)
				continue
			}
			mkey := key{
				address: UnmarshalAddress(row.Address).String(),
				slotID:  uint(row.Slotid),
			}
			report, ok := reportMap[mkey]
			if ok && report.Version >= row.Version {
				continue
			}
			reportMap[mkey] = row
		}
	}

	reportRows := make([]*Row, 0)
	for _, row := range reportMap {
		reportRows = append(reportRows, row)

		if len(reportRows) >= int(c.config.MaxReportEntries) {
			break
		}
	}

	report, err := MarshalRows(reportRows)
	if err != nil {
		return false, nil, err
	}

	promReportingPluginsReportRowsCount.WithLabelValues(c.config.ProductName).Set(float64(len(reportRows)))

	return true, report, nil
}

func (c *plugin) ShouldAcceptFinalizedReport(ctx context.Context, _ types.ReportTimestamp, report types.Report) (bool, error) {
	promReportingPluginShouldAccept.WithLabelValues(c.config.ProductName).Inc()

	reportRows, err := UnmarshalRows(report)
	if err != nil {
		return false, errors.Wrap(err, "failed to UnmarshalRows in ShouldAcceptFinalizedReport()")
	}

	for _, row := range reportRows {
		ormRow := &s4.Row{
			Address:    UnmarshalAddress(row.Address),
			SlotId:     uint(row.Slotid),
			Payload:    row.Payload,
			Version:    row.Version,
			Expiration: row.Expiration,
			Confirmed:  true,
			Signature:  row.Signature,
		}
		err = c.orm.Update(ormRow, pg.WithParentCtx(ctx))
		if err != nil && !errors.Is(err, s4.ErrVersionTooLow) {
			c.logger.Errorw("Failed to Update a row in ShouldAcceptFinalizedReport()", "err", err)
			continue
		}
	}

	// If ShouldAcceptFinalizedReport returns false, ShouldTransmitAcceptedReport will not be called.
	return false, nil
}

func (c *plugin) ShouldTransmitAcceptedReport(context.Context, types.ReportTimestamp, types.Report) (bool, error) {
	return false, nil
}

func (c *plugin) Close() error {
	return nil
}

func convertRow(from *s4.Row) *Row {
	return &Row{
		Address:    from.Address.Bytes(),
		Slotid:     uint32(from.SlotId),
		Version:    from.Version,
		Expiration: from.Expiration,
		Payload:    from.Payload,
		Signature:  from.Signature,
	}
}

func convertRows(from []*s4.Row) []*Row {
	rows := make([]*Row, len(from))
	for i, row := range from {
		rows[i] = convertRow(row)
	}
	return rows
}
