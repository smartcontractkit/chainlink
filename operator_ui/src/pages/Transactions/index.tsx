import React from 'react'
import { Route, Switch, useRouteMatch } from 'react-router-dom'

import { TransactionScreen } from 'screens/Transaction/TransactionScreen'
import { TransactionsScreen } from 'screens/Transactions/TransactionsScreen'

export const TransactionsPage = function () {
  const { path } = useRouteMatch()

  return (
    <Switch>
      <Route path={`${path}/:id`}>
        <TransactionScreen />
      </Route>

      <Route exact path={path}>
        <TransactionsScreen />
      </Route>
    </Switch>
  )
}
