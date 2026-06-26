package ocrcommon

import (
	commonocr "github.com/smartcontractkit/chainlink-common/pkg/ocrcommon"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
)

// DiscovererDatabase is a key-value store for p2p announcements. The
// implementation now lives in chainlink-common so it can be reused; this alias
// preserves the existing core import path.
type DiscovererDatabase = commonocr.DiscovererDatabase

const (
	// ocrDiscovererTable is the name of the table used to store OCR announcements
	ocrDiscovererTable = "ocr_discoverer_announcements"
	// don2donDiscovererTable is the name of the table used to store DON2DON announcements
	don2donDiscovererTable = "don2don_discoverer_announcements"
)

// NewOCRDiscovererDatabase creates a new DiscovererDatabase for OCR announcements
func NewOCRDiscovererDatabase(ds sqlutil.DataSource, peerID string) *DiscovererDatabase {
	return commonocr.NewDiscovererDatabase(ds, peerID, ocrDiscovererTable)
}

// NewDON2DONDiscovererDatabase creates a new DiscovererDatabase for DON2DON announcements
func NewDON2DONDiscovererDatabase(ds sqlutil.DataSource, peerID string) *DiscovererDatabase {
	return commonocr.NewDiscovererDatabase(ds, peerID, don2donDiscovererTable)
}
