import { ethers } from 'ethers'
import {
  getLogs,
  formatAnswer,
  decodeLog,
  createContract,
  createInfuraProvider,
} from './utils'
import _ from 'lodash'

export default class AggregatorContract {
  oracleAnswerEvent = {
    filter: {},
    listener: {},
  }
  answerIdInterval = null
  provider = null
  contract = null

  constructor(config, abi) {
    this.provider = createInfuraProvider(config.networkId)
    this.contract = createContract(config.contractAddress, this.provider, abi)
    this.address = config.contractAddress
    this.alive = true
    this.abi = abi
    this.config = config
  }

  kill() {
    try {
      if (!this.alive) return false
      clearInterval(this.answerIdInterval)
      this.removeListener(
        this.oracleAnswerEvent.filter,
        this.oracleAnswerEvent.listener,
      )
      this.contract = null
      this.address = null
      this.alive = false
      this.config = null
    } catch (error) {
      //
    }
  }

  removeListener(filter, eventListener) {
    if (!this.alive) return

    this.provider.removeListener(filter, eventListener)
  }

  async oracles() {
    const addresses = []

    for (let i = 0; i <= 45; i++) {
      try {
        const address = await this.contract.oracles(i)
        addresses.push(address)
      } catch (err) {
        break
      }
    }
    return addresses
  }

  async jobId(index) {
    try {
      const jobIds = await this.contract.jobIds(index)
      return ethers.utils.toUtf8String(jobIds)
    } catch {
      //
    }
  }

  async latestAnswer() {
    const latestAnswer = await this.contract.currentAnswer()
    return formatAnswer(
      latestAnswer,
      this.config.multiply,
      this.config.decimalPlaces,
    )
  }

  async latestAnswerTimestamp() {
    const updatedHeight = await this.contract.updatedHeight()
    const block = await this.provider.getBlock(Number(updatedHeight))
    return block.timestamp
  }

  async nextAnswerId() {
    if (!this.alive) return
    const answerCounter = await this.provider.getStorageAt(this.address, 13)
    const bigNumberify = ethers.utils.bigNumberify(answerCounter)
    return bigNumberify.toNumber()
  }

  async latestCompletedAnswer() {
    const completedAnswer = await this.contract.latestCompletedAnswer()
    return completedAnswer.toNumber()
  }

  async minimumAnswers() {
    const minimumAnswers = await this.contract.minimumResponses()
    return minimumAnswers.toNumber()
  }

  async addBlockTimestampToLogs(logs) {
    if (_.isEmpty(logs)) return logs

    const blockTimePromises = logs.map(log =>
      this.provider.getBlock(log.meta.blockNumber),
    )
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

  async oracleAnswerLogs({ answerId, fromBlock }) {
    const answerIdHex = ethers.utils.hexlify(answerId)

    const oracleAnswerByIdFilter = {
      ...this.contract.filters.ResponseReceived(null, answerIdHex, null),
      fromBlock,
      toBlock: 'latest',
    }

    const logs = await getLogs(
      {
        provider: this.provider,
        filter: oracleAnswerByIdFilter,
        eventInterface: this.contract.interface.events.ResponseReceived,
      },
      decodedLog => ({
        answerFormatted: formatAnswer(
          decodedLog.response,
          this.config.multiply,
          this.config.decimalPlaces,
        ),
        answer: Number(decodedLog.response),
        answerId: Number(decodedLog.answerId),
        sender: decodedLog.sender,
      }),
    )

    return logs
  }

  async chainlinkRequestedLogs(pastBlocks = 40) {
    const fromBlock = await this.provider
      .getBlockNumber()
      .then(b => b - pastBlocks)

    const chainlinkRequestedFilter = {
      ...this.contract.filters.ChainlinkRequested(null),
      fromBlock,
      toBlock: 'latest',
    }

    const logs = await getLogs({
      provider: this.provider,
      filter: chainlinkRequestedFilter,
      eventInterface: this.contract.interface.events.ChainlinkRequested,
    })

    return logs
  }

  async listenNextAnswerId(callback) {
    clearInterval(this.answerIdInterval)
    this.answerIdInterval = setInterval(async () => {
      try {
        const answerId = await this.nextAnswerId()
        return callback(answerId)
      } catch {
        console.error('Error: Failed to fetch nextAnswerId')
      }
    }, 4000)
  }

  async listenOracleAnswerEvent(callback) {
    if (!this.alive) return

    this.removeListener(
      this.oracleAnswerEvent.filter,
      this.oracleAnswerEvent.listener,
    )

    this.oracleAnswerEvent.filter = {
      ...this.contract.filters.ResponseReceived(null, null, null),
    }

    return this.provider.on(
      this.oracleAnswerEvent.filter,
      (this.oracleAnswerEvent.listener = async log => {
        const logged = decodeLog(
          {
            log,
            eventInterface: this.contract.interface.events.ResponseReceived,
          },
          decodedLog => ({
            answerFormatted: formatAnswer(
              decodedLog.response,
              this.config.multiply,
              this.config.decimalPlaces,
            ),
            answer: Number(decodedLog.response),
            answerId: Number(decodedLog.answerId),
            sender: decodedLog.sender,
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

  async answerUpdatedLogs({ fromBlock }) {
    const answerUpdatedFilter = {
      ...this.contract.filters.AnswerUpdated(null, null),
      fromBlock,
      toBlock: 'latest',
    }

    const logs = await getLogs(
      {
        provider: this.provider,
        filter: answerUpdatedFilter,
        eventInterface: this.contract.interface.events.AnswerUpdated,
      },
      decodedLog => ({
        answerFormatted: formatAnswer(
          decodedLog.current,
          this.config.multiply,
          this.config.decimalPlaces,
        ),
        answer: Number(decodedLog.current),
        answerId: Number(decodedLog.answerId),
      }),
    )

    return logs
  }
}
