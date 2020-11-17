import React from 'react'
import universal from 'react-universal-component'
import {
  Route,
  Switch,
  Redirect,
  BrowserRouter as Router,
} from 'react-router-dom'
import CssBaseline from '@material-ui/core/CssBaseline'
import Private from './Private'
import Loading from 'components/Loading'
import { useOperatorUiSelector } from 'reducers'

// Asynchronously load routes that are chunked via code-splitting
// 'import' as a function must take a string. It can't take a variable.
const uniOpts = { loading: Loading }
const SignIn = universal(import('./pages/SignIn'), uniOpts)

const Layout = () => {
  // Remove the server-side injected CSS.
  const jssStyles = document.getElementById('jss-server-side')

  React.useEffect(() => {
    if (jssStyles && jssStyles.parentNode) {
      jssStyles.parentNode.removeChild(jssStyles)
    }
  }, [jssStyles])

  const redirectTo = useOperatorUiSelector((state) => state.redirect.to)

  return (
    <Router>
      <CssBaseline />

      <Switch>
        <Route exact path="/signin" component={SignIn} />
        {redirectTo && <Redirect to={redirectTo} />}
        <Route component={Private} />
      </Switch>
    </Router>
  )
}

export default Layout
