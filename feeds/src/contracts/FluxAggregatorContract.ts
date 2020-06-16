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
import { Log, Filter, TransactionResponse, Block } from 'ethers/providers'

interface EventListener {
  filter: ethers.providers.EventType
  listener: ethers.providers.Listener
}

interface Meta {
  timestamp?: number
  gasPrice?: string
  transactionHash: string
  blockNumber: number
}

interface DecodedLog {
  meta: Meta
}

interface SubmissionReceivedEventLog extends Log {
  answer: number
  round: number
  oracle: string
}

interface NewRoundEventLog extends Log {
  roundId: number
  startedBy: number
  startedAt: number
}

interface AnswerUpdatedLog extends Log {
  current: number
  roundId: number
  timestamp: number
}

export interface AnswerUpdatedLogFormat extends DecodedLog {
  answerFormatted: string
  answer: number
  answerId: number
  timestamp: number
}

export interface SubmissionReceivedEventLogFormat extends DecodedLog {
  answerFormatted: string
  answer: number
  answerId: number
  sender: string
  timestamp: number
}

export interface NewRoundEventLogFormat extends DecodedLog {
  startedBy: number
  startedAt: number
  answerId: number
}

export default class FluxContract {
  private submissionReceivedEvent: EventListener = {
    filter: {},
    listener: () => {},
  }

  private newRoundEvent: EventListener = {
    filter: {},
    listener: () => {},
  }

  private answerIdInterval: ReturnType<typeof setTimeout | any> = null
  private alive: boolean
  private options: FeedConfig

