import React from 'react'
import { Route, RouteProps, Redirect, useLocation } from 'react-router-dom'
import { useDispatch } from 'react-redux'
import { useOperatorUiSelector } from 'reducers'
import { RouterActionType } from 'reducers/actions'

export const PrivateRoute = (props: RouteProps) => {
  const dispatch = useDispatch()
  const { pathname } = useLocation()

  React.useEffect(() => {
    dispatch({
      type: RouterActionType.MATCH_ROUTE,
      pathname,
    })
  }, [dispatch, pathname])

  const authenticated = useOperatorUiSelector(
    (state) => state.authentication.allowed,
  )

  if (authenticated) {
    return <Route {...props} />
  }

  return <Redirect to="/signin" />
}

export default PrivateRoute
