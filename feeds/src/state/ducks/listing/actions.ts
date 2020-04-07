import { ListingAnswer } from './reducers'

/**
 * listing/SET_ANSWERS
 */
export interface SetAnswersAction {
  type: 'listing/SET_ANSWERS'
  payload: any
}

export function setAnswers(payload: ListingAnswer[]): SetAnswersAction {
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
  payload: any
}

export function setHealthPrice(payload: any): SetHealthPriceAction {
  return {
    type: 'listing/SET_HEALTH_PRICE',
    payload,
  }
}
