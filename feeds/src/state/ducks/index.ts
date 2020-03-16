export { default as aggregation } from './aggregation'
export { default as networkGraph } from './networkGraph'
export { default as listing } from './listing'

import { State as AggregationState } from './aggregation/reducers'
import { State as NetworkGraphState } from './aggregation/reducers'
import { State as LisitingState } from './aggregation/reducers'

export type Ducks = {
  aggregation: AggregationState
  listing: LisitingState
  networkGraph: NetworkGraphState
}
