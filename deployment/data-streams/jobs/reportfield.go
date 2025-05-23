package jobs

type ReportFields interface {
	GetStreamType() StreamType
}

type ReportFieldLLO struct {
	ResultPath string
	// StreamID allows assigning an own stream ID to the report field.
	StreamID *string
}

type QuoteReportFields struct {
	Bid       ReportFieldLLO
	Benchmark ReportFieldLLO
	Ask       ReportFieldLLO
}

func (quote QuoteReportFields) GetStreamType() StreamType {
	return StreamTypeQuote
}

type MedianReportFields struct {
	Benchmark ReportFieldLLO
}

func (median MedianReportFields) GetStreamType() StreamType {
	return StreamTypeMedian
}
