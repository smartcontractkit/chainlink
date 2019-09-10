import { HeartbeatContract, provider } from './contract'
import { ethers } from 'ethers'
import {
  formatEthPrice,
  getLogs,
  getLogsWithoutTimestamp,
  getLogsFromEvent,
  getLogsFromEventWithoutTimestamp,
  removeListener
} from './utils'

import { nextAnswerId } from './api'

export async function oracleResponseById(answerId, pastBlocks = 100) {
  const answerIdHex = ethers.utils.hexlify(answerId)

  const oracleResponseByIdFilter = {
    ...HeartbeatContract.filters.ResponseReceived(null, answerIdHex, null),
    fromBlock: provider.getBlockNumber().then(b => b - pastBlocks),
    toBlock: 'latest'
  }

  const logs = await getLogs({
    name: 'ResponseReceived',
    filter: oracleResponseByIdFilter,
    eventInterface: HeartbeatContract.interface.events.ResponseReceived,
    cb: decodedLog => ({
      responseFormatted: formatEthPrice(decodedLog.response),
      response: Number(decodedLog.response),
      answerId: Number(decodedLog.answerId),
      sender: decodedLog.sender
    })
  })

  return logs
}

let oracleResponseEventFilter
let oracleResponseEventListener

export function listenOracleResponseEvent(callback) {
  removeListener(oracleResponseEventFilter, oracleResponseEventListener)

  oracleResponseEventFilter = {
    ...HeartbeatContract.filters.ResponseReceived(null, null, null)
  }

  provider.on(
    oracleResponseEventFilter,
    (oracleResponseEventListener = async log => {
      const logged = await getLogsFromEvent({
        name: 'ResponseReceived',
        log,
        eventInterface: HeartbeatContract.interface.events.ResponseReceived,
        cb: decodedLog => ({
          responseFormatted: formatEthPrice(decodedLog.response),
          response: Number(decodedLog.response),
          answerId: Number(decodedLog.answerId),
          sender: decodedLog.sender
        })
      })

      return callback ? callback(logged) : logged
    })
  )
}

export async function chainlinkRequested(pastBlocks = 40) {
  const fromBlock = await provider.getBlockNumber().then(b => b - pastBlocks)

  const chainlinkRequestedFilter = {
    ...HeartbeatContract.filters.ChainlinkRequested(null),
    fromBlock,
    toBlock: 'latest'
  }

  const logs = await getLogs({
    name: 'ChainlinkRequested',
    filter: chainlinkRequestedFilter,
    eventInterface: HeartbeatContract.interface.events.ChainlinkRequested
  })

  return logs
}

let chainlinkRequestedEventFilter
let chainlinkRequestedEventListener

export function listenChainlinkRequestedEvent(callback) {
  removeListener(chainlinkRequestedEventFilter, chainlinkRequestedEventListener)

  chainlinkRequestedEventFilter = {
    ...HeartbeatContract.filters.ChainlinkRequested(null)
  }

  provider.once(
    chainlinkRequestedEventFilter,
    (chainlinkRequestedEventListener = async log => {
      const logged = await getLogsFromEvent({
        name: 'ChainlinkRequested',
        log,
        eventInterface: HeartbeatContract.interface.events.ChainlinkRequested
      })

      return callback ? callback(logged) : logged
    })
  )
}

let answerIdTimer

export function listenNextAnswerId(callback) {
  clearInterval(answerIdTimer)
  answerIdTimer = setInterval(async () => {
    const answerId = await nextAnswerId()
    return callback(answerId)
  }, 8000)
}

export async function answerUpdated(pastBlocks = 6700) {
  const answerUpdatedFilter = {
    ...HeartbeatContract.filters.AnswerUpdated(null, null),
    fromBlock: provider.getBlockNumber().then(b => b - pastBlocks),
    toBlock: 'latest'
  }

  const logs = await getLogsWithoutTimestamp({
    name: 'AnswerUpdated',
    filter: answerUpdatedFilter,
    eventInterface: HeartbeatContract.interface.events.AnswerUpdated,
    cb: decodedLog => ({
      responseFormatted: formatEthPrice(decodedLog.current),
      response: String(decodedLog.current),
      answerId: Number(decodedLog.answerId)
    })
  })

  return logs
}

let answerUpdatedEventFilter
let answerUpdatedEventListener

export function listenAnswerUpdatedEvent(callback) {
  removeListener(answerUpdatedEventFilter, answerUpdatedEventListener)

  answerUpdatedEventFilter = {
    ...HeartbeatContract.filters.AnswerUpdated(null, null)
  }

  provider.on(
    answerUpdatedEventFilter,
    (answerUpdatedEventListener = async log => {
      const logged = await getLogsFromEventWithoutTimestamp({
        name: 'AnswerUpdated',
        log,
        eventInterface: HeartbeatContract.interface.events.AnswerUpdated,
        cb: decodedLog => ({
          responseFormatted: formatEthPrice(decodedLog.current),
          response: String(decodedLog.current),
          answerId: Number(decodedLog.answerId)
        })
      })

      return callback ? callback(logged) : logged
    })
  )
}
