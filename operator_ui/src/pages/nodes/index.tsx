import React from 'react'
import { Route, Switch, useRouteMatch } from 'react-router-dom'

import { NodeScreen } from 'screens/Node/NodeScreen'
import { NodesScreen } from 'screens/Nodes/NodesScreen'

export const NodesPage = function () {
  const { path } = useRouteMatch()

  return (
    <Switch>
      <Route path={`${path}/:id`}>
        <NodeScreen />
      </Route>

      <Route exact path={path}>
        <NodesScreen />
      </Route>
    </Switch>
  )
}
