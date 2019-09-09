import * as actions from './actions'

import {
  oracleAddresses,
  currentAnswer,
  latestCompletedAnswer,
  nextAnswerId,
  minimumResponses,
  updateHeight
} from 'state/contract/api'

import {
  oracleResponseById,
  chainlinkRequested,
  listenOracleResponseEvent,
  listenNextAnswerId
} from 'state/contract/events'

const fetchOracles = () => {
  return async (dispatch, getState) => {
    if (getState().aggregation.oracles) {
      return
    }
    dispatch(actions.requestOracles())
    try {
      let payload = await oracleAddresses()
      dispatch(actions.successOracles(payload))
    } catch (error) {}
  }
}

const fetchLatestCompletedAnswerId = () => {
  return async (dispatch, getState) => {
    try {
      let payload = await latestCompletedAnswer()
      dispatch(actions.setLatestCompletedAnswerId(payload))
      return payload
    } catch (error) {}
  }
}

const fetchNextAnswerId = () => {
  return async (dispatch, getState) => {
    try {
      let payload = await nextAnswerId()
      dispatch(actions.setNextAnswerId(payload))
      return payload
    } catch (error) {}
  }
}

const fetchCurrentAnswer = () => {
  return async (dispatch, getState) => {
    try {
      let payload = await currentAnswer()
      dispatch(actions.setCurrentAnswer(payload))
    } catch (error) {}
  }
}

const fetchUpdateHeight = () => {
  return async (dispatch, getState) => {
    try {
      let payload = await updateHeight()
      dispatch(actions.setUpdateHeight(payload))
    } catch (error) {}
  }
}

const fetchOracleResponseById = answerId => {
  return async (dispatch, getState) => {
    try {
      let payload = await oracleResponseById(answerId)
      dispatch(actions.setOracleResponse(payload))
    } catch (error) {}
  }
}

const fetchRequestTime = () => {
  return async (dispatch, getState) => {
    try {
      let payload = await chainlinkRequested()
      const latestBlock = payload[payload.length - 1]
      dispatch(actions.setRequestTime(latestBlock.meta.timestamp))
    } catch (error) {}
  }
}

const fetchMinimumResponses = () => {
  return async (dispatch, getState) => {
    try {
      let payload = await minimumResponses()
      dispatch(actions.setMinumumResponses(payload))
    } catch (error) {}
  }
}

// TODO
// const fetchAnswerHistory = () => {
//   return async (dispatch, getState) => {
//     try {
//       let payload = await answerUpdated()

//       const uniquePayload = _.uniqBy(payload, e => {
//         return e.answerId
//       })

//       const formattedPayload = uniquePayload.map(e => ({
//         answerId: e.answerId,
//         response: e.response,
//         responseFormatted: e.responseFormatted,
//         blockNumber: e.meta.blockNumber,
//         timestamp: e.meta.timestamp
//       }))

//       dispatch(actions.setAnswerHistory(formattedPayload))
//     } catch (error) {}
//   }
// }

const fetchInitData = () => {
  return async (dispatch, getState) => {
    await fetchMinimumResponses()(dispatch)
    await fetchOracles()(dispatch, getState)

    const nextAnswerId = await fetchNextAnswerId()(dispatch)
    fetchOracleResponseById(nextAnswerId - 1)(dispatch)
    fetchRequestTime()(dispatch, getState)
    fetchLatestCompletedAnswerId()(dispatch)
    fetchCurrentAnswer()(dispatch)
    fetchUpdateHeight()(dispatch)
    initListeners()(dispatch, getState)
  }
}

const initListeners = () => {
  return async (dispatch, getState) => {
    /**
     * Listen to next answer id
     * - change next answer id
     * - reset oracle data
     * - reset request time
     */

    listenNextAnswerId(async responseNextAnswerId => {
      const { nextAnswerId } = getState().aggregation

      if (responseNextAnswerId > nextAnswerId) {
        dispatch(actions.setNextAnswerId(responseNextAnswerId))
        dispatch(actions.setOracleResponse([]))
        fetchRequestTime()(dispatch, getState)
      }
    })

    /**
     * Listen to oracles response
     * - compare answerId
     * - add unique oracles response data
     */

    listenOracleResponseEvent(async responseLog => {
      const { nextAnswerId, minimumResponses } = getState().aggregation

      if (responseLog.answerId === nextAnswerId - 1) {
        const storeLogs = getState().aggregation.oracleResponse || []
        const uniqueLogs = storeLogs.filter(l => {
          return l.meta.transactionHash !== responseLog.meta.transactionHash
        })
        dispatch(actions.setOracleResponse([...uniqueLogs, ...[responseLog]]))

        if (uniqueLogs.length >= minimumResponses) {
          fetchLatestCompletedAnswerId()(dispatch)
          fetchCurrentAnswer()(dispatch)
          fetchUpdateHeight()(dispatch)
        }
      }
    })
  }
}

export { fetchInitData, fetchOracles, fetchLatestCompletedAnswerId }
