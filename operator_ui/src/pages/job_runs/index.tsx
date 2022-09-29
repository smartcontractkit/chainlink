import React from 'react'
import { Route, Switch, useRouteMatch } from 'react-router-dom'

import { JobRunScreen } from 'screens/JobRun/JobRunScreen'
import { JobRunsScreen } from 'screens/JobRuns/JobRunsScreen'

export const JobRunsPage = function () {
  const { path } = useRouteMatch()

  return (
    <Switch>
      <Route exact path={path}>
        <JobRunsScreen />
      </Route>

      <Route path={`${path}/:id`}>
        <JobRunScreen />
      </Route>
    </Switch>
  )
}
