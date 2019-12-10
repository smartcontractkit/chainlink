import { ethers } from 'ethers'
import {
  getLogs,
  formatEthPrice,
  decodeLog,
  createContract,
  createInfuraProvider,
} from './utils'
import _ from 'lodash'

export default class AggregationContract {
  oracleResponseEvent = {
    filter: {},
    listener: {},
  }
  answerIdInterval
  provider
  contract

  constructor(options, abi) {
    this.provider = createInfuraProvider(options.network)
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
        this.oracleResponseEvent.filter,
        this.oracleResponseEvent.listener,
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

  async currentAnswer() {
    const currentAnswer = await this.contract.currentAnswer()
    return formatEthPrice(currentAnswer)
  }

  async updateHeight() {
    const updatedHeight = await this.contract.updatedHeight()
    const block = await this.provider.getBlock(Number(updatedHeight))
    return block.timestamp
  }

  async nextAnswerId() {
    if (!this.alive) return
    const answerCounter = await this.provider.getStorageAt(this.address, 13)
    const bigNumberify = ethers.utils.bigNumberify(answerCounter)
    return Number(bigNumberify)
  }

  async latestCompletedAnswer() {
    const currentAnswer = await this.contract.latestCompletedAnswer()
    return Number(currentAnswer)
  }

  async minimumResponses() {
    const minimumResponses = await this.contract.minimumResponses()
    return Number(minimumResponses)
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

  async oracleResponseLogs({ answerId, fromBlock }) {
    const answerIdHex = ethers.utils.hexlify(answerId)

    const oracleResponseByIdFilter = {
      ...this.contract.filters.ResponseReceived(null, answerIdHex, null),
      fromBlock,
      toBlock: 'latest',
    }

    const logs = await getLogs(
      {
        provider: this.provider,
        filter: oracleResponseByIdFilter,
        eventInterface: this.contract.interface.events.ResponseReceived,
      },
      decodedLog => ({
        responseFormatted: formatEthPrice(decodedLog.response),
        response: Number(decodedLog.response),
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
      const answerId = await this.nextAnswerId()
      return callback(answerId)
    }, 4000)
  }

  async listenOracleResponseEvent(callback) {
    if (!this.alive) return

    this.removeListener(
      this.oracleResponseEvent.filter,
      this.oracleResponseEvent.listener,
    )

    this.oracleResponseEvent.filter = {
      ...this.contract.filters.ResponseReceived(null, null, null),
    }

    return this.provider.on(
      this.oracleResponseEvent.filter,
      (this.oracleResponseEvent.listener = async log => {
        const logged = decodeLog(
          {
            log,
            eventInterface: this.contract.interface.events.ResponseReceived,
          },
          decodedLog => ({
            responseFormatted: formatEthPrice(decodedLog.response),
            response: Number(decodedLog.response),
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
        responseFormatted: formatEthPrice(decodedLog.current),
        response: Number(decodedLog.current),
        answerId: Number(decodedLog.answerId),
      }),
    )

    return logs
  }
}
