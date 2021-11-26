import React from 'react'
import { Route, Switch, useRouteMatch } from 'react-router-dom'

import { JobsScreen } from '../../screens/Jobs/JobsScreen'

export const JobsPage = function () {
  const { path } = useRouteMatch()

  return (
    <Switch>
      <Route exact path={path}>
        <JobsScreen />
      </Route>
    </Switch>
  )
}
