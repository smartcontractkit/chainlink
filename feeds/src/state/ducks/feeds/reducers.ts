import 'core-js/stable/object/from-entries'
import { FeedConfig } from 'feeds'
import { Actions } from 'state/actions'
import feeds from '../../../feeds.json'

export interface State {
  items: Record<FeedConfig['contractAddress'], FeedConfig>
  order: string[]
}

export const INITIAL_STATE: State = {
  items: Object.fromEntries(feeds.map(f => [f.contractAddress, f])),
  order: feeds.map(f => f.contractAddress),
}

const reducer = (state: State = INITIAL_STATE, action: Actions) => {
  switch (action.type) {
    default:
      return state
  }
}

export default reducer
