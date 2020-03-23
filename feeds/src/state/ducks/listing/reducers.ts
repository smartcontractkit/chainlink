import { ListingAnswer } from './operations'
import { Actions } from 'state/actions'

export interface HealthCheck {
  currentPrice: number
}

export interface State {
  answers?: ListingAnswer[]
  healthChecks: Record<string, HealthCheck>
}

export const INITIAL_STATE: State = {
  answers: undefined,
  healthChecks: {},
}

const reducer = (state: State = INITIAL_STATE, action: Actions) => {
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
