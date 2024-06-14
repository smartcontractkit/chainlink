package ocr2vrf_actions

import ocr2vrftypes "github.com/smartcontractkit/ocr2vrf/types"

type DKGKeyConfig struct {
	DKGEncryptionPublicKey string
	DKGSigningPublicKey    string
}

type DKGConfig struct {
	DKGKeyConfigs      []DKGKeyConfig
	DKGKeyID           string
	DKGContractAddress string
}

type VRFBeaconConfig struct {
	VRFBeaconAddress  string
	ConfDelays        []string
	CoordinatorConfig *ocr2vrftypes.CoordinatorConfig
}

type OCR2Config struct {
	OnchainPublicKeys    []string
	OffchainPublicKeys   []string
	PeerIds              []string
	ConfigPublicKeys     []string
	TransmitterAddresses []string
	Schedule             []int
}

type OCR2VRFPluginConfig struct {
	OCR2Config            OCR2Config
	DKGConfig             DKGConfig
	VRFBeaconConfig       VRFBeaconConfig
	VRFCoordinatorAddress string
	LinkEthFeedAddress    string
}