  private contract: ethers.Contract | null
  provider: ethers.providers.JsonRpcProvider
  address: string

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
      this.alive = false
    } catch {
      console.error('Cannot delete FluxContract')
    }
  }

  removeListener(
    filter: ethers.providers.EventType,
    eventListener: ethers.providers.Listener,
  ): void {
    if (!this.alive) return

    this.provider.removeListener(filter, eventListener)
  }

  async oracles(): Promise<string[]> {
    if (!this.contract) {
      throw Error('Contract instance does not exist')
    }

    return await this.contract.getOracles()
  }

  async minimumAnswers(): Promise<number> {
    if (!this.contract) {
      throw Error('Contract instance does not exist')
    }

    return await this.contract.minAnswerCount()
  }

  async latestRound(): Promise<number> {
    if (!this.contract) {
      throw Error('Contract instance does not exist')
    }

    const latestRound = await this.contract.latestRound()
    return latestRound.toNumber()
  }

  async reportingRound(): Promise<number> {
    if (!this.contract) {
      throw Error('Contract instance does not exist')
    }

    const reportingRound = await this.contract.reportingRound()
    return reportingRound.toNumber()
  }

  async latestAnswer(): Promise<string> {
    if (!this.contract) {
      throw Error('Contract instance does not exist')
    }

    const latestAnswer = await this.contract.latestAnswer()
    return formatAnswer(
      latestAnswer,
      this.options.multiply,
      this.options.decimalPlaces,
      this.options.formatDecimalPlaces,
    )
  }

  async latestTimestamp(): Promise<number> {
    if (!this.contract) {
      throw Error('Contract instance does not exist')
    }

    const latestTimestamp = await this.contract.latestTimestamp()
    return latestTimestamp.toNumber()
  }

  async getAnswer(answerId: number): Promise<string> {
    if (!this.contract) {
      throw Error('Contract instance does not exist')
    }

    const getAnswer = await this.contract.getAnswer(answerId)
    return formatAnswer(
      getAnswer,
      this.options.multiply,
      this.options.decimalPlaces,
      this.options.formatDecimalPlaces,
    )
  }

  async getTimestamp(answerId: number): Promise<number> {
    if (!this.contract) {
      throw Error('Contract instance does not exist')
    }

    const timestamp = await this.contract.getTimestamp(answerId)
    return timestamp.toNumber()
  }

  async description(): Promise<string> {
    if (!this.contract) {
      throw Error('Contract instance does not exist')
    }

    const description = await this.contract.description()
    return ethers.utils.parseBytes32String(description)
  }

  async decimals(): Promise<number> {
    if (!this.contract) {
      throw Error('Contract instance does not exist')
    }

    return await this.contract.decimals()
  }

  async listenSubmissionReceivedEvent(callback: Function | undefined) {
    if (!this.contract) {
      throw Error('Contract instance does not exist')
    }

    this.removeListener(
      this.submissionReceivedEvent.filter,
      this.submissionReceivedEvent.listener,
    )

    this.submissionReceivedEvent.filter = {
      ...this.contract.filters.SubmissionReceived(null, null, null),
    }

    return this.provider.on(
      this.submissionReceivedEvent.filter,
      (this.submissionReceivedEvent.listener = async (log: Log) => {
        if (!this.contract) {
          throw Error('Contract instance does not exist')
        }

        const logged = decodeLog(
          {
            log,
            eventInterface: this.contract.interface.events.SubmissionReceived,
          },
          (decodedLog: SubmissionReceivedEventLog) => ({
            answerFormatted: formatAnswer(
              decodedLog.answer,
              this.options.multiply,
              this.options.decimalPlaces,
              this.options.formatDecimalPlaces,
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

  async listenNewRoundEvent(callback: Function | undefined) {
    if (!this.contract) {
      throw Error('Contract instance does not exist')
    }

    this.removeListener(this.newRoundEvent.filter, this.newRoundEvent.listener)

    this.newRoundEvent.filter = {
      ...this.contract.filters.NewRound(null, null, null),
    }

    return this.provider.on(
      this.newRoundEvent.filter,
      (this.newRoundEvent.listener = async (log: Log) => {
        if (!this.contract) {
          throw Error('Contract instance does not exist')
        }

        const logged = decodeLog(
          {
            log,
            eventInterface: this.contract.interface.events.NewRound,
          },
          (decodedLog: NewRoundEventLog) => ({
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
    if (!this.contract) {
      throw Error('Contract instance does not exist')
    }

    const newRoundFilter: Filter = {
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
      (decodedLog: NewRoundEventLog) => ({
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
    if (!this.contract) {
      throw Error('Contract instance does not exist')
    }

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
      (decodedLog: SubmissionReceivedEventLog) => ({
        answerFormatted: formatAnswer(
          decodedLog.answer,
          this.options.multiply,
          this.options.decimalPlaces,
          this.options.formatDecimalPlaces,
        ),
        answer: Number(decodedLog.answer),
        answerId: Number(decodedLog.round),
        sender: decodedLog.oracle,
      }),
    )

    return logs
  }

  async answerUpdatedLogs({
    fromBlock,
  }: {
    fromBlock: number
  }): Promise<AnswerUpdatedLogFormat[]> {
    if (!this.contract) {
      throw Error('Contract instance does not exist')
    }

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
      (decodedLog: AnswerUpdatedLog) => ({
        answerFormatted: formatAnswer(
          decodedLog.current,
          this.options.multiply,
          this.options.decimalPlaces,
          this.options.formatDecimalPlaces,
        ),
        answer: Number(decodedLog.current),
        answerId: Number(decodedLog.roundId),
        timestamp: Number(decodedLog.timestamp),
      }),
    )

    return logs
  }

  async addBlockTimestampToLogs(
    logs: DecodedLog[],
  ): Promise<Array<DecodedLog>> {
    if (_.isEmpty(logs)) return logs

    const blockTimePromises = logs.map((log: any) =>
      this.provider.getBlock(log.meta.blockNumber),
    )
    const blockTimes: Block[] = await Promise.all(blockTimePromises)

    return logs.map((l: DecodedLog, i: number) => {
      l.meta.timestamp = blockTimes[i].timestamp
      return l
    })
  }

  async addGasPriceToLogs(logs: DecodedLog[]): Promise<Array<DecodedLog>> {
    if (!logs) return logs

    const logsWithGasPriceOps = logs.map(async (log: DecodedLog) => {
      const tx: TransactionResponse = await this.provider.getTransaction(
        log.meta.transactionHash,
      )
      // eslint-disable-next-line require-atomic-updates
      log.meta.gasPrice = ethers.utils.formatUnits(tx.gasPrice, 'gwei')
      return log
    })

    return Promise.all(logsWithGasPriceOps)
  }
}
