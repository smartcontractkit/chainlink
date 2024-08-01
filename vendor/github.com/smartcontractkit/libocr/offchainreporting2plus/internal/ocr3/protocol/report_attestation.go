package protocol

import (
	"context"
	"crypto/rand"
	"math"
	"math/big"
	"runtime"
	"sort"
	"sync"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3/scheduler"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

func RunReportAttestation[RI any](
	ctx context.Context,

	chNetToReportAttestation <-chan MessageToReportAttestationWithSender[RI],
	chOutcomeGenerationToReportAttestation <-chan EventToReportAttestation[RI],
	chReportAttestationToTransmission chan<- EventToTransmission[RI],
	config ocr3config.SharedConfig,
	contractTransmitter ocr3types.ContractTransmitter[RI],
	logger loghelper.LoggerWithContext,
	netSender NetworkSender[RI],
	onchainKeyring ocr3types.OnchainKeyring[RI],
	reportingPlugin ocr3types.ReportingPlugin[RI],
) {
	sched := scheduler.NewScheduler[EventMissingOutcome[RI]]()
	defer sched.Close()

	newReportAttestationState(ctx, chNetToReportAttestation,
		chOutcomeGenerationToReportAttestation, chReportAttestationToTransmission,
		config, contractTransmitter, logger, netSender, onchainKeyring, reportingPlugin, sched).run()
}

const expiryMinRounds int = 10
const expiryDuration = 1 * time.Minute
const expiryMaxRounds int = 50

const lookaheadMinRounds int = 4
const lookaheadDuration = 30 * time.Second
const lookaheadMaxRounds int = 10

type reportAttestationState[RI any] struct {
	ctx context.Context

	chNetToReportAttestation               <-chan MessageToReportAttestationWithSender[RI]
	chOutcomeGenerationToReportAttestation <-chan EventToReportAttestation[RI]
	chReportAttestationToTransmission      chan<- EventToTransmission[RI]
	config                                 ocr3config.SharedConfig
	contractTransmitter                    ocr3types.ContractTransmitter[RI]
	logger                                 loghelper.LoggerWithContext
	netSender                              NetworkSender[RI]
	onchainKeyring                         ocr3types.OnchainKeyring[RI]
	reportingPlugin                        ocr3types.ReportingPlugin[RI]

	scheduler *scheduler.Scheduler[EventMissingOutcome[RI]]
	// reap() is used to prevent unbounded state growth of rounds
	rounds map[uint64]*round[RI]
	// highest sequence number for which we have attested reports
	highestAttestedSeqNr uint64
	// highest sequence number for which we have received report signatures
	// from each oracle
	highestReportSignaturesSeqNr []uint64
}

type round[RI any] struct {
	verifiedCertifiedCommit *CertifiedCommit               // only stores certifiedCommit whose qc has been verified
	reportsWithInfo         []ocr3types.ReportWithInfo[RI] // cache result of ReportingPlugin.Reports(certifiedCommit.SeqNr, certifiedCommit.Outcome)
	oracles                 []oracle                       // always initialized to be of length n
	startedFetch            bool
	complete                bool
}

// oracle contains information about interactions with oracles (self & others)
type oracle struct {
	signatures      [][]byte
	validSignatures *bool
	weRequested     bool
	theyServiced    bool
	weServiced      bool
}

func (repatt *reportAttestationState[RI]) run() {
	repatt.logger.Info("ReportAttestation: running", nil)

	for {
		select {
		case msg := <-repatt.chNetToReportAttestation:
			msg.msg.processReportAttestation(repatt, msg.sender)
		case ev := <-repatt.chOutcomeGenerationToReportAttestation:
			ev.processReportAttestation(repatt)
		case ev := <-repatt.scheduler.Scheduled():
			ev.processReportAttestation(repatt)
		case <-repatt.ctx.Done():
		}

		// ensure prompt exit
		select {
		case <-repatt.ctx.Done():
			repatt.logger.Info("ReportAttestation: exiting", nil)
			repatt.scheduler.Close()
			return
		default:
		}
	}
}

func (repatt *reportAttestationState[RI]) messageReportSignatures(
	msg MessageReportSignatures[RI],
	sender commontypes.OracleID,
) {
	if repatt.isBeyondExpiry(msg.SeqNr) {
		repatt.logger.Debug("ignoring MessageReportSignatures for expired seqNr", commontypes.LogFields{
			"seqNr":  msg.SeqNr,
			"sender": sender,
		})
		return
	}

	if repatt.highestReportSignaturesSeqNr[sender] < msg.SeqNr {
		repatt.highestReportSignaturesSeqNr[sender] = msg.SeqNr
	}

	if repatt.isBeyondLookahead(msg.SeqNr) {
		repatt.logger.Debug("ignoring MessageReportSignatures for seqNr beyond lookahead", commontypes.LogFields{
			"seqNr":  msg.SeqNr,
			"sender": sender,
		})
		return
	}

	if _, ok := repatt.rounds[msg.SeqNr]; !ok {
		repatt.rounds[msg.SeqNr] = &round[RI]{
			nil,
			nil,
			make([]oracle, repatt.config.N()),
			false,
			false,
		}
	}

	if len(repatt.rounds[msg.SeqNr].oracles[sender].signatures) != 0 {
		repatt.logger.Debug("ignoring MessageReportSignatures with duplicate signature", commontypes.LogFields{
			"seqNr":  msg.SeqNr,
			"sender": sender,
		})
		return
	}

	repatt.rounds[msg.SeqNr].oracles[sender].signatures = msg.ReportSignatures

	repatt.tryComplete(msg.SeqNr)
}

func (repatt *reportAttestationState[RI]) eventMissingOutcome(ev EventMissingOutcome[RI]) {
	if repatt.rounds[ev.SeqNr].verifiedCertifiedCommit != nil {
		repatt.logger.Debug("dropping EventMissingOutcome, already have Outcome", commontypes.LogFields{
			"seqNr": ev.SeqNr,
		})
		return
	}

	repatt.tryRequestCertifiedCommit(ev.SeqNr)
}

func (repatt *reportAttestationState[RI]) messageCertifiedCommitRequest(msg MessageCertifiedCommitRequest[RI], sender commontypes.OracleID) {
	if repatt.rounds[msg.SeqNr] == nil || repatt.rounds[msg.SeqNr].verifiedCertifiedCommit == nil {
		repatt.logger.Debug("dropping MessageCertifiedCommitRequest for outcome with unknown certified commit", commontypes.LogFields{
			"seqNr":  msg.SeqNr,
			"sender": sender,
		})
		return
	}

	if repatt.rounds[msg.SeqNr].oracles[sender].weServiced {
		repatt.logger.Warn("dropping duplicate MessageCertifiedCommitRequest", commontypes.LogFields{
			"seqNr":  msg.SeqNr,
			"sender": sender,
		})
		return
	}

	repatt.rounds[msg.SeqNr].oracles[sender].weServiced = true

	repatt.logger.Debug("sending MessageCertifiedCommit", commontypes.LogFields{
		"seqNr": msg.SeqNr,
		"to":    sender,
	})
	repatt.netSender.SendTo(MessageCertifiedCommit[RI]{*repatt.rounds[msg.SeqNr].verifiedCertifiedCommit}, sender)
}

func (repatt *reportAttestationState[RI]) messageCertifiedCommit(msg MessageCertifiedCommit[RI], sender commontypes.OracleID) {
	if repatt.rounds[msg.CertifiedCommit.SeqNr] == nil {
		repatt.logger.Warn("dropping MessageCertifiedCommit for unknown seqNr", commontypes.LogFields{
			"seqNr":  msg.CertifiedCommit.SeqNr,
			"sender": sender,
		})
		return
	}

	oracle := &repatt.rounds[msg.CertifiedCommit.SeqNr].oracles[sender]
	if !(oracle.weRequested && !oracle.theyServiced) {
		repatt.logger.Warn("dropping unexpected MessageCertifiedCommit", commontypes.LogFields{
			"seqNr":        msg.CertifiedCommit.SeqNr,
			"sender":       sender,
			"weRequested":  oracle.weRequested,
			"theyServiced": oracle.theyServiced,
		})
		return
	}

	oracle.theyServiced = true

	if repatt.rounds[msg.CertifiedCommit.SeqNr].verifiedCertifiedCommit != nil {
		repatt.logger.Debug("dropping redundant MessageCertifiedCommit", commontypes.LogFields{
			"seqNr":  msg.CertifiedCommit.SeqNr,
			"sender": sender,
		})
		return
	}

	if err := msg.CertifiedCommit.Verify(repatt.config.ConfigDigest, repatt.config.OracleIdentities, repatt.config.ByzQuorumSize()); err != nil {
		repatt.logger.Warn("dropping MessageCertifiedCommit with invalid certified commit", commontypes.LogFields{
			"seqNr":  msg.CertifiedCommit.SeqNr,
			"sender": sender,
		})
		return
	}

	repatt.logger.Debug("received valid MessageCertifiedCommit", commontypes.LogFields{
		"seqNr":  msg.CertifiedCommit.SeqNr,
		"sender": sender,
	})

	repatt.receivedVerifiedCertifiedCommit(msg.CertifiedCommit)
}

func (repatt *reportAttestationState[RI]) tryRequestCertifiedCommit(seqNr uint64) {
	candidates := make([]commontypes.OracleID, 0, repatt.config.N())
	for oracleID, oracle := range repatt.rounds[seqNr].oracles {
		// avoid duplicate requests
		if oracle.weRequested {
			continue
		}
		// avoid requesting from oracles that haven't sent MessageReportSignatures
		if len(oracle.signatures) == 0 {
			continue
		}
		candidates = append(candidates, commontypes.OracleID(oracleID))
	}

	if len(candidates) == 0 {

		return
	}

	randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(candidates))))
	if err != nil {
		repatt.logger.Critical("unexpected error returned by rand.Int", commontypes.LogFields{
			"error": err,
		})
		return
	}
	randomCandidate := candidates[int(randomIndex.Int64())]
	repatt.rounds[seqNr].oracles[randomCandidate].weRequested = true
	repatt.logger.Debug("sending MessageCertifiedCommitRequest", commontypes.LogFields{
		"seqNr": seqNr,
		"to":    randomCandidate,
	})
	repatt.netSender.SendTo(MessageCertifiedCommitRequest[RI]{seqNr}, randomCandidate)
	repatt.scheduler.ScheduleDelay(EventMissingOutcome[RI]{seqNr}, repatt.config.DeltaCertifiedCommitRequest)
}

