import { FeedConfig } from 'config'
import { ethers } from 'ethers'
import _ from 'lodash'
import FluxAggregatorAbi from '../../../contracts/FluxAggregatorAbi.json'
import FluxAggregatorContract from '../../../contracts/FluxAggregatorContract'
import * as actions from './actions'

export default class FluxOperations {
  static contractInstance: any

  static fetchOracleList() {
    return async (dispatch: any, getState: any) => {
      if (getState().aggregator.oracleList) {
        return
      }
      try {
        const payload = await FluxOperations.contractInstance.oracles()
        dispatch(actions.setOracleList(payload))
      } catch {
        console.error('Could not fetch oracle list ')
      }
    }
  }

  static fetchLatestCompletedAnswerId() {
    return async (dispatch: any) => {
      try {
        const payload = await FluxOperations.contractInstance.latestRound()
        dispatch(actions.setLatestCompletedAnswerId(payload))
        return payload
      } catch {
        console.error('Could not fetch latest completed answer id')
      }
    }
  }

  static fetchLatestAnswer() {
    return async (dispatch: any) => {
      try {
        const payload = await FluxOperations.contractInstance.latestAnswer()
        dispatch(actions.setLatestAnswer(payload))
      } catch {
        console.error('Could not fetch latest answer')
      }
    }
  }

  static fetchLatestAnswerTimestamp() {
    return async (dispatch: any) => {
      try {
        const payload = await FluxOperations.contractInstance.latestTimestamp()
        dispatch(actions.setLatestAnswerTimestamp(payload))
        return payload
      } catch {
        console.error('Could not fetch latest answer timestamp')
      }
    }
  }

  static fetchOracleAnswersById(request: any) {
    return async (dispatch: any, getState: any) => {
      try {
        const currentLogs = getState().aggregator.oracleAnswers
        const logs = await FluxOperations.contractInstance.submissionReceivedLogs(
          request,
        )
        const withTimestamp = await FluxOperations.contractInstance.addBlockTimestampToLogs(
          logs,
        )
        const withGasAndTimeStamp = await FluxOperations.contractInstance.addGasPriceToLogs(
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

  static fetchLatestRequestTimestamp = (request: any) => {
    return async (dispatch: any) => {
      try {
        const logs = await FluxOperations.contractInstance.newRoundLogs(request)
        const startedAt = logs?.[logs.length - 1].startedAt
        dispatch(actions.setLatestRequestTimestamp(startedAt))
      } catch {
        console.error('Could not fetch request time')
      }
    }
  }

  static fetchMinimumAnswers() {
    return async (dispatch: any) => {
      try {
        const payload = await FluxOperations.contractInstance.minimumAnswers()
        dispatch(actions.setMinumumAnswers(payload))
      } catch {
        console.error('Could not fetch minimum answers')
      }
    }
  }

  static fetchAnswerHistory(fromBlock: number) {
    return async (dispatch: any) => {
      try {
        const payload = await FluxOperations.contractInstance.answerUpdatedLogs(
          {
            fromBlock,
          },
        )
        const uniquePayload = _.uniqBy(payload, (e: any) => {
          return e.answerId
        })

        dispatch(actions.setAnswerHistory(uniquePayload))
      } catch {
        console.error('Could not fetch answer history')
      }
    }
  }

  static initListeners() {
    return async (dispatch: any, getState: any) => {
      FluxOperations.contractInstance.listenSubmissionReceivedEvent(
        async (responseLog: any) => {
          const { minimumAnswers } = getState().aggregator
          const oracleAnswers = getState().aggregator.oracleAnswers || []
          const updatedAnswers = oracleAnswers.map((response: any) => {
            return response.sender === responseLog.sender
              ? responseLog
              : response
          })

          dispatch(actions.setOracleAnswers(updatedAnswers))

          const latestIdAnswers = _.filter(updatedAnswers, {
            answerId: responseLog.answerId,
          })

          if (latestIdAnswers.length >= minimumAnswers) {
            FluxOperations.fetchLatestAnswer()(dispatch)
            FluxOperations.fetchLatestAnswerTimestamp()(dispatch)
          }
        },
      )

      FluxOperations.contractInstance.listenNewRoundEvent(
        async (responseLog: any) => {
          await FluxOperations.fetchLatestCompletedAnswerId()(dispatch)
          dispatch(actions.setPendingAnswerId(responseLog.answerId))
          dispatch(actions.setLatestRequestTimestamp(responseLog.startedAt))
        },
      )
    }
  }

  /**
   * Initialise aggregator contract and fill the store with all necessery data for a visualisation page.
   * @param config FeedsConfig
   */

  static initContract(config: FeedConfig) {
    return async (dispatch: any, getState: any) => {
      try {
        FluxOperations.contractInstance?.kill()
      } catch {
        console.error('Could not close the contract instance')
      }

      try {
        ethers.utils.getAddress(config.contractAddress)
      } catch (error) {
        throw new Error('Wrong contract address')
      }

      dispatch(actions.setContractAddress(config.contractAddress))

      FluxOperations.contractInstance = new FluxAggregatorContract(
        config,
        FluxAggregatorAbi,
      )

      // Oracle addresses
      await FluxOperations.fetchOracleList()(dispatch, getState)

      // Minimum oracle responses
      FluxOperations.fetchMinimumAnswers()(dispatch)

      // Set answer Id
      const reportingAnswerId = await FluxOperations.contractInstance.reportingRound()
      dispatch(actions.setPendingAnswerId(reportingAnswerId))

      // Current answers
      await FluxOperations.fetchLatestAnswerTimestamp()(dispatch)

      // Fetch previous answers
      const currentBlockNumber = await FluxOperations.contractInstance.provider.getBlockNumber()
      const latestAnswerId = await FluxOperations.contractInstance.latestRound()
      const fromBlock =
        currentBlockNumber <= 6700 ? 0 : currentBlockNumber - 6700 // ~6700 blocks per day

      await FluxOperations.fetchOracleAnswersById({
        round: latestAnswerId,
        fromBlock,
      })(dispatch, getState)

      // Fetch latest answers
      await FluxOperations.fetchOracleAnswersById({
        round: reportingAnswerId,
        fromBlock,
      })(dispatch, getState)

      /**
       * Oracle Latest Request Time
       * Used to calculate hearbeat countdown timer
       */
      if (config.heartbeat) {
        FluxOperations.fetchLatestRequestTimestamp({
          round: reportingAnswerId,
          fromBlock,
        })(dispatch)
      }

      // Current answer
      FluxOperations.fetchLatestAnswer()(dispatch)

      // initalise listeners
      FluxOperations.initListeners()(dispatch, getState)

      if (config.history) {
        FluxOperations.fetchAnswerHistory(fromBlock)(dispatch)
      }
    }
  }

  static fetchJobId(address: any) {
    return async (_dispatch: any, getState: any) => {
      const { oracleList } = getState().aggregator
      try {
        const index = oracleList.indexOf(address)
        return FluxOperations.contractInstance.jobId(index)
      } catch {
        console.error('Could not fetch a job id')
      }
    }
  }
}
