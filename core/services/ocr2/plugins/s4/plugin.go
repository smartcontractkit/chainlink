package s4

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/services/s4"
)

var (
	promStoragePluginUpdatesCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "storage_plugin_updates",
		Help: "Number of storage updates fetched from other nodes",
	}, []string{})

	promStorageTotalByteSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "storage_total_byte_size",
		Help: "Current byte size of data stored in S4",
	}, []string{})
)

type plugin struct {
	logger       commontypes.Logger
	config       *PluginConfig
	orm          s4.ORM
	addressRange *s4.AddressRange
}

type key struct {
	address string
	slotID  uint
}

var _ types.ReportingPlugin = (*plugin)(nil)

func NewReportingPlugin(logger commontypes.Logger, config *PluginConfig, orm s4.ORM) (types.ReportingPlugin, error) {
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
		logger:       logger,
		config:       config,
		orm:          orm,
		addressRange: addressRange,
	}, nil
}

func (c *plugin) Query(ctx context.Context, ts types.ReportTimestamp) (types.Query, error) {
	promReportingPluginQuery.WithLabelValues(c.config.ProductName).Inc()

	snapshot, err := c.orm.GetSnapshot(ctx, c.addressRange)
	if err != nil {
		return nil, errors.Wrap(err, "failed to GetVersions in Query()")
	}

	var storageTotalByteSize uint64
	rows := make([]*SnapshotRow, len(snapshot))
	for i, v := range snapshot {
		rows[i] = &SnapshotRow{
			Address: v.Address.Bytes(),
			Slotid:  uint32(v.SlotId),
			Version: v.Version,
		}

		storageTotalByteSize += v.PayloadSize
	}

	queryBytes, err := MarshalQuery(rows, c.addressRange)
	if err != nil {
		return nil, err
	}

	promReportingPluginsQueryRowsCount.WithLabelValues(c.config.ProductName).Set(float64(len(rows)))
	promReportingPluginsQueryByteSize.WithLabelValues(c.config.ProductName).Set(float64(len(queryBytes)))

	promStorageTotalByteSize.WithLabelValues().Set(float64(storageTotalByteSize))

	c.addressRange.Advance()

	c.logger.Debug("S4StorageReporting Query", commontypes.LogFields{
		"epoch":         ts.Epoch,
		"round":         ts.Round,
		"nSnapshotRows": len(rows),
	})

	return queryBytes, err
}

func (c *plugin) Observation(ctx context.Context, ts types.ReportTimestamp, query types.Query) (types.Observation, error) {
	promReportingPluginObservation.WithLabelValues(c.config.ProductName).Inc()

	now := time.Now().UTC()
	count, err := c.orm.DeleteExpired(ctx, c.config.MaxDeleteExpiredEntries, now)
	if err != nil {
		return nil, errors.Wrap(err, "failed to DeleteExpired in Observation()")
	}
	promReportingPluginsExpiredRows.WithLabelValues(c.config.ProductName).Add(float64(count))

	returnObservation := func(rows []*s4.Row) (types.Observation, error) {
		promReportingPluginsObservationRowsCount.WithLabelValues(c.config.ProductName).Set(float64(len(rows)))
		return MarshalRows(convertRows(rows))
	}

	unconfirmedRows, err := c.orm.GetUnconfirmedRows(ctx, c.config.MaxObservationEntries)
	if err != nil {
		return nil, errors.Wrap(err, "failed to GetUnconfirmedRows in Observation()")
	}

	if uint(len(unconfirmedRows)) >= c.config.MaxObservationEntries {
		return returnObservation(unconfirmedRows[:c.config.MaxObservationEntries])
	}

	maxRemainingRows := int(c.config.MaxObservationEntries) - len(unconfirmedRows)
	remainingRows := make([]*s4.Row, 0)

	queryRows, addressRange, err := UnmarshalQuery(query)
	if err != nil {
		c.logger.Error("Failed to unmarshal query (likely malformed)", commontypes.LogFields{"err": err})
	} else {
		snapshot, err := c.orm.GetSnapshot(ctx, addressRange)
		if err != nil {
			c.logger.Error("ORM GetSnapshot error", commontypes.LogFields{"err": err})
		} else {
			type rkey struct {
				address *big.Big
				slotID  uint
			}

			snapshotVersionsMap := snapshotToVersionMap(snapshot)
			toBeAdded := make([]rkey, 0)
			// Add rows from query snapshot that have a higher version locally.
			for _, qr := range queryRows {
				address := UnmarshalAddress(qr.Address)
				k := key{address: address.String(), slotID: uint(qr.Slotid)}
				if version, ok := snapshotVersionsMap[k]; ok && version > qr.Version {
					toBeAdded = append(toBeAdded, rkey{address: address, slotID: uint(qr.Slotid)})
				}
				delete(snapshotVersionsMap, k)
			}

			if len(toBeAdded) > maxRemainingRows {
				toBeAdded = toBeAdded[:maxRemainingRows]
			} else {
				// Add rows from query address range that exist locally but are missing from query snapshot.
				for _, sr := range snapshot {
					if !sr.Confirmed {
						continue
					}
					k := key{address: sr.Address.String(), slotID: sr.SlotId}
					if _, ok := snapshotVersionsMap[k]; ok {
						toBeAdded = append(toBeAdded, rkey{address: sr.Address, slotID: sr.SlotId})
						if len(toBeAdded) == maxRemainingRows {
							break
						}
					}
				}
			}

			for _, k := range toBeAdded {
				row, err := c.orm.Get(ctx, k.address, k.slotID)
				if err == nil {
					remainingRows = append(remainingRows, row)
				} else if !errors.Is(err, s4.ErrNotFound) {
					c.logger.Error("ORM Get error", commontypes.LogFields{"err": err})
				}
			}
		}
	}

	c.logger.Debug("S4StorageReporting Observation", commontypes.LogFields{
		"epoch":            ts.Epoch,
		"round":            ts.Round,
		"nUnconfirmedRows": len(unconfirmedRows),
		"nRemainingRows":   len(remainingRows),
	})

	return returnObservation(append(unconfirmedRows, remainingRows...))
}

