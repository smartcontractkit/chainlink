import { FeedConfig } from 'config'

export const FETCH_FEEDS_BEGIN = 'listing/FETCH_FEEDS_BEGIN'
export const FETCH_FEEDS_SUCCESS = 'listing/FETCH_FEEDS_SUCCESS'
export const FETCH_FEEDS_ERROR = 'listing/FETCH_FEEDS_ERROR'
export const FETCH_ANSWER_SUCCESS = 'listing/FETCH_ANSWER_SUCCESS'
export const FETCH_HEALTH_PRICE_SUCCESS = 'listing/FETCH_HEALTH_PRICE_SUCCESS'

export interface FetchFeedsBeginAction {
  type: typeof FETCH_FEEDS_BEGIN
}

export interface FetchFeedsSuccessAction {
  type: typeof FETCH_FEEDS_SUCCESS
  payload: FeedConfig[]
}

export interface FetchFeedsErrorAction {
  type: typeof FETCH_FEEDS_ERROR
  error: Error
}

export interface ListingAnswer {
  answer: string
  config: FeedConfig
}

export interface FetchAnswerSuccessAction {
  type: typeof FETCH_ANSWER_SUCCESS
  payload: ListingAnswer
}

export interface HealthPrice {
  price: number
  config: FeedConfig
}

export interface FetchHealthPriceSuccessAction {
  type: typeof FETCH_HEALTH_PRICE_SUCCESS
  payload: HealthPrice
}

export type ListingActionTypes =
  | FetchFeedsBeginAction
  | FetchFeedsSuccessAction
  | FetchFeedsErrorAction
  | FetchAnswerSuccessAction
  | FetchHealthPriceSuccessAction
