import {
  REQUEST_CONFIGURATION,
  RECEIVE_CONFIGURATION_SUCCESS,
  RECEIVE_CONFIGURATION_ERROR
} from 'actions'

const initialState = {
  config: {},
  fetching: false,
  networkError: false
}

export default (state = initialState, action = {}) => {
  switch (action.type) {
    case REQUEST_CONFIGURATION:
      return Object.assign(
        {},
        state,
        {
          fetching: true,
          networkError: false
        }
      )
    case RECEIVE_CONFIGURATION_SUCCESS:
      return Object.assign(
        {},
        state,
        {
          config: action.config,
          fetching: false,
          networkError: false
        }
      )
    case RECEIVE_CONFIGURATION_ERROR:
      return Object.assign(
        {},
        state,
        {
          fetching: false,
          networkError: !!action.networkError
        }
      )
    default:
      return state
  }
}
