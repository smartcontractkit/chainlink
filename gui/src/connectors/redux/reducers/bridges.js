import {
  REQUEST_BRIDGES,
  RECEIVE_BRIDGES_SUCCESS,
  RECEIVE_BRIDGES_ERROR,
  REQUEST_BRIDGE,
  RECEIVE_BRIDGE_SUCCESS,
  RECEIVE_BRIDGE_ERROR
} from 'actions'

const initialState = {
  items: {},
  currentPage: [],
  count: 0,
  networkError: false
}

export default (state = initialState, action = {}) => {
  switch (action.type) {
    case REQUEST_BRIDGES:
      return Object.assign({}, state, { networkError: false })
    case RECEIVE_BRIDGES_SUCCESS: {
      const newItems = action.items.reduce((acc, i) => {
        acc[i.id] = i
        return acc
      }, {})

      return Object.assign({}, state, {
        items: Object.assign({}, state.items, newItems),
        currentPage: action.items.map(b => b.id),
        count: action.count,
        networkError: false
      })
    }
    case RECEIVE_BRIDGES_ERROR:
      return Object.assign({}, state, { networkError: !!action.networkError })
    case REQUEST_BRIDGE:
      return Object.assign({}, state, { networkError: false })
    case RECEIVE_BRIDGE_SUCCESS:
      return Object.assign({}, state, {
        items: Object.assign({}, state.items, { [action.item.id]: action.item })
      })
    case RECEIVE_BRIDGE_ERROR:
      return Object.assign({}, state, { networkError: !!action.networkError })
    default:
      return state
  }
}
