import { ethers } from 'ethers'
import {
  getLogs,
  formatAnswer,
  decodeLog,
  createContract,
  createInfuraProvider,
} from './utils'
import _ from 'lodash'

export default class PrepaidContract {
  answerUpdatedEvent = {
    filter: {},
    listener: {},
  }
  submissionReceivedEvent = {
    filter: {},
    listener: {},
  }

  newRoundEvent = {
    filter: {},
    listener: {},
  }

  answerIdInterval
  provider
  contract

  constructor(options, abi) {
    this.provider = createInfuraProvider(options.networkId)
    this.contract = createContract(options.contractAddress, this.provider, abi)
    this.address = options.contractAddress
    this.alive = true
    this.abi = abi
    this.options = options
  }

  kill() {
    try {
      if (!this.alive) return false
      clearInterval(this.answerIdInterval)
      this.removeListener(
        this.answerUpdatedEvent.filter,
        this.answerUpdatedEvent.listener,
      )
      this.contract = null
      this.address = null
      this.alive = false
      this.options = null
    } catch (error) {
      //
    }
  }

  removeListener(filter, eventListener) {
    if (!this.alive) return

    this.provider.removeListener(filter, eventListener)
  }

  async oracles() {
    return await this.contract.getOracles()
  }

  async minimumAnswers() {
    const minAnswers = await this.contract.minAnswerCount()
    return Number(minAnswers)
  }

  async latestRound() {
    const latestRound = await this.contract.latestRound()
    return Number(latestRound)
  }

  async reportingRound() {
    const reportingRound = await this.contract.reportingRound()
    return Number(reportingRound)
  }

  async latestAnswer() {
    const latestAnswer = await this.contract.latestAnswer()
    return formatAnswer(
      latestAnswer,
      this.options.multiply,
      this.options.decimalPlaces,
    )
  }

  async latestTimestamp() {
    const latestTimestamp = await this.contract.latestTimestamp()
    return Number(latestTimestamp)
  }

  async getAnswer(answerId) {
    const getAnswer = await this.contract.getAnswer(answerId)
    return formatAnswer(
      getAnswer,
      this.options.multiply,
      this.options.decimalPlaces,
    )
  }

  async getTimestamp(answerId) {
    const timestamp = await this.contract.getTimestamp(answerId)
    return Number(timestamp)
  }

  async description() {
    const description = await this.contract.description()
    return ethers.utils.parseBytes32String(description)
  }

  async decimals() {
    const decimals = await this.contract.decimals()
    return Number(ethers.utils.bigNumberify(decimals))
  }

  async latestSubmission(oracles) {
    const submissions = oracles.map(async oracle => {
      const submission = await this.contract.latestSubmission(oracle)
      return {
        responseFormatted: formatAnswer(
          submission[0],
          this.options.multiply,
          this.options.decimalPlaces,
        ),
        response: Number(submission[0]),
        answerId: Number(submission[1]),
        sender: oracle,
      }
    })
    return Promise.all(submissions)
  }

  async listenSubmissionReceivedEvent(callback) {
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
      (this.submissionReceivedEvent.listener = async log => {
        const logged = decodeLog(
          {
            log,
            eventInterface: this.contract.interface.events.SubmissionReceived,
          },
          decodedLog => ({
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

  async listenNewRoundEvent(callback) {
    if (!this.alive) return

    this.removeListener(this.newRoundEvent.filter, this.newRoundEvent.listener)

    this.newRoundEvent.filter = {
      ...this.contract.filters.NewRound(null, null, null),
    }

    return this.provider.on(
      this.newRoundEvent.filter,
      (this.newRoundEvent.listener = async log => {
        const logged = decodeLog(
          {
            log,
            eventInterface: this.contract.interface.events.NewRound,
          },
          decodedLog => ({
            answerId: Number(decodedLog.roundId),
            startedBy: decodedLog.startedBy,
            startedAt: Number(decodedLog.startedAt),
          }),
        )

        return callback ? callback(logged) : logged
      }),
    )
  }

  async newRoundLogs({ fromBlock, round }) {
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
      decodedLog => ({
        answerId: Number(decodedLog.roundId),
        startedBy: decodedLog.startedBy,
        startedAt: Number(decodedLog.startedAt),
      }),
    )

    return logs
  }

  async submissionReceivedLogs({ fromBlock, round }) {
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
      decodedLog => ({
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

  async answerUpdatedLogs({ fromBlock }) {
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
      decodedLog => ({
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

  async addBlockTimestampToLogs(logs) {
    if (_.isEmpty(logs)) return logs

    const blockTimePromises = []

    for (let i = 0; i < logs.length; i++) {
      blockTimePromises.push(this.provider.getBlock(logs[i].meta.blockNumber))
    }
    const blockTimes = await Promise.all(blockTimePromises)

    return logs.map((l, i) => {
      l.meta.timestamp = blockTimes[i].timestamp
      return l
    })
  }

  async addGasPriceToLogs(logs) {
    if (!logs) return logs

    const logsWithGasPriceOps = logs.map(async log => {
      const tx = await this.provider.getTransaction(log.meta.transactionHash)
      // eslint-disable-next-line require-atomic-updates
      log.meta.gasPrice = ethers.utils.formatUnits(tx.gasPrice, 'gwei')
      return log
    })

    return Promise.all(logsWithGasPriceOps)
  }
}
