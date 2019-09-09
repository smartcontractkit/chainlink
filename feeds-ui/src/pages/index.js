import React from 'react'
import { Switch, Route, Redirect, BrowserRouter } from 'react-router-dom'

import NetworkPage from './NetworkPage'

const AppRoutes = () => {
  return (
    <BrowserRouter>
      <Switch>
        <Route exact path="/" component={NetworkPage} />
        <Redirect to="/" />
      </Switch>
    </BrowserRouter>
  )
}

export default AppRoutes
