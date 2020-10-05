package protocol

import (
	"context"

	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/internal/config"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"
	"github.com/smartcontractkit/chainlink/libocr/subprocesses"
)

func RunOracle(
	ctx context.Context,

	config config.SharedConfig,
	contractTransmitter types.ContractTransmitter,
	database types.Database,
	datasource types.DataSource,
	id types.OracleID,
	keys types.PrivateKeys,
	localConfig types.LocalConfig,
	logger types.Logger,
	netEndpoint NetworkEndpoint,
) {
	o := oracleState{
		ctx: ctx,

		Config:              config,
		contractTransmitter: contractTransmitter,
		database:            database,
		datasource:          datasource,
		id:                  id,
		localConfig:         localConfig,
		logger:              logger,
		netEndpoint:         netEndpoint,
		PrivateKeys:         keys,
	}
	o.run()
}

type oracleState struct {
	ctx context.Context

	Config              config.SharedConfig
	configTracker       types.ContractConfigTracker
	contractTransmitter types.ContractTransmitter
	database            types.Database
	datasource          types.DataSource
	id                  types.OracleID
	localConfig         types.LocalConfig
	logger              types.Logger
	netEndpoint         NetworkEndpoint
	PrivateKeys         types.PrivateKeys

	chNetToPacemaker        chan<- MessageToPacemakerWithSender
	chNetToReportGeneration chan<- MessageToReportGenerationWithSender
	childCancel             context.CancelFunc
	childCtx                context.Context
	subprocesses            subprocesses.Subprocesses
}

func (o *oracleState) run() {
	o.logger.Info("Running", nil)

	chNetToPacemaker := make(chan MessageToPacemakerWithSender)
	o.chNetToPacemaker = chNetToPacemaker

	chNetToReportGeneration := make(chan MessageToReportGenerationWithSender)
	o.chNetToReportGeneration = chNetToReportGeneration

	chReportGenerationToTransmission := make(chan EventToTransmission)

	o.childCtx, o.childCancel = context.WithCancel(context.Background())
	defer o.childCancel()

	o.subprocesses.Go(func() {
		RunPacemaker(
			o.childCtx,
			&o.subprocesses,

			chNetToPacemaker,
			chNetToReportGeneration,
			chReportGenerationToTransmission,
			o.Config,
			o.contractTransmitter,
			o.database,
			o.datasource,
			o.id,
			o.localConfig,
			o.logger,
			o.netEndpoint,
			o.PrivateKeys,
		)
	})
	o.subprocesses.Go(func() {
		RunTransmission(
			o.childCtx,
			&o.subprocesses,

			o.Config,
			chReportGenerationToTransmission,
			o.database,
			o.id,
			o.localConfig,
			o.logger,
			o.contractTransmitter,
		)
	})

	chNet := o.netEndpoint.Receive()

	chDone := o.ctx.Done()
	for {
		select {
		case msg := <-chNet:
			msg.Msg.process(o, msg.Sender)
		case <-chDone:
		}

		select {
		case <-chDone:
			o.logger.Debug("Oracle: winding down", nil)
			o.childCancel()
			o.subprocesses.Wait()
			o.logger.Debug("Oracle: exiting", nil)
			return
		default:
		}
	}
}
