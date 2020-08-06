import * as jsonapi from '@chainlink/json-api-client'
import { Dispatch } from 'redux'
import _ from 'lodash'
import moment from 'moment'
import { ethers } from 'ethers'
import { FeedConfig, OracleNode, Config } from '../../../config'
import { Networks } from '../../../utils'
import * as actions from './actions'
import AggregatorAbi from '../../../contracts/AggregatorAbi.json'
import AggregatorAbiV2 from '../../../contracts/AggregatorAbi.v2.json'
import AggregatorContract from '../../../contracts/AggregatorContract'
import AggregatorContractV2 from '../../../contracts/AggregatorContractV2'

/**
 * feed
 */
const NETWORK_PATHS: Record<string, Networks> = {
  ropsten: Networks.ROPSTEN,
  mainnet: Networks.MAINNET,
  kovan: Networks.KOVAN,
}

export function fetchFeedByPair(pairPath: string, networkPath = 'mainnet') {
  return async (dispatch: Dispatch) => {
    dispatch(actions.fetchFeedByPairBegin())

    jsonapi
      .fetchWithTimeout(Config.feedsJson(), {})
      .then((r: Response) => r.json())
      .then((json: FeedConfig[]) => {
        const networkId = NETWORK_PATHS[networkPath] ?? Networks.MAINNET
        const feed = json.find(
          f => f.path === pairPath && f.networkId === networkId,
        )

        if (feed) {
          dispatch(actions.fetchFeedByPairSuccess(feed))
        } else {
          dispatch(actions.fetchFeedByPairError('Not Found'))
        }
      })
      .catch(e => {
        dispatch(actions.fetchFeedByPairError(e.toString()))
      })
  }
}

export function fetchFeedByAddress(contractAddress: string) {
  return async (dispatch: Dispatch) => {
    dispatch(actions.fetchFeedByAddressBegin())

    jsonapi
      .fetchWithTimeout(Config.feedsJson(), {})
      .then((r: Response) => r.json())
      .then((json: FeedConfig[]) => {
        const feed = json.find(f => f.contractAddress === contractAddress)

        if (feed) {
          dispatch(actions.fetchFeedByAddressSuccess(feed))
        } else {
          dispatch(actions.fetchFeedByAddressError('Not Found'))
        }
      })
      .catch(e => {
        dispatch(actions.fetchFeedByAddressError(e.toString()))
      })
  }
}

/**
 * oracle nodes
 */
export function fetchOracleNodes() {
  return async (dispatch: Dispatch) => {
    dispatch(actions.fetchOracleNodesBegin())

    jsonapi
      .fetchWithTimeout(Config.nodesJson(), {})
      .then((r: Response) => r.json())
      .then((json: OracleNode[]) => {
        dispatch(actions.fetchOracleNodesSuccess(json))
      })
      .catch(e => {
        dispatch(actions.fetchFeedByPairError(e.toString()))
      })
  }
}

/**
 * oracles
 */
let contractInstance: any

function fetchOracleList() {
  return async (dispatch: any, getState: any) => {
    if (getState().aggregator.oracleList) {
      return
    }
    try {
      const payload = await contractInstance.oracles()
      dispatch(actions.setOracleList(payload))
    } catch {
      console.error('Could not fetch oracle list ')
    }
  }
}

function fetchLatestCompletedAnswerId() {
  return async (dispatch: any) => {
    try {
      const payload = await contractInstance.latestCompletedAnswer()
      dispatch(actions.setLatestCompletedAnswerId(payload))
      return payload
    } catch {
      console.error('Could not fetch latest completed answer id')
    }
  }
}

function fetchLatestAnswer() {
  return async (dispatch: any) => {
    try {
      const payload = await contractInstance.latestAnswer()
      dispatch(actions.setLatestAnswer(payload))
    } catch {
      console.error('Could not fetch latest answer')
    }
  }
}

