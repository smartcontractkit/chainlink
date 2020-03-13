import * as types from './types'

export const setAnswers = payload => ({
  type: types.SET_ANSWERS,
  payload,
})

export const setHealthPrice = payload => ({
  type: types.SET_HEALTH_PRICE,
  payload,
})
