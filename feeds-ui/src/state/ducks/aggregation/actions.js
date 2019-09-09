import * as types from './types'

export const requestOracles = () => ({
  type: types.ORACLES_REQUEST
})

export const successOracles = payload => ({
  type: types.ORACLES_SUCCESS,
  payload
})

export const setCurrentAnswer = payload => ({
  type: types.CURRENT_ANSWER,
  payload
})

export const setLatestCompletedAnswerId = payload => ({
  type: types.LATEST_COMPLETED_ANSWER_ID,
  payload
})

export const setPendingAnswerId = payload => ({
  type: types.PENDING_ANSWER_ID,
  payload
})

export const setNextAnswerId = payload => ({
  type: types.NEXT_ANSWER_ID,
  payload
})

export const setOracleResponse = payload => ({
  type: types.ORACLE_RESPONSE,
  payload
})

export const setRequestTime = payload => ({
  type: types.REQUEST_TIME,
  payload
})

export const setMinumumResponses = payload => ({
  type: types.MINUMUM_RESPONSES,
  payload
})

export const setUpdateHeight = payload => ({
  type: types.UPDATE_HEIGHT,
  payload
})

export const setAnswerHistory = payload => ({
  type: types.ANSWER_HISTORY,
  payload
})
