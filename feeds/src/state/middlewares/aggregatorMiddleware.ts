import { Middleware } from 'redux'
import { Actions } from 'state/actions'

const aggregationMiddleware: Middleware = store => next => (
  action: Actions,
) => {
  if (ignore(action)) {
    return next(action)
  }

  if (!store.getState().aggregation.contractAddress) {
    return
  }

  return next(action)
}

const IGNORE_AGGREGATION_TYPES: Array<string> = [
  'aggregation/CLEAR_STATE',
  'aggregation/OPTIONS',
  'aggregation/CONTRACT_ADDRESS',
]

function ignore(action: Actions) {
  return (
    IGNORE_AGGREGATION_TYPES.includes(action.type) ||
    !action.type.startsWith('aggregation')
  )
}

export default aggregationMiddleware
