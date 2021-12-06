import React from 'react'
import { Route, Switch, useRouteMatch } from 'react-router-dom'

import { JobScreen } from 'screens/Job/JobScreen'
import { JobsScreen } from 'screens/Jobs/JobsScreen'

export const JobsPage = function () {
  const { path } = useRouteMatch()

  return (
    <Switch>
      <Route path={`${path}/:id`}>
        <JobScreen />
      </Route>

      <Route exact path={path}>
        <JobsScreen />
      </Route>
    </Switch>
  )
}
