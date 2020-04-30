import { FeedConfig } from 'config'

/**
 * listing/FETCH_FEEDS_BEGIN
 */
export interface FetchFeedsBeginAction {
  type: 'listing/FETCH_FEEDS_BEGIN'
}

export function fetchFeedsBegin(): FetchFeedsBeginAction {
  return {
    type: 'listing/FETCH_FEEDS_BEGIN',
  }
}

/**
 * listing/FETCH_FEEDS_SUCCESS
 */
export interface FetchFeedsSuccessAction {
  type: 'listing/FETCH_FEEDS_SUCCESS'
  payload: FeedConfig[]
}

export function fetchFeedsSuccess(
  payload: FeedConfig[],
): FetchFeedsSuccessAction {
  return {
    type: 'listing/FETCH_FEEDS_SUCCESS',
    payload,
  }
}

/**
 * listing/FETCH_FEEDS_ERROR
 */
export interface FetchFeedsErrorAction {
  type: 'listing/FETCH_FEEDS_ERROR'
  error: Error
}

export function fetchFeedsError(error: Error): FetchFeedsErrorAction {
  return {
    type: 'listing/FETCH_FEEDS_ERROR',
    error,
  }
}

/**
 * listing/FETCH_ANSWER_SUCCESS
 */
export interface ListingAnswer {
  answer: string
  config: FeedConfig
}

export interface FetchAnswerSuccessAction {
  type: 'listing/FETCH_ANSWER_SUCCESS'
  payload: ListingAnswer
}

export function fetchAnswerSuccess(
  payload: ListingAnswer,
): FetchAnswerSuccessAction {
  return {
    type: 'listing/FETCH_ANSWER_SUCCESS',
    payload,
  }
}

/**
 * listing/FETCH_HEALTH_PRICE_SUCCESS
 */
export interface HealthPrice {
  price: number
  config: FeedConfig
}

export interface FetchHealthPriceSuccessAction {
  type: 'listing/FETCH_HEALTH_PRICE_SUCCESS'
  payload: HealthPrice
}

export function fetchHealthPriceSuccess(
  payload: any,
): FetchHealthPriceSuccessAction {
  return {
    type: 'listing/FETCH_HEALTH_PRICE_SUCCESS',
    payload,
  }
}
