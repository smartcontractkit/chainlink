import React from 'react'
import { Route, useRouteMatch } from 'react-router-dom'

import { BridgeScreen } from '../../screens/Bridge/BridgeScreen'
import { BridgesScreen } from '../../screens/Bridges/BridgesScreen'

export const BridgesPage = function () {
  const { path } = useRouteMatch()

  return (
    <>
      <Route exact path={`${path}/:id`}>
        <BridgeScreen />
      </Route>

      <Route exact path={path}>
        <BridgesScreen />
      </Route>
    </>
  )
}
