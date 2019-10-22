import { Action, Dispatch } from 'redux'
import { ThunkAction } from 'redux-thunk'
import httpStatus from 'http-status-codes'
import * as api from '../api'
import { State as AppState } from '../reducers'

export function signIn(
  username: string,
  password: string,
): ThunkAction<Promise<void>, AppState, void, Action<string>> {
  return (dispatch: Dispatch) => {
    return api.signIn(username, password).then(status => {
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

export function signOut(): ThunkAction<
  Promise<void>,
  AppState,
  void,
  Action<string>
> {
  return (dispatch: Dispatch) => {
    return api.signOut().then(() => {
      dispatch({ type: 'ADMIN_SIGNOUT_SUCCEEDED' })
    })
  }
}
