package s4

import (
	"context"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/s4"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"google.golang.org/protobuf/proto"
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

func NewReportingPlugin(logger logger.Logger, config *PluginConfig, orm s4.ORM) (types.ReportingPlugin, error) {
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

	versions, err := c.orm.GetVersions(c.addressRange, pg.WithParentCtx(ctx))
	if err != nil {
		return nil, errors.Wrap(err, "failed to GetVersions in Query()")
	}

	query := &Query{
		Versions: make([]*VersionRow, len(versions)),
		AddressRange: &AddressRange{
			MinAddress: MarshalAddress(c.addressRange.MinAddress),
			MaxAddress: MarshalAddress(c.addressRange.MaxAddress),
		},
	}

	for i := range versions {
		query.Versions[i] = convertVersionRow(versions[i])
	}

	queryBytes, err := proto.Marshal(query)
	if err != nil {
		return nil, err
	}

	promReportingPluginsQueryRowsCount.WithLabelValues(c.config.Product).Set(float64(len(versions)))
	promReportingPluginsQueryByteSize.WithLabelValues(c.config.Product).Set(float64(len(queryBytes)))

	c.addressRange.Advance()

	return queryBytes, err
}

func (c *plugin) Observation(ctx context.Context, _ types.ReportTimestamp, query types.Query) (types.Observation, error) {
	promReportingPluginObservation.WithLabelValues(c.config.Product).Inc()

	if err := c.orm.DeleteExpired(pg.WithParentCtx(ctx)); err != nil {
		return nil, errors.Wrap(err, "failed to DeleteExpired in Observation()")
	}

	observationRows := make([]*Row, 0)
	unconfirmedRows, err := c.orm.GetUnconfirmedRows(c.config.MaxObservationEntries, pg.WithParentCtx(ctx))
	if err != nil {
		return nil, errors.Wrap(err, "failed to GetUnconfirmedRows in Observation()")
	}

	if uint(len(unconfirmedRows)) < c.config.MaxObservationEntries {
		versionRows, addressRange, err := UnmarshalQuery(query)
		if err != nil || addressRange == nil {
			c.logger.Errorw("Failed to UnmarshalQuery, likely data is malformed", "err", err)
		} else {
			if c.addressRange.Interval().Cmp(addressRange.Interval()) != 0 {
				c.logger.Errorw("Address interval does not match, likely query is malformed", "current", c.addressRange.Interval(), "query", addressRange.Interval())
			} else {
				maxObservationRows := int(c.config.MaxObservationEntries) - len(unconfirmedRows)
				for _, vr := range versionRows {
					address, err := UnmarshalAddress(vr.Address)
					if err != nil {
						c.logger.Errorw("Failed to unmarshal address from Query", "err", err, "address", vr.Address)
						continue
					}
					row, err := c.orm.Get(address, uint(vr.Slotid), pg.WithParentCtx(ctx))
					if err == nil && row.Version > vr.Version {
						observationRows = append(observationRows, convertRow(row))
					} else if err != nil && !errors.Is(err, s4.ErrNotFound) {
						c.logger.Errorw("ORM Get error", "err", err)
					}
					if len(observationRows) >= maxObservationRows {
						break
					}
				}
			}
		}
	}

	observationEntriesCount := len(unconfirmedRows) + len(observationRows)
	if observationEntriesCount > int(c.config.MaxObservationEntries) {
		observationEntriesCount = int(c.config.MaxObservationEntries)
	}

	rows := make([]*Row, observationEntriesCount)
	for i := 0; i < observationEntriesCount; i++ {
		if i < len(unconfirmedRows) {
			rows[i] = convertRow(unconfirmedRows[i])
		} else {
			j := i - len(unconfirmedRows)
			rows[i] = observationRows[j]
		}
	}

	promReportingPluginsObservationRowsCount.WithLabelValues(c.config.Product).Set(float64(len(rows)))

	return MarshalRows(rows, nil)
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
			if err := verifySignature(row); err != nil {
				promReportingPluginWrongSigCount.WithLabelValues(c.config.Product).Inc()
				c.logger.Errorw("Report detected invalid signature", "err", err, "oracleID", ao.Observer)
				continue
			}
			mkey := key{
				address: row.Address,
				slot:    uint(row.Slotid),
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

func convertVersionRow(from *s4.VersionRow) *VersionRow {
	return &VersionRow{
		Address: from.Address.Hex(),
		Slotid:  uint32(from.SlotId),
		Version: from.Version,
	}
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
