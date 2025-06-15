package bindings

import "encoding/binary"

// This function is not EVM specific, it's generic and should be provided by CRE
func GenerateReport(chainID uint32, userData []byte) commonReport {
	// This isn't right, but is fine for testing for now
	var allData []byte
	allData = binary.LittleEndian.AppendUint32(allData, chainID)
	allData = append(allData, userData...)
	return commonReport{RawReport: allData}

}

// This is not EVM specific, it's generic
type commonReport struct {
	RawReport     []byte
	ReportContext []byte
	Signatures    [][]byte
	ID            []byte
}
