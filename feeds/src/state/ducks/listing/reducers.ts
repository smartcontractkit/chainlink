import 'core-js/stable/object/from-entries'
import { Reducer } from 'redux'
import { FeedConfig } from 'config'
import { Actions } from 'state/actions'

export interface HealthCheck {
  currentPrice: number
}

export interface State {
  loadingFeeds: boolean
  feedItems: Record<FeedConfig['contractAddress'], FeedConfig>
  feedOrder: Array<FeedConfig['contractAddress']>
  answers: Record<FeedConfig['contractAddress'], string>
  healthChecks: Record<FeedConfig['contractAddress'], HealthCheck>
}

export const INITIAL_STATE: State = {
  loadingFeeds: false,
  feedItems: {},
  feedOrder: [],
  answers: {},
  healthChecks: {},
}

const reducer: Reducer<State, Actions> = (state = INITIAL_STATE, action) => {
  switch (action.type) {
    case 'listing/FETCH_FEEDS_BEGIN': {
      return {
        ...state,
        loadingFeeds: true,
      }
    }

    case 'listing/FETCH_FEEDS_SUCCESS': {
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

    case 'listing/FETCH_FEEDS_ERROR': {
      return {
        ...state,
        loadingFeeds: false,
      }
    }

    case 'listing/FETCH_ANSWER_SUCCESS':
      return {
        ...state,
        answers: {
          ...state.answers,
          [action.payload.config.contractAddress]: action.payload.answer,
        },
      }

    case 'listing/FETCH_HEALTH_PRICE_SUCCESS': {
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
