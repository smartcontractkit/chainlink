import React from 'react'
import { Root, Routes } from 'react-static'
import universal from 'react-universal-component'
import { Route, Switch, Redirect } from 'react-router-dom'
import CssBaseline from '@material-ui/core/CssBaseline'
import PrivateRoute from './PrivateRoute'
import Private from './Private'
import Loading from 'components/Loading'
import { connect } from 'react-redux'
import { bindActionCreators } from 'redux'

// Asynchronously load routes that are chunked via code-splitting
// 'import' as a function must take a string. It can't take a variable.
const uniOpts = { loading: Loading }
const SignIn = universal(import('./containers/SignIn'), uniOpts)
const SignOut = universal(import('./containers/SignOut'), uniOpts)

class Layout extends React.Component {
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
      <Root>
        <CssBaseline />

        <Switch>
          <Route exact path='/signin' component={SignIn} />
          <PrivateRoute exact path='/signout' component={SignOut} />
          {redirectTo && <Redirect to={redirectTo} />}
          <Route component={Private} />
          <Route render={() => <Routes />} />
        </Switch>
      </Root>
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

export default ConnectedLayout
