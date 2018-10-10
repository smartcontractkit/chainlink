import {
  REQUEST_BRIDGESPEC,
  RECEIVE_BRIDGESPEC_SUCCESS,
  RECEIVE_BRIDGESPEC_ERROR
} from 'actions'

const initialState = {
  name: '',
  url: '',
  confirmations: '0',
  fetching: false,
  networkError: false
}

export default (state = initialState, action = {}) => {
  switch (action.type) {
    case REQUEST_BRIDGESPEC:
      return Object.assign(
        {},
        state,
        {
          fetching: true,
          networkError: false
        }
      )
    case RECEIVE_BRIDGESPEC_SUCCESS:
      return Object.assign(
        {},
        state,
        {
          name: action.name,
          url: action.url,
          confirmations: action.confirmations,
          minimumContractPayment: action.minimumContractPayment,
          incomingToken: action.incomingToken,
          outgoingToken: action.outgoingToken,
          fetching: false,
          networkError: false
        }
      )
    case RECEIVE_BRIDGESPEC_ERROR:
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