func (repatt *reportAttestationState[RI]) tryComplete(seqNr uint64) {
	if repatt.rounds[seqNr].complete {
		repatt.logger.Debug("cannot complete, already completed", commontypes.LogFields{
			"seqNr": seqNr,
		})
		return
	}

	if repatt.rounds[seqNr].verifiedCertifiedCommit == nil {
		oraclesThatSentSignatures := 0
		for _, oracle := range repatt.rounds[seqNr].oracles {
			if len(oracle.signatures) == 0 {
				continue
			}
			oraclesThatSentSignatures++
		}

		if oraclesThatSentSignatures <= repatt.config.F {
			repatt.logger.Debug("cannot complete, missing CertifiedCommit and signatures", commontypes.LogFields{
				"oraclesThatSentSignatures": oraclesThatSentSignatures,
				"seqNr":                     seqNr,
				"threshold":                 repatt.config.F + 1,
			})
		} else if !repatt.rounds[seqNr].startedFetch {
			repatt.rounds[seqNr].startedFetch = true
			repatt.scheduler.ScheduleDelay(EventMissingOutcome[RI]{seqNr}, repatt.config.DeltaCertifiedCommitRequest)
		}
		return
	}

	reportsWithInfo := repatt.rounds[seqNr].reportsWithInfo
	goodSigs := 0
	var aossPerReport [][]types.AttributedOnchainSignature = make([][]types.AttributedOnchainSignature, len(reportsWithInfo))
	for oracleID := range repatt.rounds[seqNr].oracles {
		oracle := &repatt.rounds[seqNr].oracles[oracleID]
		if len(oracle.signatures) == 0 {
			continue
		}
		if oracle.validSignatures == nil {
			validSignatures := repatt.verifySignatures(
				repatt.config.OracleIdentities[oracleID].OnchainPublicKey,
				seqNr,
				reportsWithInfo,
				oracle.signatures,
			)
			oracle.validSignatures = &validSignatures
			if !validSignatures {
				// Other less common causes include actually invalid signatures.
				repatt.logger.Warn("report signatures failed to verify. This is commonly caused by non-determinism in the ReportingPlugin", commontypes.LogFields{
					"sender": oracleID,
					"seqNr":  seqNr,
				})
			}
		}
		if oracle.validSignatures != nil && *oracle.validSignatures {
			goodSigs++

			for i := range reportsWithInfo {
				aossPerReport[i] = append(aossPerReport[i], types.AttributedOnchainSignature{
					oracle.signatures[i],
					commontypes.OracleID(oracleID),
				})
			}
		}
		if goodSigs > repatt.config.F {
			break
		}
	}

	if goodSigs <= repatt.config.F {
		repatt.logger.Debug("cannot complete, insufficient number of signatures", commontypes.LogFields{
			"seqNr":     seqNr,
			"goodSigs":  goodSigs,
			"threshold": repatt.config.F + 1,
		})
		return
	}

	if repatt.highestAttestedSeqNr < seqNr {
		repatt.highestAttestedSeqNr = seqNr
	}

	repatt.rounds[seqNr].complete = true

	repatt.logger.Debug("sending attested reports to transmission protocol", commontypes.LogFields{
		"seqNr":   seqNr,
		"reports": len(reportsWithInfo),
	})

	for i := range reportsWithInfo {
		select {
		case repatt.chReportAttestationToTransmission <- EventAttestedReport[RI]{
			seqNr,
			i,
			AttestedReportMany[RI]{
				reportsWithInfo[i],
				aossPerReport[i],
			},
		}:
		case <-repatt.ctx.Done():
		}
	}

	repatt.reap()
}

