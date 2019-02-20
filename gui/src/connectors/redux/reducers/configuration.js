import {
  REQUEST_CONFIGURATION,
  RECEIVE_CONFIGURATION_SUCCESS,
  RECEIVE_CONFIGURATION_ERROR
} from 'actions'

const initialState = {
  config: {},
  networkError: false
}

export default (state = initialState, action = {}) => {
  switch (action.type) {
    case REQUEST_CONFIGURATION:
      return Object.assign({}, state, { networkError: false })
    case RECEIVE_CONFIGURATION_SUCCESS:
      return Object.assign({}, state, {
        config: action.config,
        networkError: false
      })
    case RECEIVE_CONFIGURATION_ERROR:
      return Object.assign({}, state, { networkError: !!action.networkError })
    default:
      return state
  }
}
