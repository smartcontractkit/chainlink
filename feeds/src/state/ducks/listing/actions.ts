import { FeedConfig } from 'config'
import {
  FETCH_FEEDS_BEGIN,
  FETCH_FEEDS_SUCCESS,
  FETCH_FEEDS_ERROR,
  ListingAnswer,
  FETCH_ANSWER_SUCCESS,
  FETCH_HEALTH_PRICE_SUCCESS,
  ListingActionTypes,
  FETCH_ANSWER_TIMESTAMP_SUCCESS,
} from './types'

export function fetchFeedsBegin(): ListingActionTypes {
  return {
    type: FETCH_FEEDS_BEGIN,
  }
}

export function fetchFeedsSuccess(payload: FeedConfig[]): ListingActionTypes {
  return {
    type: FETCH_FEEDS_SUCCESS,
    payload,
  }
}

export function fetchFeedsError(error: Error): ListingActionTypes {
  return {
    type: FETCH_FEEDS_ERROR,
    error,
  }
}

export function fetchAnswerSuccess(payload: ListingAnswer): ListingActionTypes {
  return {
    type: FETCH_ANSWER_SUCCESS,
    payload,
  }
}

export function fetchHealthPriceSuccess(payload: any): ListingActionTypes {
  return {
    type: FETCH_HEALTH_PRICE_SUCCESS,
    payload,
  }
}

export function fetchAnswerTimestampSuccess(payload: any): ListingActionTypes {
  return {
    type: FETCH_ANSWER_TIMESTAMP_SUCCESS,
    payload,
  }
}
