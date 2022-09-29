import React from 'react'
import {
  Route,
  Switch,
  Redirect,
  BrowserRouter as Router,
} from 'react-router-dom'
import CssBaseline from '@material-ui/core/CssBaseline'
import Private from './Private'
import { useOperatorUiSelector } from 'reducers'
import SignIn from 'pages/SignIn'

const Layout = () => {
  const redirectTo = useOperatorUiSelector((state) => state.redirect.to)

  return (
    <Router>
      <CssBaseline />

      <Switch>
        <Route exact path="/signin">
          <SignIn />
        </Route>

        {redirectTo && <Redirect to={redirectTo} />}

        <Route component={Private} />
      </Switch>
    </Router>
  )
}

export default Layout
