package s4

import (
	"context"
	"errors"
	"math/big"

	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/s4"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"google.golang.org/protobuf/proto"
)

type plugin struct {
	config          *PluginConfig
	orm             s4.ORM
	minQueryAddress *big.Int
	maxQueryAddress *big.Int
	addressInteval  *big.Int
}

type key struct {
	address string
	slot    uint
}

var _ types.ReportingPlugin = (*plugin)(nil)

func NewConsensusPlugin(config *PluginConfig, orm s4.ORM) (*plugin, error) {
	if config.NSnapshotShards == 0 {
		return nil, errors.New("number of snapshots shards cannot be zero")
	}
	if config.MaxObservationEntries == 0 {
		return nil, errors.New("max number of observation entries cannot be zero")
	}
	if config.MaxReportEntries == 0 {
		return nil, errors.New("max number of report entries cannot be zero")
	}

	addressRange := new(big.Int).Sub(s4.MaxAddress, s4.MinAddress)
	divisor := big.NewInt(int64(config.NSnapshotShards))
	addressInteval := new(big.Int).Div(addressRange, divisor)

	return &plugin{
		config:          config,
		orm:             orm,
		minQueryAddress: s4.MinAddress,
		maxQueryAddress: new(big.Int).Add(s4.MinAddress, addressInteval),
		addressInteval:  addressInteval,
	}, nil
}

func (c *plugin) advanceAddressRange() {
	if c.config.NSnapshotShards > 1 {
		c.minQueryAddress = new(big.Int).Add(c.minQueryAddress, c.addressInteval)
		c.maxQueryAddress = new(big.Int).Add(c.maxQueryAddress, c.addressInteval)
		if c.maxQueryAddress.Cmp(s4.MaxAddress) > 0 {
			c.maxQueryAddress = s4.MaxAddress
		}
		if c.minQueryAddress.Cmp(s4.MaxAddress) >= 0 {
			c.minQueryAddress = s4.MinAddress
			c.maxQueryAddress = new(big.Int).Add(c.minQueryAddress, c.addressInteval)
		}
	}
}

func (c *plugin) Query(ctx context.Context, _ types.ReportTimestamp) (types.Query, error) {
	snapshot, err := c.orm.GetSnapshot(c.minQueryAddress, c.maxQueryAddress, pg.WithParentCtx(ctx))
	if err != nil {
		return nil, err
	}

	if len(snapshot) == 0 {
		return nil, nil
	}

	rows := make([]*Row, len(snapshot))
	for i, row := range snapshot {
		rows[i] = convertRow(row)
	}
	query, err := marshalRows(rows, c.minQueryAddress, c.maxQueryAddress)
	if err != nil {
		return nil, err
	}

	c.advanceAddressRange()

	return query, err
}

func (c *plugin) Observation(ctx context.Context, _ types.ReportTimestamp, query types.Query) (types.Observation, error) {
	if err := c.orm.DeleteExpired(pg.WithParentCtx(ctx)); err != nil {
		return nil, err
	}

	queryRows, minAddress, maxAddress, err := unmarshalRows(query)
	if err != nil {
		return nil, err
	}

	snapshot, err := c.orm.GetSnapshot(minAddress, maxAddress, pg.WithParentCtx(ctx))
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

	return marshalRows(observation, minAddress, maxAddress)
}

func (c *plugin) Report(_ context.Context, _ types.ReportTimestamp, _ types.Query, aos []types.AttributedObservation) (bool, types.Report, error) {
	reportMap := make(map[key]*Row)

	for _, ao := range aos {
		observationRows, _, _, err := unmarshalRows(ao.Observation)
		if err != nil {
			return false, nil, err
		}

		for _, row := range observationRows {
			mkey := key{
				address: row.Address,
				slot:    uint(row.Slotid),
			}
			report, ok := reportMap[mkey]
			if ok && report.Version > row.Version {
				continue
			}
			if err := verifySignature(row); err != nil {
				return false, nil, err
			}
			reportMap[mkey] = row
		}
	}

	reportRows := make([]*Row, 0)
	for _, row := range reportMap {
		reportRows = append(reportRows, row)
	}
	if len(reportRows) > int(c.config.MaxReportEntries) {
		reportRows = reportRows[:c.config.MaxReportEntries]
	}

	report, err := marshalRows(reportRows, nil, nil)
	if err != nil {
		return false, nil, err
	}

	return true, report, nil
}

func (c *plugin) ShouldAcceptFinalizedReport(ctx context.Context, _ types.ReportTimestamp, report types.Report) (bool, error) {
	reportRows, _, _, err := unmarshalRows(report)
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
			return false, err
		}
	}

	return true, nil
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

func marshalRows(rows []*Row, minAddress, maxAddress *big.Int) ([]byte, error) {
	rr := &Rows{
		Rows: rows,
	}
	if minAddress != nil {
		minAddressStr, err := minAddress.MarshalText()
		if err != nil {
			return nil, err
		}
		rr.MinAddress = string(minAddressStr)
	}
	if maxAddress != nil {
		maxAddressStr, err := maxAddress.MarshalText()
		if err != nil {
			return nil, err
		}
		rr.MaxAddress = string(maxAddressStr)
	}
	return proto.Marshal(rr)
}

func unmarshalRows(data []byte) ([]*Row, *big.Int, *big.Int, error) {
	rows := &Rows{}
	if err := proto.Unmarshal(data, rows); err != nil {
		return nil, nil, nil, err
	}
	minAddress := new(big.Int)
	maxAddress := new(big.Int)
	if rows.MinAddress != "" {
		if err := minAddress.UnmarshalText([]byte(rows.MinAddress)); err != nil {
			return nil, nil, nil, err
		}
	}
	if rows.MaxAddress != "" {
		if err := maxAddress.UnmarshalText([]byte(rows.MaxAddress)); err != nil {
			return nil, nil, nil, err
		}
	}
	return rows.Rows, minAddress, maxAddress, nil
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
