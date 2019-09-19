import { Dispatch } from 'redux'
import httpStatus from 'http-status-codes'
import * as api from '../api'

export function fetchOperators() {
  return (dispatch: Dispatch<any>) => {
    api.getOperators().then(status => {
      if (status === httpStatus.OK) {
        dispatch({ type: 'FETCH_OPERATORS_SUCCEEDED', data: [] })
      } else if (status === httpStatus.UNAUTHORIZED) {
        dispatch({ type: 'ADMIN_SIGNOUT_SUCCEEDED' })
      }
    })
  }
}
