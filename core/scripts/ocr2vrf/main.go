package main

import (
	"flag"
	"os"
	"time"

	"github.com/pkg/errors"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"go.dedis.ch/kyber/v3"
)

// TODO: Remove below once translator implementation is available in ocr2vrf repo
type trivialTranslation struct{ base kyber.Point }

func NewTrivialTranslation(base kyber.Point) *trivialTranslation {
	return &trivialTranslation{base}
}

func (t *trivialTranslation) TranslateKey(share kyber.Scalar) (kyber.Point, error) {
	return t.base.Clone().Mul(share, nil), nil
}

func (t *trivialTranslation) VerifyTranslation(pk1, pk2 kyber.Point) error {
	if pk1.Equal(pk2) {
		return nil
	}
	return errors.Errorf("putative translated points are not equal")
}

func (t *trivialTranslation) Name() string {
	return "trivial translator"
}

func (t *trivialTranslation) TargetGroup(
	sourceGroup kyber.Group,
) (targetGroup kyber.Group, err error) {
	return sourceGroup, nil
}

type commonSetConfigCommands struct {
	onchainPubKeys         string
	offchainPubKeys        string
	configPubKeys          string
	peerIDs                string
	transmitters           string
	schedule               string
	f                      uint
	deltaProgress          time.Duration
	deltaResend            time.Duration
	deltaRound             time.Duration
	deltaGrace             time.Duration
	deltaStage             time.Duration
	maxRounds              time.Duration
	maxDurationQuery       time.Duration
	maxDurationObservation time.Duration
	maxDurationReport      time.Duration
	maxDurationAccept      time.Duration
	maxDurationTransmit    time.Duration
}

type dkgSetConfigCommands struct {
	commonSetConfigCommands
	dkgEncryptionPubKeys string
	dkgSigningPubKeys    string
	keyID                string
}

func main() {
	e := helpers.SetupEnv()

	switch os.Args[1] {

	case "dkg-deploy":
		deployDKG(e)

	case "vrf-deploy":
		cmd := flag.NewFlagSet("vrf-deploy", flag.ExitOnError)
		dkgAddress := cmd.String("dkg-address", "", "dkg contract address")
		keyID := cmd.String("key-id", "", "key ID")
		helpers.ParseArgs(cmd, os.Args[2:], "dkg-address", "key-id")
		deployVRF(e, *dkgAddress, *keyID)

	case "dkg-add-client":
		cmd := flag.NewFlagSet("dkg-add-client", flag.ExitOnError)
		dkgAddress := cmd.String("dkg-address", "", "DKG contract address")
		keyID := cmd.String("key-id", "", "key ID")
		clientAddress := cmd.String("client-address", "", "client address")
		helpers.ParseArgs(cmd, os.Args[2:], "dkg-address", "key-id", "client-address")
		addClientToDKG(e, *dkgAddress, *keyID, *clientAddress)

	case "dkg-remove-client":
		cmd := flag.NewFlagSet("dkg-add-client", flag.ExitOnError)
		dkgAddress := cmd.String("dkg-address", "", "DKG contract address")
		keyID := cmd.String("key-id", "", "key ID")
		clientAddress := cmd.String("client-address", "", "client address")
		helpers.ParseArgs(cmd, os.Args[2:], "dkg-address", "key-id", "client-address")
		removeClientFromDKG(e, *dkgAddress, *keyID, *clientAddress)

	case "dkg-set-config":
		cmd := flag.NewFlagSet("dkg-set-config", flag.ExitOnError)
		dkgAddress := cmd.String("dkg-address", "", "DKG contract address")
		keyID := cmd.String("key-id", "", "key ID")
		onchainPubKeys := cmd.String("onchain-pub-keys", "", "comma-separated list of OCR on-chain pubkeys")
		offchainPubKeys := cmd.String("offchain-pub-keys", "", "comma-separated list of OCR off-chain pubkeys")
		configPubKeys := cmd.String("config-pub-keys", "", "comma-separated list of OCR config pubkeys")
		peerIDs := cmd.String("peer-ids", "", "comma-separated list of peer IDs")
		transmitters := cmd.String("transmitters", "", "comma-separated list transmitters")
		dkgEncryptionPubKeys := cmd.String("dkg-encryption-pub-keys", "", "comma-separated list of DKG encryption pubkeys")
		dkgSigningPubKeys := cmd.String("dkg-signing-pub-keys", "", "comma-separated list of DKG signing pubkeys")
		schedule := cmd.String("schedule", "", "comma-separted list of transmission schedule")
		f := cmd.Uint("f", 1, "number of faulty oracles")
		deltaProgress := cmd.Duration("delta-progress", 30*time.Second, "duration of delta progress")
		deltaResend := cmd.Duration("delta-resend", 10*time.Second, "duration of delta resend")
		deltaRound := cmd.Duration("delta-round", 10*time.Second, "duration of delta round")
		deltaGrace := cmd.Duration("delta-grace", 20*time.Second, "duration of delta grace")
		deltaStage := cmd.Duration("delta-stage", 20*time.Second, "duration of delta grace")
		maxRounds := cmd.Duration("max-rounds", 3*time.Second, "maximum number of rounds")
		maxDurationQuery := cmd.Duration("max-duration-query", 5*time.Second, "maximum duration of query")
		maxDurationObservation := cmd.Duration("max-duration-observation", 5*time.Second, "maximum duration of observation method")
		maxDurationReport := cmd.Duration("max-duration-report", 5*time.Second, "maximum duration of report method")
		maxDurationAccept := cmd.Duration("max-duration-accept", 5*time.Second, "maximum duration of shouldAcceptFinalizedReport method")
		maxDurationTransmit := cmd.Duration("max-duration-transmit", 5*time.Second, "maximum duration of shouldTransmitAcceptedReport method")

		helpers.ParseArgs(cmd,
			os.Args[2:],
			"dkg-address",
			"key-id",
			"onchain-pub-keys",
			"offchain-pub-keys",
			"config-pub-keys",
			"peer-ids",
			"transmitters",
			"dkg-encryption-pub-keys",
			"dkg-signing-pub-keys",
			"schedule")

		commands := dkgSetConfigCommands{
			commonSetConfigCommands: commonSetConfigCommands{
				onchainPubKeys:         *onchainPubKeys,
				offchainPubKeys:        *offchainPubKeys,
				configPubKeys:          *configPubKeys,
				peerIDs:                *peerIDs,
				transmitters:           *transmitters,
				schedule:               *schedule,
				f:                      *f,
				deltaProgress:          *deltaProgress,
				deltaResend:            *deltaResend,
				deltaRound:             *deltaRound,
				deltaGrace:             *deltaGrace,
				deltaStage:             *deltaStage,
				maxRounds:              *maxRounds,
				maxDurationQuery:       *maxDurationQuery,
				maxDurationObservation: *maxDurationObservation,
				maxDurationReport:      *maxDurationReport,
				maxDurationAccept:      *maxDurationAccept,
				maxDurationTransmit:    *maxDurationTransmit,
			},
			dkgEncryptionPubKeys: *dkgEncryptionPubKeys,
			dkgSigningPubKeys:    *dkgSigningPubKeys,
			keyID:                *keyID,
		}

		setDKGConfig(e, *dkgAddress, commands)

	default:
		panic("unrecognized subcommand: " + os.Args[1])
	}
}
