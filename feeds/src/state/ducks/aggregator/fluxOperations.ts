import * as actions from './actions'
import _ from 'lodash'
import { ethers } from 'ethers'
import FluxAggregatorAbi from '../../../contracts/FluxAggregatorAbi.json'
import FluxAggregatorContract from '../../../contracts/FluxAggregatorContract'

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
      const payload = await contractInstance.latestRound()
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
      const payload = await contractInstance.latestTimestamp()
      dispatch(actions.setLatestAnswerTimestamp(payload))
      return payload
    } catch {
      console.error('Could not fetch latest answer timestamp')
    }
  }
}

const fetchOracleAnswersById = (request: any) => {
  return async (dispatch: any, getState: any) => {
    try {
      const currentLogs = getState().aggregator.oracleAnswers || []
      const logs = await contractInstance.submissionReceivedLogs(request)
      const withTimestamp = await contractInstance.addBlockTimestampToLogs(logs)
      const withGasAndTimeStamp = await contractInstance.addGasPriceToLogs(
        withTimestamp,
      )

      const uniquePayload = _.uniqBy(
        [...withGasAndTimeStamp, ...currentLogs],
        l => {
          return l.sender
        },
      )

      dispatch(actions.setOracleAnswers(uniquePayload))
    } catch {
      console.error('Could not fetch oracle answers')
    }
  }
}

const fetchLatestRequestTimestamp = (request: any) => {
  return async (dispatch: any) => {
    try {
      const logs = await contractInstance.newRoundLogs(request)
      const startedAt = logs.length && logs[logs.length - 1].startedAt
      dispatch(actions.setLatestRequestTimestamp(startedAt))
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

function fetchAnswerHistory(fromBlock: number) {
  return async (dispatch: any) => {
    try {
      const payload = await contractInstance.answerUpdatedLogs({ fromBlock })
      const uniquePayload = _.uniqBy(payload, (e: any) => {
        return e.answerId
      })

      dispatch(actions.setAnswerHistory(uniquePayload))
    } catch {
      console.error('Could not fetch answer history')
    }
  }
}

function initListeners() {
  return async (dispatch: any, getState: any) => {
    contractInstance.listenSubmissionReceivedEvent(async (responseLog: any) => {
      const { minimumAnswers } = getState().aggregator
      const oracleAnswers = getState().aggregator.oracleAnswers || []
      const updatedAnswers = oracleAnswers.map((response: any) => {
        return response.sender === responseLog.sender ? responseLog : response
      })

      dispatch(actions.setOracleAnswers(updatedAnswers))

      const latestIdAnswers = _.filter(updatedAnswers, {
        answerId: responseLog.answerId,
      })

      if (latestIdAnswers.length >= minimumAnswers) {
        fetchLatestAnswer()(dispatch)
        fetchLatestAnswerTimestamp()(dispatch)
      }
    })

    contractInstance.listenNewRoundEvent(async (responseLog: any) => {
      await fetchLatestCompletedAnswerId()(dispatch)
      dispatch(actions.setPendingAnswerId(responseLog.answerId))
      dispatch(actions.setLatestRequestTimestamp(responseLog.startedAt))
    })
  }
}

const initContract = (config: any) => {
  return async (dispatch: any, getState: any) => {
    dispatch(actions.clearState())

    try {
      if (contractInstance) {
        contractInstance.kill()
      }
    } catch {
      console.error('Could not close the contract instance')
    }

    try {
      ethers.utils.getAddress(config.contractAddress)
    } catch (error) {
      throw new Error('Wrong contract address')
    }

    dispatch(actions.setConfig(config))
    dispatch(actions.setContractAddress(config.contractAddress))

    contractInstance = new FluxAggregatorContract(config, FluxAggregatorAbi)

    // Oracle addresses

    await fetchOracleList()(dispatch, getState)

    // Minimum oracle responses

    fetchMinimumAnswers()(dispatch)

    // Set answer Id

    const reportingAnswerId = await contractInstance.reportingRound()
    dispatch(actions.setPendingAnswerId(reportingAnswerId))

    // Current answers

    await fetchLatestAnswerTimestamp()(dispatch)

    // Fetch previous answers

    const currentBlockNumber = await contractInstance.provider.getBlockNumber()
    const latestAnswerId = await contractInstance.latestRound()
    const fromBlock = currentBlockNumber <= 6700 ? 0 : currentBlockNumber - 6700 // ~6700 blocks per day

    await fetchOracleAnswersById({
      round: latestAnswerId,
      fromBlock,
    })(dispatch, getState)

    // Fetch latest answers

    await fetchOracleAnswersById({
      round: reportingAnswerId,
      fromBlock,
    })(dispatch, getState)

    /**
     * Oracle Latest Request Time
     * Used to calculate hearbeat countdown timer
     */

    if (config.heartbeat) {
      fetchLatestRequestTimestamp({
        round: reportingAnswerId,
        fromBlock,
      })(dispatch)
    }

    // Current answer

    fetchLatestAnswer()(dispatch)

    // initalise listeners

    initListeners()(dispatch, getState)

    if (config.history) {
      fetchAnswerHistory(fromBlock)(dispatch)
    }
  }
}

function clearState() {
  return async (dispatch: any) => {
    try {
      contractInstance.kill()
    } catch {
      console.error('Could not clear the contract')
    }

    dispatch(actions.clearState())
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

export { initContract, clearState, fetchJobId }
