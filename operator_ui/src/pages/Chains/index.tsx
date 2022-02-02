import React from 'react'
import { Route, Switch, useRouteMatch } from 'react-router-dom'

import ChainsNew from 'src/screens/Chain/New'
import { ChainScreen } from 'src/screens/Chain/ChainScreen'
import { ChainsScreen } from 'src/screens/Chains/ChainsScreen'

export const ChainsPage = function () {
  const { path } = useRouteMatch()

  return (
    <Switch>
      <Route path={`${path}/new`}>
        <ChainsNew />
      </Route>

      <Route path={`${path}/:id`}>
        <ChainScreen />
      </Route>

      <Route exact path={path}>
        <ChainsScreen />
      </Route>
    </Switch>
  )
}
