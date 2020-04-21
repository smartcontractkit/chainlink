import { ethers } from 'ethers'
import {
  getLogs,
  formatAnswer,
  decodeLog,
  createContract,
  createInfuraProvider,
} from './utils'
import _ from 'lodash'
import { FeedConfig } from 'config'

export default class FluxContract {
  private submissionReceivedEvent: any = {
    filter: {},
    listener: {},
  }

  private newRoundEvent: any = {
    filter: {},
    listener: {},
  }

  private answerIdInterval: ReturnType<typeof setTimeout | any> = null
  private alive: boolean
  private options: FeedConfig | any

  private provider: ethers.providers.JsonRpcProvider
  private contract: ethers.Contract | any
  address: string | any

  constructor(options: FeedConfig, abi: any) {
    this.provider = createInfuraProvider(options.networkId)
    this.contract = createContract(options.contractAddress, this.provider, abi)
    this.alive = true
    this.options = options
    this.address = options.contractAddress
  }

  kill(): void {
    try {
      if (!this.alive) return
      clearInterval(this.answerIdInterval)
      this.removeListener(
        this.submissionReceivedEvent.filter,
        this.submissionReceivedEvent.listener,
      )
      this.removeListener(
        this.newRoundEvent.filter,
        this.newRoundEvent.listener,
      )
      this.contract = null
      this.address = null
      this.alive = false
      this.options = null
    } catch {
      console.error('Cannot delete FluxContract')
    }
  }

  removeListener(filter: any, eventListener: any): void {
    if (!this.alive) return

    this.provider.removeListener(filter, eventListener)
  }

  async oracles(): Promise<string[]> {
    return await this.contract.getOracles()
  }

  async minimumAnswers(): Promise<number> {
    return await this.contract.minAnswerCount()
  }

  async latestRound(): Promise<number> {
    const latestRound = await this.contract.latestRound()
    this.decimals()
    return latestRound.toNumber()
  }

  async reportingRound(): Promise<number> {
    const reportingRound = await this.contract.reportingRound()
    return reportingRound.toNumber()
  }

  async latestAnswer(): Promise<string> {
    const latestAnswer = await this.contract.latestAnswer()
    return formatAnswer(
      latestAnswer,
      this.options.multiply,
      this.options.decimalPlaces,
    )
  }

  async latestTimestamp(): Promise<number> {
    const latestTimestamp = await this.contract.latestTimestamp()
    return latestTimestamp.toNumber()
  }

  async getAnswer(answerId: number): Promise<string> {
    const getAnswer = await this.contract.getAnswer(answerId)
    return formatAnswer(
      getAnswer,
      this.options.multiply,
      this.options.decimalPlaces,
    )
  }

  async getTimestamp(answerId: number): Promise<number> {
    const timestamp = await this.contract.getTimestamp(answerId)
    return timestamp.toNumber()
  }

  async description(): Promise<string> {
    const description = await this.contract.description()
    return ethers.utils.parseBytes32String(description)
  }

  async decimals(): Promise<number> {
    return await this.contract.decimals()
  }

  async listenSubmissionReceivedEvent(callback: any) {
    if (!this.alive) return

    this.removeListener(
      this.submissionReceivedEvent.filter,
      this.submissionReceivedEvent.listener,
    )

    this.submissionReceivedEvent.filter = {
      ...this.contract.filters.SubmissionReceived(null, null, null),
    }

    return this.provider.on(
      this.submissionReceivedEvent.filter,
      (this.submissionReceivedEvent.listener = async (log: any) => {
        const logged = decodeLog(
          {
            log,
            eventInterface: this.contract.interface.events.SubmissionReceived,
          },
          (decodedLog: any) => ({
            answerFormatted: formatAnswer(
              decodedLog.answer,
              this.options.multiply,
              this.options.decimalPlaces,
            ),
            answer: Number(decodedLog.answer),
            answerId: Number(decodedLog.round),
            sender: decodedLog.oracle,
          }),
        )
        const logWithTimestamp = await this.addBlockTimestampToLogs([logged])
        const logWithGasPrice = await this.addGasPriceToLogs(
          logWithTimestamp,
        ).then(l => l[0])

        return callback ? callback(logWithGasPrice) : logWithGasPrice
      }),
    )
  }