func (repatt *reportAttestationState[RI]) verifySignatures(publicKey types.OnchainPublicKey, seqNr uint64, reportsWithInfo []ocr3types.ReportWithInfo[RI], signatures [][]byte) bool {
	if len(reportsWithInfo) != len(signatures) {
		return false
	}

	n := runtime.GOMAXPROCS(0)
	if (len(reportsWithInfo)+3)/4 < n {
		n = (len(reportsWithInfo) + 3) / 4
	}

	var wg sync.WaitGroup
	wg.Add(n)

	var mutex sync.Mutex
	allValid := true

	for k := 0; k < n; k++ {
		k := k

		go func() {
			defer wg.Done()
			for i := k; i < len(reportsWithInfo); i += n {
				if i%n != k {
					panic("bug")
				}

				mutex.Lock()
				allValidCopy := allValid
				mutex.Unlock()

				if !allValidCopy {
					return
				}

				if !repatt.onchainKeyring.Verify(publicKey, repatt.config.ConfigDigest, seqNr, reportsWithInfo[i], signatures[i]) {
					mutex.Lock()
					allValid = false
					mutex.Unlock()
					return
				}
			}
		}()
	}

	wg.Wait()

	return allValid
}

func (repatt *reportAttestationState[RI]) eventCommittedOutcome(ev EventCommittedOutcome[RI]) {
	repatt.receivedVerifiedCertifiedCommit(ev.CertifiedCommit)
}

