import {
  RECEIVE_SESSION_FAIL
} from 'actions'

const initialState = []
const SIGN_IN_FAIL_MSG = 'Your email or password are incorrect. Please try again'

export default (state = initialState, action = {}) => {
  switch (action.type) {
    case RECEIVE_SESSION_FAIL:
      return [SIGN_IN_FAIL_MSG]
    default:
      return state
  }
}
