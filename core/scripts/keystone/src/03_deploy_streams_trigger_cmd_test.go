package src

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

var (
	chainID         = int64(123456)
	feedID          = fmt.Sprintf("%x", [32]byte{0: 1})
	feedName        = "BTC/USD"
	verifierAddress = fmt.Sprintf("0x%x", [20]byte{0: 7})
)

func TestCreateMercuryV3Job(t *testing.T) {
	ocrKeyBundleID := "ocr_key_bundle_id"
	nodeCSAKey := "node_csa_key"
	bridgeName := "bridge_name"
	linkFeedID := fmt.Sprintf("%x", [32]byte{0: 2})
	nativeFeedID := fmt.Sprintf("%x", [32]byte{0: 3})
	u, err := url.Parse("https://crib-henry-keystone-node1.main.stage.cldev.sh")
	if err != nil {
		t.Fatal(err)
	}

	jobConfigData := MercuryV3JobSpecData{
		BootstrapHost:   u.Hostname(),
		VerifierAddress: verifierAddress,
		OCRKeyBundleID:  ocrKeyBundleID,
		NodeCSAKey:      nodeCSAKey,
		Bridge:          bridgeName,
		FeedName:        feedName,
		FeedID:          feedID,
		LinkFeedID:      linkFeedID,
		NativeFeedID:    nativeFeedID,
		ChainID:         chainID,
	}
	_, output := createMercuryV3Job(jobConfigData)

	snaps.MatchSnapshot(t, output)
}

func TestCreateMercuryBootstrapJob(t *testing.T) {
	jobConfigData := MercuryV3BootstrapJobSpecData{
		FeedName:            feedName,
		FeedID:          feedID,
		ChainID:         chainID,
		VerifierAddress: verifierAddress,
	}

	_, output := createMercuryV3BootstrapJob(jobConfigData)

	snaps.MatchSnapshot(t, output)
}
func TestCreateKeystoneWorkflowJob(t *testing.T) {
	workflowConfig := WorkflowJobSpecConfig{
		JobSpecName:          "keystone_workflow",
		WorkflowOwnerAddress: "0x1234567890abcdef1234567890abcdef12345678",
		FeedIDs:              []string{"feed1", "feed2", "feed3"},
		TargetID:             "target_id",
		TargetAddress:        "0xabcdefabcdefabcdefabcdefabcdefabcdef",
	}

	output := createKeystoneWorkflowJob(workflowConfig)

	snaps.MatchSnapshot(t, output)
}

