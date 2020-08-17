import 'core-js/stable/object/from-entries'
import { Reducer } from 'redux'
import { FeedConfig } from 'config'
import {
  FETCH_FEEDS_BEGIN,
  FETCH_FEEDS_SUCCESS,
  FETCH_FEEDS_ERROR,
  FETCH_ANSWER_SUCCESS,
  FETCH_HEALTH_PRICE_SUCCESS,
  FETCH_ANSWER_TIMESTAMP_SUCCESS,
  ListingActionTypes,
} from './types'

export interface HealthCheck {
  currentPrice: number
}

export interface State {
  loadingFeeds: boolean
  feedItems: Record<FeedConfig['contractAddress'], FeedConfig>
  feedOrder: Array<FeedConfig['contractAddress']>
  answers: Record<FeedConfig['contractAddress'], string>
  healthChecks: Record<FeedConfig['contractAddress'], HealthCheck>
  answersTimestamp: Record<FeedConfig['contractAddress'], number>
}

export const INITIAL_STATE: State = {
  loadingFeeds: false,
  feedItems: {},
  feedOrder: [],
  answers: {},
  answersTimestamp: {},
  healthChecks: {},
}

const reducer: Reducer<State, ListingActionTypes> = (
  state = INITIAL_STATE,
  action,
) => {
  switch (action.type) {
    case FETCH_FEEDS_BEGIN: {
      return {
        ...state,
        loadingFeeds: true,
      }
    }

    case FETCH_FEEDS_SUCCESS: {
      const listedFeeds: Array<[
        FeedConfig['contractAddress'],
        FeedConfig,
      ]> = action.payload
        .filter(f => f.listing)
        .map(f => [f.contractAddress, f])
      const feedItems = Object.fromEntries(listedFeeds)
      const feedOrder = listedFeeds.map(([a]) => a)

      return {
        ...state,
        loadingFeeds: false,
        feedItems,
        feedOrder,
      }
    }

    case FETCH_FEEDS_ERROR: {
      return {
        ...state,
        loadingFeeds: false,
      }
    }

    case FETCH_ANSWER_SUCCESS:
      return {
        ...state,
        answers: {
          ...state.answers,
          [action.payload.config.contractAddress]: action.payload.answer,
        },
      }

    case FETCH_ANSWER_TIMESTAMP_SUCCESS:
      return {
        ...state,
        answersTimestamp: {
          ...state.answersTimestamp,
          [action.payload.config.contractAddress]: action.payload.timestamp,
        },
      }

    case FETCH_HEALTH_PRICE_SUCCESS: {
      const { config, price } = action.payload
      const healthCheck: HealthCheck = { currentPrice: price }
      const healthChecks: Record<string, HealthCheck> = {
        ...state.healthChecks,
        [config.contractAddress]: healthCheck,
      }

      return {
        ...state,
        healthChecks,
      }
    }

    default:
      return state
  }
}

export default reducer