func (c *plugin) Report(_ context.Context, ts types.ReportTimestamp, _ types.Query, aos []types.AttributedObservation) (bool, types.Report, error) {
	promReportingPluginReport.WithLabelValues(c.config.ProductName).Inc()

	reportMap := make(map[key]*Row)
	reportKeys := []key{}

	for _, ao := range aos {
		observationRows, err := UnmarshalRows(ao.Observation)
		if err != nil {
			return false, nil, errors.Wrap(err, "failed to UnmarshalRows in Report()")
		}

		for _, row := range observationRows {
			if err := row.VerifySignature(); err != nil {
				promReportingPluginWrongSigCount.WithLabelValues(c.config.ProductName).Inc()
				c.logger.Error("Report detected invalid signature", commontypes.LogFields{"err": err, "oracleID": ao.Observer})
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
			reportKeys = append(reportKeys, mkey)
		}
	}

	reportRows := make([]*Row, 0)
	for _, key := range reportKeys {
		row := reportMap[key]
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
	c.logger.Debug("S4StorageReporting Report", commontypes.LogFields{
		"epoch":         ts.Epoch,
		"round":         ts.Round,
		"nReportRows":   len(reportRows),
		"nObservations": len(aos),
	})

	return true, report, nil
}

func (c *plugin) ShouldAcceptFinalizedReport(ctx context.Context, ts types.ReportTimestamp, report types.Report) (bool, error) {
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

		now := time.Now().UnixMilli()
		if now > ormRow.Expiration {
			c.logger.Error("Received an expired entry in a report, not saving", commontypes.LogFields{
				"expirationTs": ormRow.Expiration,
				"nowTs":        now,
			})
			continue
		}

		err = c.orm.Update(ctx, ormRow)
		if err != nil && !errors.Is(err, s4.ErrVersionTooLow) {
			c.logger.Error("Failed to Update a row in ShouldAcceptFinalizedReport()", commontypes.LogFields{"err": err})
			continue
		}
		promStoragePluginUpdatesCount.WithLabelValues().Inc()
	}

	c.logger.Debug("S4StorageReporting ShouldAcceptFinalizedReport", commontypes.LogFields{
		"epoch":       ts.Epoch,
		"round":       ts.Round,
		"nReportRows": len(reportRows),
	})

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

func snapshotToVersionMap(rows []*s4.SnapshotRow) map[key]uint64 {
	m := make(map[key]uint64)
	for _, row := range rows {
		if row.Confirmed {
			m[key{address: row.Address.String(), slotID: row.SlotId}] = row.Version
		}
	}
	return m
}