function fetchLatestAnswerTimestamp() {
  return async (dispatch: any) => {
    try {
      const payload = await contractInstance.latestAnswerTimestamp()
      dispatch(actions.setLatestAnswerTimestamp(payload))
      return payload || {}
    } catch {
      console.error('Could not fetch latest answer timestamp')
    }
  }
}

const fetchOracleAnswersById = (request: any) => {
  return async (dispatch: any, getState: any) => {
    try {
      const currentLogs = getState().aggregator.oracleAnswers

      const logs = await contractInstance.oracleAnswerLogs(request)
      const withTimestamp = await contractInstance.addBlockTimestampToLogs(logs)
      const withGasAndTimeStamp = await contractInstance.addGasPriceToLogs(
        withTimestamp,
      )

      const uniquePayload = _.uniqBy(
        [...withGasAndTimeStamp, ...currentLogs],
        l => l.sender,
      )

      dispatch(actions.setOracleAnswers(uniquePayload))
    } catch {
      console.error('Could not fetch oracle answers')
    }
  }
}

const fetchLatestRequestTimestamp = (config: FeedConfig) => {
  return async (dispatch: any) => {
    try {
      // calculate last update time
      const pastBlocks = config.heartbeat
        ? Math.floor(config.heartbeat / 13)
        : 40
      const logs = await contractInstance.chainlinkRequestedLogs(pastBlocks)
      const latestLog = logs.length && logs[logs.length - 1].meta.blockNumber

      const block = latestLog
        ? await contractInstance.provider.getBlock(latestLog)
        : null

      dispatch(
        actions.setLatestRequestTimestamp(block ? block.timestamp : null),
      )
    } catch {
      console.error('Could not fetch request time')
    }
  }
}

function fetchMinimumAnswers() {
  return async (dispatch: any) => {
    try {
      const payload = await contractInstance.minimumAnswers()
      dispatch(actions.setMinumumAnswers(payload))
    } catch {
      console.error('Could not fetch minimum answers')
    }
  }
}

function fetchAnswerHistory(config: FeedConfig) {
  return async (dispatch: any) => {
    try {
      const fromBlock = await contractInstance.provider
        .getBlockNumber()
        .then((b: any) => b - 6700 * (config.historyDays ?? 1)) // 6700 block is ~24 hours

      const payload = await contractInstance.answerUpdatedLogs({ fromBlock })
      const uniquePayload = _.uniqBy(payload, (e: any) => {
        return e.answerId
      })

      let history

      if (contractInstance.config.contractVersion === 2) {
        history = uniquePayload
      } else {
        const withTimestamp = await contractInstance.addBlockTimestampToLogs(
          uniquePayload,
        )

        history = withTimestamp.map((e: any) => ({
          answerId: e.answerId,
          answer: e.answer,
          answerFormatted: e.answerFormatted,
          timestamp: e.meta.timestamp,
        }))
      }

      dispatch(actions.setAnswerHistory(history))
    } catch {
      console.error('Could not fetch answer history')
    }
  }
}

function initListeners() {
  return async (dispatch: any, getState: any) => {
    /**
     * Listen to next answer id
     * - change next answer id
     * - reset oracle data
     * - reset request time (hardcode current time)
     */

    contractInstance.listenNextAnswerId(async (responseNextAnswerId: any) => {
      const { nextAnswerId } = getState().aggregator
      if (responseNextAnswerId > nextAnswerId) {
        dispatch(actions.setNextAnswerId(responseNextAnswerId))
        dispatch(actions.setPendingAnswerId(responseNextAnswerId - 1))

        // reset hearbeat countdown timer
        dispatch(actions.setLatestRequestTimestamp(moment().unix()))
      }
    })

    /**
     * Listen to oracles response
     * - compare answerId
     * - add unique oracles response data
     */

    contractInstance.listenOracleAnswerEvent(async (responseLog: any) => {
      const { nextAnswerId, minimumAnswers } = getState().aggregator

      if (responseLog.answerId === nextAnswerId - 1) {
        const storeLogs = getState().aggregator.oracleAnswers || []
        const uniqueLogs = storeLogs.filter((l: any) => {
          return l.meta.transactionHash !== responseLog.meta.transactionHash
        })

        const updateLogs = uniqueLogs.map((l: any) =>
          l.sender === responseLog.sender ? responseLog : l,
        )

        const senderIndex = _.findIndex(uniqueLogs, {
          sender: responseLog.sender,
        })

        if (senderIndex < 0) {
          updateLogs.push(responseLog)
        }

        dispatch(actions.setOracleAnswers(updateLogs))

        const latestIdAnswers = _.filter(updateLogs, {
          answerId: responseLog.answerId,
        })

        if (latestIdAnswers.length >= minimumAnswers) {
          fetchLatestCompletedAnswerId()(dispatch)
          fetchLatestAnswer()(dispatch)
          fetchLatestAnswerTimestamp()(dispatch)
        }
      }
    })
  }
}

