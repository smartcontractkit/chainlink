import reducer from './reducers'

import * as aggregatorActions from './actions'
import * as aggregatorOperations from './operations'
import fluxAggregatorOperations from './fluxOperations'
import * as aggregatorSelectors from './selectors'

export {
  aggregatorActions,
  aggregatorOperations,
  aggregatorSelectors,
  fluxAggregatorOperations,
}

export default reducer
