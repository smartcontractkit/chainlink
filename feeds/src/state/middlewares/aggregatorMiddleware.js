import * as aggregationTypes from '../ducks/aggregation/types'

const aggregationMiddleware = store => next => action => {
  const ignoreWhen =
    action.type === aggregationTypes.CLEAR_STATE ||
    action.type === aggregationTypes.OPTIONS ||
    action.type === aggregationTypes.CONTRACT_ADDRESS ||
    !action.type.startsWith('aggregation')

  if (ignoreWhen) {
    return next(action)
  }

  if (!store.getState().aggregation.contractAddress) {
    return
  }

  return next(action)
}

export default aggregationMiddleware
