import {
  REQUEST_ACCOUNT_BALANCE,
  RECEIVE_ACCOUNT_BALANCE_SUCCESS,
  RECEIVE_ACCOUNT_BALANCE_ERROR
} from 'actions'

const initialState = {
  eth: null,
  link: null,
  networkError: false
}

export default (state = initialState, action = {}) => {
  switch (action.type) {
    case REQUEST_ACCOUNT_BALANCE:
      return Object.assign(
        {},
        state,
        {networkError: false}
      )
    case RECEIVE_ACCOUNT_BALANCE_SUCCESS:
      return Object.assign(
        {},
        state,
        {
          eth: action.eth,
          link: action.link,
          networkError: false
        }
      )
    case RECEIVE_ACCOUNT_BALANCE_ERROR:
      return Object.assign(
        {},
        state,
        {networkError: !!action.networkError}
      )
    default:
      return state
  }
}
