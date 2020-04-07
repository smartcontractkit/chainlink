import { FeedConfig, getFeedsConfig } from 'config'
import 'core-js/stable/object/from-entries'
import { Actions } from 'state/actions'

export interface State {
  items: Record<FeedConfig['contractAddress'], FeedConfig>
  order: string[]
}

export const INITIAL_STATE: State = {
  items: Object.fromEntries(getFeedsConfig().map(f => [f.contractAddress, f])),
  order: getFeedsConfig().map(f => f.contractAddress),
}

const reducer = (state: State = INITIAL_STATE, action: Actions) => {
  switch (action.type) {
    default:
      return state
  }
}

export default reducer
