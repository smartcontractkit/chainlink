import * as actions from './actions'
import _ from 'lodash'
import moment from 'moment'
import { ethers } from 'ethers'
import AggregationAbi from 'contracts/AggregationAbi.json'
import AggregationAbiV2 from 'contracts/AggregationAbi.v2.json'
import AggregationContract from 'contracts/AggregationContract'
import AggregationContractV2 from 'contracts/AggregationContractV2'

let contractInstance

const fetchOracles = () => {
  return async (dispatch, getState) => {
    if (getState().aggregation.oracles) {
      return
    }
    try {
      const payload = await contractInstance.oracles()
      dispatch(actions.setOracles(payload))
    } catch {
      console.error('Could not fetch oracles ')
    }
  }
}

const fetchLatestCompletedAnswerId = () => {
  return async dispatch => {
    try {
      const payload = await contractInstance.latestCompletedAnswer()
      dispatch(actions.setLatestCompletedAnswerId(payload))
      return payload
    } catch {
      console.error('Could not fetch latest completed answer id')
    }
  }
}

const fetchCurrentAnswer = () => {
  return async dispatch => {
    try {
      const payload = await contractInstance.currentAnswer()
      dispatch(actions.setCurrentAnswer(payload))
    } catch {
      console.error('Could not fetch latest completed answer id')
    }
  }
}

const fetchUpdateHeight = () => {
  return async dispatch => {
    try {
      const payload = await contractInstance.updateHeight()
      dispatch(actions.setUpdateHeight(payload))
      return payload || {}
    } catch {
      console.error('Could not fetch update height')
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
    } catch {
      console.error('Could not fetch oracle responses')
    }
  }
}

const fetchRequestTime = options => {
  return async dispatch => {
    try {
      // calculate last update time
      const pastBlocks = options.counter ? Math.floor(options.counter / 13) : 40
      const logs = await contractInstance.chainlinkRequestedLogs(pastBlocks)
      const latestLog = logs.length && logs[logs.length - 1].meta.blockNumber

      const block = latestLog
        ? await contractInstance.provider.getBlock(latestLog)
        : null

      dispatch(actions.setRequestTime(block ? block.timestamp : null))
    } catch {
      console.error('Could not fetch request time')
    }
  }
}

const fetchMinimumResponses = () => {
  return async dispatch => {
    try {
      const payload = await contractInstance.minimumResponses()
      dispatch(actions.setMinumumResponses(payload))
    } catch {
      console.error('Could not fetch minimum responses')
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

      let history

      if (contractInstance.options.contractVersion === 2) {
        history = uniquePayload
      } else {
        const withTimestamp = await contractInstance.addBlockTimestampToLogs(
          uniquePayload,
        )

        history = withTimestamp.map(e => ({
          answerId: e.answerId,
          response: e.response,
          responseFormatted: e.responseFormatted,
          timestamp: e.meta.timestamp,
        }))
      }

      dispatch(actions.setAnswerHistory(history))
    } catch {
      console.error('Could not fetch answer history')
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
      if (contractInstance) {
        contractInstance.kill()
      }
    } catch {
      console.error('Could not close the contract instance')
    }

    try {
      ethers.utils.getAddress(options.contractAddress)
    } catch (error) {
      throw new Error('Wrong contract address')
    }

    dispatch(actions.setOptions(options))
    dispatch(actions.setContractAddress(options.contractAddress))

    if (options.contractVersion === 2) {
      contractInstance = new AggregationContractV2(options, AggregationAbiV2)
    } else {
      contractInstance = new AggregationContract(options, AggregationAbi)
    }

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

    fetchRequestTime(options)(dispatch)

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
    } catch {
      console.error('Could not clear state')
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
    } catch {
      console.error('Could not fetch a job id')
    }
  }
}

export { initContract, clearState, fetchJobId }
