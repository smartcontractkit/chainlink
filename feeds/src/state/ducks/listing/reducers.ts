import { FeedConfig } from 'config'
import { Reducer } from 'redux'
import { Actions } from 'state/actions'

export interface ListingAnswer {
  answer?: string
  config: FeedConfig
}

export interface HealthPrice {
  config: FeedConfig
  price: number
}

export interface HealthCheck {
  currentPrice: number
}

export interface State {
  answers: ListingAnswer[]
  healthChecks: Record<FeedConfig['contractAddress'], HealthCheck>
}

export const INITIAL_STATE: State = {
  answers: [],
  healthChecks: {},
}

const reducer: Reducer<State, Actions> = (state = INITIAL_STATE, action) => {
  switch (action.type) {
    case 'listing/SET_ANSWERS':
      return {
        ...state,
        answers: action.payload,
      }

    case 'listing/SET_HEALTH_PRICE': {
      const { config, price } = action.payload
      const healthCheck: HealthCheck = { currentPrice: price }
      const healthChecks: Record<string, HealthCheck> = {
        ...state.healthChecks,
        ...{ [config.contractAddress]: healthCheck },
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
