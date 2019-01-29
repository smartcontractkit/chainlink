import React, { PureComponent } from 'react'
import CssBaseline from '@material-ui/core/CssBaseline'
import { bindActionCreators } from 'redux'
import { connect } from 'react-redux'
import { Router, Route, Switch } from 'react-static'
import { Redirect } from 'react-router'
import Routes from 'react-static-routes'
import { hot } from 'react-hot-loader'
import universal from 'react-universal-component'
import Loading from 'components/Loading'
import Private from './Private'
import PrivateRoute from './PrivateRoute'

const uniOpts = { loading: Loading }
const SignIn = universal(import('./containers/SignIn'), uniOpts)
const SignOut = universal(import('./containers/SignOut'), uniOpts)

class Layout extends PureComponent {
  // Remove the server-side injected CSS.
  componentDidMount () {
    const jssStyles = document.getElementById('jss-server-side')
    if (jssStyles && jssStyles.parentNode) {
      jssStyles.parentNode.removeChild(jssStyles)
    }
  }

  render () {
    const { redirectTo } = this.props

    return (
      <React.Fragment>
        <CssBaseline />

        <Router>
          <Switch>
            <Route exact path='/signin' component={SignIn} />
            <PrivateRoute exact path='/signout' component={SignOut} />
            {redirectTo && <Redirect to={redirectTo} />}
            <Route component={Private} />
            <Routes />
          </Switch>
        </Router>
      </React.Fragment>
    )
  }
}

const mapStateToProps = state => ({
  redirectTo: state.redirect.to
})

const mapDispatchToProps = dispatch => bindActionCreators(
  {},
  dispatch
)

export const ConnectedLayout = connect(
  mapStateToProps,
  mapDispatchToProps
)(Layout)

export default hot(module)(ConnectedLayout)
