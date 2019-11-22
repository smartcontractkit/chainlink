import { ethers } from 'ethers'
import { getLogs, formatEthPrice, decodeLog } from './utils'
import AggregationAbi from 'contracts/AggregationAbi.json'
import _ from 'lodash'

const infuraKey = process.env.REACT_APP_INFURA_KEY

const createInfuraProvider = (network = 'mainnet') => {
  const provider = new ethers.providers.JsonRpcProvider(
    `https://${network}.infura.io/v3/${infuraKey}`,
  )
  provider.pollingInterval = 8000

  return provider
}

const createContract = (address, provider) =>
  new ethers.Contract(address, AggregationAbi, provider)

export default class AggregationContract {
  oracleResponseEvent = {
    filter: {},
    listener: {},
  }
  answerIdInterval

  constructor(address, name, symbol, network) {
    this.provider = createInfuraProvider(network)
    this.contract = createContract(address, this.provider)
    this.name = name
    this.symbol = symbol
    this.address = address
    this.alive = true
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
      this.name = null
      this.symbol = null
      this.address = null
      this.alive = false
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
    } catch (error) {
      //
    }
  }

  async currentAnswer() {
    const currentAnswer = await this.contract.currentAnswer()
    return formatEthPrice(currentAnswer)
  }

  async updateHeight() {
    const updatedHeight = await this.contract.updatedHeight()
    const block = await this.provider.getBlock(updatedHeight.toNumber())
    return {
      block: updatedHeight.toNumber(),
      timestamp: block.timestamp,
    }
  }

  async nextAnswerId() {
    if (!this.alive) return
    const answerCounter = await this.provider.getStorageAt(this.address, 13)
    const bigNumberify = ethers.utils.bigNumberify(answerCounter)
    return bigNumberify.toNumber()
  }

  async latestCompletedAnswer() {
    const currentAnswer = await this.contract.latestCompletedAnswer()
    return currentAnswer.toNumber()
  }

  async minimumResponses() {
    const minimumResponses = await this.contract.minimumResponses()
    return minimumResponses.toNumber()
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
        return callback ? callback(logWithTimestamp[0]) : logWithTimestamp[0]
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
