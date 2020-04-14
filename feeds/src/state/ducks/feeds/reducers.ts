import 'core-js/stable/object/from-entries'
import { FeedConfig, getFeedsConfig } from '../../../config'
import { Actions } from 'state/actions'
import { Networks } from 'utils'

export interface State {
  items: Record<FeedConfig['contractAddress'], FeedConfig>
  order: Array<FeedConfig['contractAddress']>
  pairPaths: [string, Networks, FeedConfig['contractAddress']][]
}

export const INITIAL_STATE: State = {
  items: Object.fromEntries(getFeedsConfig().map(f => [f.contractAddress, f])),
  order: getFeedsConfig().map(f => f.contractAddress),
  pairPaths: getFeedsConfig().map(f => [
    f.path,
    f.networkId,
    f.contractAddress,
  ]),
}

const reducer = (state: State = INITIAL_STATE, action: Actions) => {
  switch (action.type) {
    default:
      return state
  }
}

export default reducer
