package s4

import (
	"context"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/s4"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type plugin struct {
	logger       logger.Logger
	config       *PluginConfig
	orm          s4.ORM
	addressRange *s4.AddressRange
}

type key struct {
	address string
	slot    uint
}

var _ types.ReportingPlugin = (*plugin)(nil)

func NewReportingPlugin(logger logger.Logger, config *PluginConfig, orm s4.ORM) (*plugin, error) {
	if config.MaxObservationEntries == 0 {
		return nil, errors.New("max number of observation entries cannot be zero")
	}

	addressRange, err := s4.NewInitialAddressRangeForIntervals(config.NSnapshotShards)
	if err != nil {
		return nil, err
	}

	return &plugin{
		logger:       logger.Named("OCR2-S4").With("product", config.Product),
		config:       config,
		orm:          orm,
		addressRange: addressRange,
	}, nil
}

func (c *plugin) Query(ctx context.Context, _ types.ReportTimestamp) (types.Query, error) {
	promReportingPluginQuery.WithLabelValues(c.config.Product).Inc()

	snapshot, err := c.orm.GetSnapshot(c.addressRange, pg.WithParentCtx(ctx))
	if err != nil {
		return nil, errors.Wrap(err, "failed to GetSnapshot in Query()")
	}

	rows := make([]*Row, len(snapshot))
	for i, row := range snapshot {
		rows[i] = convertRow(row)
	}
	query, err := MarshalRows(rows, c.addressRange)
	if err != nil {
		return nil, err
	}

	promReportingPluginsQueryRowsCount.WithLabelValues(c.config.Product).Set(float64(len(rows)))
	promReportingPluginsQueryByteSize.WithLabelValues(c.config.Product).Set(float64(len(query)))

	c.addressRange.Advance()

	return query, err
}

func (c *plugin) Observation(ctx context.Context, _ types.ReportTimestamp, query types.Query) (types.Observation, error) {
	promReportingPluginObservation.WithLabelValues(c.config.Product).Inc()

	if err := c.orm.DeleteExpired(pg.WithParentCtx(ctx)); err != nil {
		return nil, errors.Wrap(err, "failed to DeleteExpired in Observation()")
	}

	queryRows, addressRange, err := UnmarshalRows(query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to UnmarshalRows in Observation()")
	}

	snapshot, err := c.orm.GetSnapshot(addressRange, pg.WithParentCtx(ctx))
	if err != nil {
		return nil, errors.Wrap(err, "failed to GetSnapshot in Observation()")
	}

	snapshotMap := make(map[key]*Row)
	unconfirmedMap := make(map[key]*Row)
	for _, row := range snapshot {
		r := convertRow(row)
		mkey := key{
			address: r.Address,
			slot:    uint(r.Slotid),
		}
		snapshotMap[mkey] = r
		if !row.Confirmed {
			unconfirmedMap[mkey] = r
		}
	}

	observation := make([]*Row, 0)
	for _, queryRow := range queryRows {
		mkey := key{
			address: queryRow.Address,
			slot:    uint(queryRow.Slotid),
		}
		snapshotRow, ok := snapshotMap[mkey]
		if ok && queryRow.Version < snapshotRow.Version {
			observation = append(observation, queryRow)
			delete(unconfirmedMap, mkey)
		}
	}

	for _, unconfirmed := range unconfirmedMap {
		observation = append(observation, unconfirmed)
	}
	if len(observation) > int(c.config.MaxObservationEntries) {
		observation = observation[:c.config.MaxObservationEntries]
	}

	promReportingPluginsObservationRowsCount.WithLabelValues(c.config.Product).Set(float64(len(observation)))

	return MarshalRows(observation, addressRange)
}

func (c *plugin) Report(_ context.Context, _ types.ReportTimestamp, _ types.Query, aos []types.AttributedObservation) (bool, types.Report, error) {
	promReportingPluginReport.WithLabelValues(c.config.Product).Inc()

	reportMap := make(map[key]*Row)

	for _, ao := range aos {
		observationRows, _, err := UnmarshalRows(ao.Observation)
		if err != nil {
			return false, nil, errors.Wrap(err, "failed to UnmarshalRows in Report()")
		}

		for _, row := range observationRows {
			mkey := key{
				address: row.Address,
				slot:    uint(row.Slotid),
			}
			report, ok := reportMap[mkey]
			if ok && report.Version >= row.Version {
				continue
			}
			if err := verifySignature(row); err != nil {
				promReportingPluginWrongSigCount.WithLabelValues(c.config.Product).Inc()
				c.logger.Errorw("Report round detected invalid signature", "err", err, "oracleID", ao.Observer)
				continue
			}
			reportMap[mkey] = row
		}
	}

	reportRows := make([]*Row, 0)
	for _, row := range reportMap {
		reportRows = append(reportRows, row)
	}

	report, err := MarshalRows(reportRows, nil)
	if err != nil {
		return false, nil, err
	}

	promReportingPluginsReportRowsCount.WithLabelValues(c.config.Product).Set(float64(len(reportRows)))

	return true, report, nil
}

func (c *plugin) ShouldAcceptFinalizedReport(ctx context.Context, _ types.ReportTimestamp, report types.Report) (bool, error) {
	promReportingPluginShouldAccept.WithLabelValues(c.config.Product).Inc()

	reportRows, _, err := UnmarshalRows(report)
	if err != nil {
		return false, errors.Wrap(err, "failed to UnmarshalRows in ShouldAcceptFinalizedReport()")
	}

	for _, row := range reportRows {
		address, err := UnmarshalAddress(row.Address)
		if err != nil {
			return false, err
		}
		ormRow := &s4.Row{
			Address:    address,
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
		Address:    from.Address.Hex(),
		Slotid:     uint32(from.SlotId),
		Version:    from.Version,
		Expiration: from.Expiration,
		Payload:    from.Payload,
		Signature:  from.Signature,
	}
}

func verifySignature(row *Row) error {
	address := common.HexToAddress(row.Address)
	e := &s4.Envelope{
		Address:    address.Bytes(),
		SlotID:     uint(row.Slotid),
		Payload:    row.Payload,
		Version:    row.Version,
		Expiration: row.Expiration,
	}
	signer, err := e.GetSignerAddress(row.Signature)
	if err != nil {
		return err
	}
	if signer != address {
		return s4.ErrWrongSignature
	}
	return nil
}