const initContract = (config: FeedConfig) => {
  return async (dispatch: any, getState: any) => {
    try {
      contractInstance?.kill()
    } catch {
      console.error('Could not close the contract instance')
    }

    try {
      ethers.utils.getAddress(config.contractAddress)
    } catch (error) {
      throw new Error('Wrong contract address')
    }

    dispatch(actions.setContractAddress(config.contractAddress))

    if (config.contractVersion === 2) {
      contractInstance = new AggregatorContractV2(config, AggregatorAbiV2)
    } else {
      contractInstance = new AggregatorContract(config, AggregatorAbi)
    }

    // Oracle addresses
    await fetchOracleList()(dispatch, getState)

    // Minimum oracle responses
    fetchMinimumAnswers()(dispatch)

    // Set answer Id
    const nextAnswerId = await contractInstance.nextAnswerId()
    dispatch(actions.setNextAnswerId(nextAnswerId))
    dispatch(actions.setPendingAnswerId(nextAnswerId - 1))

    // Current answers
    await fetchLatestAnswerTimestamp()(dispatch)

    // Fetch previous answers
    const currentBlockNumber = await contractInstance.provider.getBlockNumber()

    await fetchOracleAnswersById({
      answerId: nextAnswerId - 2,
      fromBlock: currentBlockNumber - 6700, // ~6700 blocks per day
    })(dispatch, getState)

    // Fetch latest answers
    fetchOracleAnswersById({
      answerId: nextAnswerId - 1,
      fromBlock: currentBlockNumber - 6700,
    })(dispatch, getState)

    /**
     * Oracle Latest Request Time
     * Used to calculate hearbeat countdown timer
     */
    if (config.heartbeat) {
      fetchLatestRequestTimestamp(config)(dispatch)
    }

    // Latest completed answer id
    fetchLatestCompletedAnswerId()(dispatch)

    // Current answer and block height
    fetchLatestAnswer()(dispatch)

    // initialise listeners
    initListeners()(dispatch, getState)

    if (config.history) {
      fetchAnswerHistory(config)(dispatch)
    }
  }
}

const fetchJobId = (address: any) => {
  return async (_dispatch: any, getState: any) => {
    const { oracleList } = getState().aggregator
    try {
      const index = oracleList.indexOf(address)
      return contractInstance.jobId(index)
    } catch {
      console.error('Could not fetch a job id')
    }
  }
}

function fetchEthGasPrice() {
  return async (dispatch: any) => {
    try {
      const data = await fetch('https://ethgasstation.info/json/ethgasAPI.json')
      const jsonData = await data.json()
      dispatch(actions.setEthGasPrice(jsonData))
    } catch {
      console.error('Could not fetch gas price')
    }
  }
}

function clearContract() {
  return async (dispatch: any) => {
    try {
      dispatch(actions.clearState())
      contractInstance?.kill()
    } catch {
      console.error('Could not close the contract instance')
    }
  }
}

export { initContract, fetchJobId, fetchEthGasPrice, clearContract }
