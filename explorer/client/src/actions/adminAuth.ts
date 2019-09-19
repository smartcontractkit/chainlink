import { Dispatch } from 'redux'
import httpStatus from 'http-status-codes'
import * as api from '../api'

export function signIn(username: string, password: string) {
  return (dispatch: Dispatch<any>) => {
    api.signIn(username, password).then(status => {
      if (status === httpStatus.OK) {
        dispatch({ type: 'ADMIN_SIGNIN_SUCCEEDED' })
      } else if (status === httpStatus.UNAUTHORIZED) {
        dispatch({ type: 'ADMIN_SIGNIN_FAILED' })
        dispatch({
          type: 'NOTIFY_ERROR',
          text: 'Invalid username and password.',
        })
      } else {
        dispatch({ type: 'ADMIN_SIGNIN_ERROR' })
      }
    })
  }
}

export function signOut() {
  return (dispatch: Dispatch<any>) => {
    api.signOut().then(() => {
      dispatch({ type: 'ADMIN_SIGNOUT_SUCCEEDED' })
    })
  }
}