func (repatt *reportAttestationState[RI]) receivedVerifiedCertifiedCommit(certifiedCommit CertifiedCommit) {
	if repatt.rounds[certifiedCommit.SeqNr] != nil && repatt.rounds[certifiedCommit.SeqNr].verifiedCertifiedCommit != nil {
		repatt.logger.Debug("dropping redundant CertifiedCommit", commontypes.LogFields{
			"seqNr": certifiedCommit.SeqNr,
		})
		return
	}

	reportsWithInfo, ok := callPlugin[[]ocr3types.ReportWithInfo[RI]](
		repatt.ctx,
		repatt.logger,
		commontypes.LogFields{"seqNr": certifiedCommit.SeqNr},
		"Reports",
		0, // Reports is a pure function and should finish "instantly"
		func(context.Context) ([]ocr3types.ReportWithInfo[RI], error) {
			return repatt.reportingPlugin.Reports(
				certifiedCommit.SeqNr,
				certifiedCommit.Outcome,
			)
		},
	)
	if !ok {
		return
	}

	repatt.logger.Debug("successfully invoked ReportingPlugin.Reports", commontypes.LogFields{
		"seqNr":   certifiedCommit.SeqNr,
		"reports": len(reportsWithInfo),
	})

	var sigs [][]byte
	for i, reportWithInfo := range reportsWithInfo {
		sig, err := repatt.onchainKeyring.Sign(repatt.config.ConfigDigest, certifiedCommit.SeqNr, reportWithInfo)
		if err != nil {
			repatt.logger.Error("error while signing report", commontypes.LogFields{
				"seqNr": certifiedCommit.SeqNr,
				"index": i,
				"error": err,
			})
			return
		}
		sigs = append(sigs, sig)
	}

	if _, ok := repatt.rounds[certifiedCommit.SeqNr]; !ok {
		repatt.rounds[certifiedCommit.SeqNr] = &round[RI]{
			nil,
			nil,
			make([]oracle, repatt.config.N()),
			false,
			false,
		}
	}
	repatt.rounds[certifiedCommit.SeqNr].verifiedCertifiedCommit = &certifiedCommit
	repatt.rounds[certifiedCommit.SeqNr].reportsWithInfo = reportsWithInfo

	repatt.logger.Debug("broadcasting MessageReportSignatures", commontypes.LogFields{
		"seqNr": certifiedCommit.SeqNr,
	})

	repatt.netSender.Broadcast(MessageReportSignatures[RI]{
		certifiedCommit.SeqNr,
		sigs,
	})

	// no need to call tryComplete since receipt of our own MessageReportSignatures will do so
}

