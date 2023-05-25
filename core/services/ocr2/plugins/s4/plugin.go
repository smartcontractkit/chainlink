package s4

import (
	"context"
	"errors"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/s4"

	"github.com/ethereum/go-ethereum/common"
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
		logger:       logger.Named("OCR2-S4-plugin"),
		config:       config,
		orm:          orm,
		addressRange: addressRange,
	}, nil
}

func (c *plugin) Query(ctx context.Context, _ types.ReportTimestamp) (types.Query, error) {
	snapshot, err := c.orm.GetSnapshot(c.addressRange, pg.WithParentCtx(ctx))
	if err != nil {
		return nil, err
	}

	rows := make([]*Row, len(snapshot))
	for i, row := range snapshot {
		rows[i] = convertRow(row)
	}
	query, err := MarshalRows(rows, c.addressRange)
	if err != nil {
		return nil, err
	}

	c.addressRange.Advance()

	return query, err
}

func (c *plugin) Observation(ctx context.Context, _ types.ReportTimestamp, query types.Query) (types.Observation, error) {
	if err := c.orm.DeleteExpired(pg.WithParentCtx(ctx)); err != nil {
		return nil, err
	}

	queryRows, addressRange, err := UnmarshalRows(query)
	if err != nil {
		return nil, err
	}

	snapshot, err := c.orm.GetSnapshot(addressRange, pg.WithParentCtx(ctx))
	if err != nil {
		return nil, err
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

	return MarshalRows(observation, addressRange)
}

func (c *plugin) Report(_ context.Context, _ types.ReportTimestamp, _ types.Query, aos []types.AttributedObservation) (bool, types.Report, error) {
	reportMap := make(map[key]*Row)

	for _, ao := range aos {
		observationRows, _, err := UnmarshalRows(ao.Observation)
		if err != nil {
			return false, nil, err
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

	return true, report, nil
}

func (c *plugin) ShouldAcceptFinalizedReport(ctx context.Context, _ types.ReportTimestamp, report types.Report) (bool, error) {
	reportRows, _, err := UnmarshalRows(report)
	if err != nil {
		return false, err
	}

	for _, row := range reportRows {
		ormRow := &s4.Row{
			Address:    row.Address,
			SlotId:     uint(row.Slotid),
			Payload:    row.Payload,
			Version:    row.Version,
			Expiration: row.Expiration,
			Confirmed:  true,
			Signature:  row.Signature,
		}
		err := c.orm.Update(ormRow, pg.WithParentCtx(ctx))
		if err != nil && !errors.Is(err, s4.ErrVersionTooLow) {
			c.logger.Errorw("ORM error while updating row", "err", err)
			continue
		}
	}

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
		Address:    from.Address,
		Slotid:     uint32(from.SlotId),
		Version:    from.Version,
		Expiration: from.Expiration,
		Payload:    from.Payload,
		Signature:  from.Signature,
	}
}

func verifySignature(row *Row) error {
	rowAddress := common.HexToAddress(row.Address)
	e := &s4.Envelope{
		Address:    rowAddress.Bytes(),
		SlotID:     uint(row.Slotid),
		Payload:    row.Payload,
		Version:    row.Version,
		Expiration: row.Expiration,
	}
	signer, err := e.GetSignerAddress(row.Signature)
	if err != nil {
		return err
	}
	if signer != rowAddress {
		return s4.ErrWrongSignature
	}
	return nil
}
