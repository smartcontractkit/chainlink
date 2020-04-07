import { HealthPrice, ListingAnswer } from './reducers'

/**
 * listing/SET_ANSWERS
 */
export interface SetAnswersAction {
  type: 'listing/SET_ANSWERS'
  payload: ListingAnswer[]
}

export function setAnswers(
  payload: SetAnswersAction['payload'],
): SetAnswersAction {
  return {
    type: 'listing/SET_ANSWERS',
    payload,
  }
}

/**
 * listing/SET_HEALTH_PRICE
 */
export interface SetHealthPriceAction {
  type: 'listing/SET_HEALTH_PRICE'
  payload: HealthPrice
}

export function setHealthPrice(
  payload: SetHealthPriceAction['payload'],
): SetHealthPriceAction {
  return {
    type: 'listing/SET_HEALTH_PRICE',
    payload,
  }
}
