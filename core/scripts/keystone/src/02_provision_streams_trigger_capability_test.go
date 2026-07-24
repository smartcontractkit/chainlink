package src

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

var (
	chainID             = int64(123456)
	configuratorAddress = fmt.Sprintf("0x%x", [20]byte{0: 7})
	configStoreAddress  = fmt.Sprintf("0x%x", [20]byte{0: 8})
)

func TestCreateStreamJob(t *testing.T) {
	t.Parallel()
	jobConfigData := StreamJobSpecData{
		FeedName: "BTC/USD",
		StreamID: 1,
		Bridge:   "bridge_name",
	}
	_, output := createStreamJob(jobConfigData)

	snaps.MatchSnapshot(t, output)
}

func TestCreateLLOJob(t *testing.T) {
	t.Parallel()
	u, err := url.Parse("https://crib-henry-keystone-node1.main.stage.cldev.sh")
	if err != nil {
		t.Fatal(err)
	}

	jobConfigData := LLOJobSpecData{
		DonID:                             streamsTriggerDonID,
		BootstrapHost:                     u.Hostname(),
		ConfiguratorAddress:               configuratorAddress,
		ChannelDefinitionsContractAddress: configStoreAddress,
		ChannelDefinitionsFromBlock:       channelDefinitionsFromBlock,
		NodeCSAKey:                        "node_csa_key",
		OCRKeyBundleID:                    "ocr_key_bundle_id",
		ChainID:                           chainID,
	}
	_, output := createLLOJob(jobConfigData)

	snaps.MatchSnapshot(t, output)
}

func TestCreateLLOBootstrapJob(t *testing.T) {
	t.Parallel()
	jobConfigData := LLOBootstrapJobSpecData{
		DonID:                             streamsTriggerDonID,
		ConfiguratorAddress:               configuratorAddress,
		ChannelDefinitionsContractAddress: configStoreAddress,
		ChannelDefinitionsFromBlock:       channelDefinitionsFromBlock,
		ChainID:                           chainID,
	}

	_, output := createLLOBootstrapJob(jobConfigData)

	snaps.MatchSnapshot(t, output)
}
