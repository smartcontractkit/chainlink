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
]

function ignore(action: Actions) {
  return (
    IGNORE_AGGREGATOR_TYPES.includes(action.type) ||
    !action.type.startsWith('aggregator')
  )
}

export default aggregationMiddleware
