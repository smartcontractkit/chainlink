import * as actions from './actions'
import _ from 'lodash'
import moment from 'moment'
import { ethers } from 'ethers'

import AggregationContract from 'contracts/AggregationContract'

let contractInstance

const fetchOracles = () => {
  return async (dispatch, getState) => {
    if (getState().aggregation.oracles) {
      return
    }
    try {
      const payload = await contractInstance.oracles()
      dispatch(actions.setOracles(payload))
    } catch (error) {
      //
    }
  }
}

const fetchLatestCompletedAnswerId = () => {
  return async dispatch => {
    try {
      const payload = await contractInstance.latestCompletedAnswer()
      dispatch(actions.setLatestCompletedAnswerId(payload))
      return payload
    } catch (error) {
      //
    }
  }
}

const fetchCurrentAnswer = () => {
  return async dispatch => {
    try {
      const payload = await contractInstance.currentAnswer()
      dispatch(actions.setCurrentAnswer(payload))
    } catch (error) {
      //
    }
  }
}

const fetchUpdateHeight = () => {
  return async dispatch => {
    try {
      const payload = await contractInstance.updateHeight()
      dispatch(actions.setUpdateHeight(payload))
      return payload || {}
    } catch (error) {
      return {}
    }
  }
}

const fetchOracleResponseById = request => {
  return async (dispatch, getState) => {
    try {
      const currentLogs = getState().aggregation.oracleResponse || []

      const logs = await contractInstance.oracleResponseLogs(request)
      const withTimestamp = await contractInstance.addBlockTimestampToLogs(logs)

      const uniquePayload = _.uniqBy([...withTimestamp, ...currentLogs], l => {
        return l.sender
      })

      dispatch(actions.setOracleResponse(uniquePayload))
    } catch (error) {
      //
    }
  }
}

const fetchRequestTime = () => {
  return async dispatch => {
    try {
      const logs = await contractInstance.chainlinkRequestedLogs()
      const latestLog = logs.length && logs[logs.length - 1].meta.blockNumber

      const block = latestLog
        ? await contractInstance.provider.getBlock(latestLog)
        : null

      dispatch(actions.setRequestTime(block ? block.timestamp : null))
    } catch (error) {
      //
    }
  }
}

const fetchMinimumResponses = () => {
  return async dispatch => {
    try {
      const payload = await contractInstance.minimumResponses()
      dispatch(actions.setMinumumResponses(payload))
    } catch (error) {
      //
    }
  }
}

const fetchAnswerHistory = () => {
  return async dispatch => {
    try {
      const fromBlock = await contractInstance.provider
        .getBlockNumber()
        .then(b => b - 6700) // 6700 block is ~24 hours

      const payload = await contractInstance.answerUpdatedLogs({ fromBlock })
      const uniquePayload = _.uniqBy(payload, e => {
        return e.answerId
      })
      const withTimestamp = await contractInstance.addBlockTimestampToLogs(
        uniquePayload,
      )

      const formattedPayload = withTimestamp.map(e => ({
        answerId: e.answerId,
        response: e.response,
        responseFormatted: e.responseFormatted,
        blockNumber: e.meta.blockNumber,
        timestamp: e.meta.timestamp,
      }))

      dispatch(actions.setAnswerHistory(formattedPayload))
    } catch (error) {
      //
    }
  }
}

const initListeners = () => {
  return async (dispatch, getState) => {
    /**
     * Listen to next answer id
     * - change next answer id
     * - reset oracle data
     * - reset request time (hardcode current time)
     */

    contractInstance.listenNextAnswerId(async responseNextAnswerId => {
      const { nextAnswerId } = getState().aggregation
      if (responseNextAnswerId > nextAnswerId) {
        dispatch(actions.setNextAnswerId(responseNextAnswerId))
        dispatch(actions.setPendingAnswerId(responseNextAnswerId - 1))
        dispatch(actions.setRequestTime(moment().unix()))
      }
    })

    /**
     * Listen to oracles response
     * - compare answerId
     * - add unique oracles response data
     */

    contractInstance.listenOracleResponseEvent(async responseLog => {
      const { nextAnswerId, minimumResponses } = getState().aggregation

      if (responseLog.answerId === nextAnswerId - 1) {
        const storeLogs = getState().aggregation.oracleResponse || []
        const uniqueLogs = storeLogs.filter(l => {
          return l.meta.transactionHash !== responseLog.meta.transactionHash
        })

        const updateLogs = uniqueLogs.map(l => {
          return l.sender === responseLog.sender ? responseLog : l
        })

        const senderIndex = _.findIndex(uniqueLogs, {
          sender: responseLog.sender,
        })

        if (senderIndex < 0) {
          updateLogs.push(responseLog)
        }

        dispatch(actions.setOracleResponse(updateLogs))

        const responseNumber = _.filter(updateLogs, {
          answerId: responseLog.answerId,
        })

        if (responseNumber.length >= minimumResponses) {
          fetchLatestCompletedAnswerId()(dispatch)
          fetchCurrentAnswer()(dispatch)
          fetchUpdateHeight()(dispatch)
        }
      }
    })
  }
}

const initContract = options => {
  return async (dispatch, getState) => {
    dispatch(actions.clearState())

    try {
      contractInstance.kill()
    } catch (error) {
      //
    }

    try {
      ethers.utils.getAddress(options.contractAddress)
    } catch (error) {
      throw new Error('Wrong contract address')
    }

    dispatch(actions.setOptions(options))
    dispatch(actions.setContractAddress(options.contractAddress))

    contractInstance = new AggregationContract(
      options.contractAddress,
      options.name,
      options.valuePrefix,
      options.network,
    )

    // Oracle addresses

    await fetchOracles()(dispatch, getState)

    // Minimum oracle responses

    fetchMinimumResponses()(dispatch)

    // Set answer Id

    const nextAnswerId = await contractInstance.nextAnswerId()
    dispatch(actions.setNextAnswerId(nextAnswerId))
    dispatch(actions.setPendingAnswerId(nextAnswerId - 1))

    // Current answers

    const height = await fetchUpdateHeight()(dispatch)

    // Fetch previous responses (counter / block time + 10)

    await fetchOracleResponseById({
      answerId: nextAnswerId - 2,
      fromBlock:
        height.block - options.counter
          ? Math.round(options.counter / 13 + 30)
          : 100,
    })(dispatch, getState)

    // Takes last block height minus 10 blocks (to make sure we get all the reqests)

    fetchOracleResponseById({
      answerId: nextAnswerId - 1,
      fromBlock: height.block - 10,
    })(dispatch, getState)

    // Initial request time

    fetchRequestTime()(dispatch)

    // Latest completed answer id

    fetchLatestCompletedAnswerId()(dispatch)

    // Current answer and block height

    fetchCurrentAnswer()(dispatch)

    // initalise listeners

    initListeners()(dispatch, getState)

    if (options.history) {
      fetchAnswerHistory()(dispatch)
    }
  }
}

const clearState = () => {
  return async dispatch => {
    try {
      contractInstance.kill()
    } catch (error) {
      //
    }

    dispatch(actions.clearState())
  }
}

const fetchJobId = address => {
  return async (dispatch, getState) => {
    const { oracles } = getState().aggregation
    try {
      const index = oracles.indexOf(address)
      return contractInstance.jobId(index)
    } catch (error) {
      //
    }
  }
}

export { initContract, clearState, fetchJobId }
