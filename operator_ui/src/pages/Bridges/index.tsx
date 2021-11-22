import React from 'react'
import { Route, useRouteMatch } from 'react-router-dom'

import { BridgesScreen } from '../../screens/Bridges/BridgesScreen'

export const BridgesPage = function () {
  const { path } = useRouteMatch()

  return (
    <>
      <Route exact path={path}>
        <BridgesScreen />
      </Route>
    </>
  )
}