func (repatt *reportAttestationState[RI]) isBeyondExpiry(seqNr uint64) bool {
	highest := repatt.highestAttestedSeqNr
	expiry := uint64(repatt.expiryRounds())
	if highest <= expiry {
		return false
	}
	return seqNr < highest-expiry
}

func (repatt *reportAttestationState[RI]) isBeyondLookahead(seqNr uint64) bool {
	highestReportSignaturesSeqNr := append([]uint64{}, repatt.highestReportSignaturesSeqNr...)
	sort.Slice(highestReportSignaturesSeqNr, func(i, j int) bool {
		return highestReportSignaturesSeqNr[i] > highestReportSignaturesSeqNr[j]
	})
	highest := highestReportSignaturesSeqNr[repatt.config.F] // (f+1)th largest seqNr
	lookahead := uint64(repatt.lookaheadRounds())
	if seqNr <= lookahead {
		return false
	}
	return highest < seqNr-lookahead
}

// reap expired entries from repatt.finalized to prevent unbounded state growth
func (repatt *reportAttestationState[RI]) reap() {
	maxActiveRoundCount := repatt.expiryRounds() + repatt.lookaheadRounds()
	// only reap if more than ~ a third of the rounds can be discarded
	if 3*len(repatt.rounds) <= 4*maxActiveRoundCount {
		return
	}
	// A long time ago in a galaxy far, far away, Go used to leak memory when
	// repeatedly adding and deleting from the same map without ever exceeding
	// some maximum length. Fortunately, this is no longer the case
	// https://go-review.googlesource.com/c/go/+/25049/
	for seqNr := range repatt.rounds {
		if repatt.isBeyondExpiry(seqNr) {
			delete(repatt.rounds, seqNr)
		}
	}
}

// The age (denoted in rounds) after which a report is considered expired and
// will automatically be dropped
func (repatt *reportAttestationState[RI]) expiryRounds() int {
	return repatt.roundWindowSize(expiryMinRounds, expiryMaxRounds, expiryDuration)
}

// The lookahead (denoted in rounds) after which a report is considered too far in the future and
// will automatically be dropped
func (repatt *reportAttestationState[RI]) lookaheadRounds() int {
	return repatt.roundWindowSize(lookaheadMinRounds, lookaheadMaxRounds, lookaheadDuration)
}

func (repatt *reportAttestationState[RI]) roundWindowSize(minWindowSize int, maxWindowSize int, windowDuration time.Duration) int {
	// number of rounds in a window of duration expirationAgeDuration
	size := math.Ceil(windowDuration.Seconds() / repatt.config.MinRoundInterval().Seconds())

	if size < float64(minWindowSize) {
		size = float64(minWindowSize)
	}
	if math.IsNaN(size) || size > float64(maxWindowSize) {
		size = float64(maxWindowSize)
	}

	return int(math.Ceil(size))
}

func newReportAttestationState[RI any](
	ctx context.Context,

	chNetToReportAttestation <-chan MessageToReportAttestationWithSender[RI],
	chOutcomeGenerationToReportAttestation <-chan EventToReportAttestation[RI],
	chReportAttestationToTransmission chan<- EventToTransmission[RI],
	config ocr3config.SharedConfig,
	contractTransmitter ocr3types.ContractTransmitter[RI],
	logger loghelper.LoggerWithContext,
	netSender NetworkSender[RI],
	onchainKeyring ocr3types.OnchainKeyring[RI],
	reportingPlugin ocr3types.ReportingPlugin[RI],
	sched *scheduler.Scheduler[EventMissingOutcome[RI]],
) *reportAttestationState[RI] {
	return &reportAttestationState[RI]{
		ctx,

		chNetToReportAttestation,
		chOutcomeGenerationToReportAttestation,
		chReportAttestationToTransmission,
		config,
		contractTransmitter,
		logger.MakeUpdated(commontypes.LogFields{"proto": "repatt"}),
		netSender,
		onchainKeyring,
		reportingPlugin,

		sched,
		map[uint64]*round[RI]{},
		0,
		make([]uint64, config.N()),
	}
}
