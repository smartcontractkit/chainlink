import { Middleware } from 'redux'
import { Actions } from 'state/actions'

const aggregationMiddleware: Middleware = store => next => (
  action: Actions,
) => {
  if (ignore(action)) {
    return next(action)
  }

  if (!store.getState().aggregator.contractAddress) {
    return
  }

  return next(action)
}

const IGNORE_AGGREGATOR_TYPES: Array<string> = [
  'aggregator/CLEAR_STATE',
  'aggregator/CONFIG',
  'aggregator/CONTRACT_ADDRESS',
  'aggregator/FETCH_FEED_BY_PAIR_BEGIN',
  'aggregator/FETCH_FEED_BY_PAIR_SUCCESS',
  'aggregator/FETCH_FEED_BY_PAIR_ERROR',
  'aggregator/FETCH_FEED_BY_ADDRESS_BEGIN',
  'aggregator/FETCH_FEED_BY_ADDRESS_SUCCESS',
  'aggregator/FETCH_FEED_BY_ADDRESS_ERROR',
]

function ignore(action: Actions) {
  return (
    IGNORE_AGGREGATOR_TYPES.includes(action.type) ||
    !action.type.startsWith('aggregator')
  )
}

export default aggregationMiddleware
