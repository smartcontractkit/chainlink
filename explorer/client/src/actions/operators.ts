import { Action, Dispatch } from 'redux'
import { ThunkAction } from 'redux-thunk'
import httpStatus from 'http-status-codes'
import { State as AppState } from '../reducers'
import * as api from '../api'

export function fetchOperators(): ThunkAction<
  Promise<void>,
  AppState,
  void,
  Action<string>
> {
  return (dispatch: Dispatch<any>) => {
    return api.getOperators().then(status => {
      if (status === httpStatus.OK) {
        dispatch({ type: 'FETCH_OPERATORS_SUCCEEDED', data: [] })
      } else if (status === httpStatus.UNAUTHORIZED) {
        dispatch({ type: 'ADMIN_SIGNOUT_SUCCEEDED' })
      }
    })
  }
}
