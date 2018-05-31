import {
  REQUEST_ACCOUNT_BALANCE,
  RECEIVE_ACCOUNT_BALANCE_SUCCESS,
  RECEIVE_ACCOUNT_BALANCE_ERROR
} from 'actions'

const initialState = {
  eth: '0',
  link: '0',
  fetching: false,
  networkError: false
}

export default (state = initialState, action = {}) => {
  switch (action.type) {
    case REQUEST_ACCOUNT_BALANCE:
      return Object.assign(
        {},
        state,
        {
          fetching: true,
          networkError: false
        }
      )
    case RECEIVE_ACCOUNT_BALANCE_SUCCESS:
      return Object.assign(
        {},
        state,
        {
          eth: action.eth,
          link: action.link,
          fetching: false,
          networkError: false
        }
      )
    case RECEIVE_ACCOUNT_BALANCE_ERROR:
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
