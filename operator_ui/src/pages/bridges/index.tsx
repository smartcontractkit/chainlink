import React from 'react'
import { Route, Switch, useRouteMatch } from 'react-router-dom'

import { BridgeScreen } from '../../screens/Bridge/BridgeScreen'
import { BridgesScreen } from '../../screens/Bridges/BridgesScreen'
import { EditBridgeScreen } from '../../screens/EditBridge/EditBridgeScreen'
import { NewBridgeScreen } from '../../screens/NewBridge/NewBridgeScreen'

export const BridgesPage = function () {
  const { path } = useRouteMatch()

  return (
    <Switch>
      <Route exact path={`${path}/new`}>
        <NewBridgeScreen />
      </Route>

      <Route path={`${path}/:id/edit`}>
        <EditBridgeScreen />
      </Route>

      <Route path={`${path}/:id`}>
        <BridgeScreen />
      </Route>

      <Route exact path={path}>
        <BridgesScreen />
      </Route>
    </Switch>
  )
}