  async listenNewRoundEvent(callback: any) {
    if (!this.alive) return

    this.removeListener(this.newRoundEvent.filter, this.newRoundEvent.listener)

    this.newRoundEvent.filter = {
      ...this.contract.filters.NewRound(null, null, null),
    }

    return this.provider.on(
      this.newRoundEvent.filter,
      (this.newRoundEvent.listener = async (log: any) => {
        const logged = decodeLog(
          {
            log,
            eventInterface: this.contract.interface.events.NewRound,
          },
          (decodedLog: any) => ({
            answerId: Number(decodedLog.roundId),
            startedBy: decodedLog.startedBy,
            startedAt: Number(decodedLog.startedAt),
          }),
        )

        return callback ? callback(logged) : logged
      }),
    )
  }

  async newRoundLogs({
    fromBlock,
    round,
  }: {
    fromBlock: number
    round: number
  }) {
    const newRoundFilter = {
      ...this.contract.filters.NewRound(round, null, null),
      fromBlock,
      toBlock: 'latest',
    }
    const logs = await getLogs(
      {
        provider: this.provider,
        filter: newRoundFilter,
        eventInterface: this.contract.interface.events.NewRound,
      },
      (decodedLog: any) => ({
        answerId: Number(decodedLog.roundId),
        startedBy: decodedLog.startedBy,
        startedAt: Number(decodedLog.startedAt),
      }),
    )

    return logs
  }

  async submissionReceivedLogs({
    fromBlock,
    round,
  }: {
    fromBlock: number
    round: number
  }) {
    const submissionReceivedFilter = {
      ...this.contract.filters.SubmissionReceived(null, round, null),
      fromBlock,
      toBlock: 'latest',
    }
    const logs = await getLogs(
      {
        provider: this.provider,
        filter: submissionReceivedFilter,
        eventInterface: this.contract.interface.events.SubmissionReceived,
      },
      (decodedLog: any) => ({
        answerFormatted: formatAnswer(
          decodedLog.answer,
          this.options.multiply,
          this.options.decimalPlaces,
        ),
        answer: Number(decodedLog.answer),
        answerId: Number(decodedLog.round),
        sender: decodedLog.oracle,
      }),
    )

    return logs
  }

  async answerUpdatedLogs({ fromBlock }: { fromBlock: number }) {
    const filter = {
      ...this.contract.filters.AnswerUpdated(null, null, null),
      fromBlock,
      toBlock: 'latest',
    }
    const logs = await getLogs(
      {
        provider: this.provider,
        filter,
        eventInterface: this.contract.interface.events.AnswerUpdated,
      },
      (decodedLog: any) => ({
        answerFormatted: formatAnswer(
          decodedLog.current,
          this.options.multiply,
          this.options.decimalPlaces,
        ),
        answer: Number(decodedLog.current),
        answerId: Number(decodedLog.roundId),
        timestamp: Number(decodedLog.timestamp),
      }),
    )

    return logs
  }

  async addBlockTimestampToLogs(logs: any) {
    if (_.isEmpty(logs)) return logs

    const blockTimePromises = logs.map((log: any) =>
      this.provider.getBlock(log.meta.blockNumber),
    )
    const blockTimes: any = await Promise.all(blockTimePromises)

    return logs.map((l: any, i: number) => {
      l.meta.timestamp = blockTimes[i].timestamp
      return l
    })
  }

  async addGasPriceToLogs(logs: any): Promise<Array<any>> {
    if (!logs) return logs

    const logsWithGasPriceOps = logs.map(async (log: any) => {
      const tx = await this.provider.getTransaction(log.meta.transactionHash)
      // eslint-disable-next-line require-atomic-updates
      log.meta.gasPrice = ethers.utils.formatUnits(tx.gasPrice, 'gwei')
      return log
    })

    return Promise.all(logsWithGasPriceOps)
  }
}
