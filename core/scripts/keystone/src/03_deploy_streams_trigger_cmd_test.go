package src

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

var (
	chainID         = int64(123456)
	feedID          = [32]byte{0: 1}
	feedName        = "BTC/USD"
	verifierAddress = [20]byte{0: 7}
)

func TestCreateMercuryV3Job(t *testing.T) {
	ocrKeyBundleID := "ocr_key_bundle_id"
	nodeCSAKey := "node_csa_key" 
	bridgeName := "bridge_name"
	linkFeedID := [32]byte{0: 2}
	nativeFeedID := [32]byte{0: 3}

	_, output := createMercuryV3Job(
		ocrKeyBundleID,
		verifierAddress,
		bridgeName,
		nodeCSAKey,
		feedName,
		feedID,
		linkFeedID,
		nativeFeedID,
		chainID,
	)

	snaps.MatchSnapshot(t, output)
}

func TestCreateMercuryBootstrapJob(t *testing.T) {
	_, output := createMercuryV3BootstrapJob(
		verifierAddress,
		feedName,
		feedID,
		chainID,
	)

	snaps.MatchSnapshot(t, output)
}
