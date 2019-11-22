import * as types from './types'

export const setOracles = payload => ({
  type: types.ORACLES,
  payload,
})

export const setCurrentAnswer = payload => ({
  type: types.CURRENT_ANSWER,
  payload,
})

export const setLatestCompletedAnswerId = payload => ({
  type: types.LATEST_COMPLETED_ANSWER_ID,
  payload,
})

export const setPendingAnswerId = payload => ({
  type: types.PENDING_ANSWER_ID,
  payload,
})

export const setNextAnswerId = payload => ({
  type: types.NEXT_ANSWER_ID,
  payload,
})

export const setOracleResponse = payload => ({
  type: types.ORACLE_RESPONSE,
  payload,
})

export const setRequestTime = payload => ({
  type: types.REQUEST_TIME,
  payload,
})

export const setMinumumResponses = payload => ({
  type: types.MINUMUM_RESPONSES,
  payload,
})

export const setUpdateHeight = payload => ({
  type: types.UPDATE_HEIGHT,
  payload,
})

export const setAnswerHistory = payload => ({
  type: types.ANSWER_HISTORY,
  payload,
})

export const setContractAddress = payload => ({
  type: types.CONTRACT_ADDRESS,
  payload,
})

export const setOptions = payload => ({
  type: types.OPTIONS,
  payload,
})

export const clearState = () => ({
  type: types.CLEAR_STATE,
})
