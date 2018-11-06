import React, { Component } from 'react'
import Routes from 'react-static-routes'
import CssBaseline from '@material-ui/core/CssBaseline'
import Grid from '@material-ui/core/Grid'
import PrivateRoute from './PrivateRoute'
import Header from 'containers/Header'
import Loading from 'components/Loading'
import Notifications from 'containers/Notifications'
import universal from 'react-universal-component'
import { Redirect } from 'react-router'
import { Router, Route, Switch } from 'react-static'
import { hot } from 'react-hot-loader'
import { withStyles } from '@material-ui/core/styles'
import { connect } from 'react-redux'
import { bindActionCreators } from 'redux'

// Asynchronously load routes that are chunked via code-splitting
// 'import' as a function must take a string. It can't take a variable.
const uniOpts = {loading: Loading}
const DashboardsIndex = universal(import('./containers/Dashboards/Index'), uniOpts)
const JobsIndex = universal(import('./containers/Jobs/Index'), uniOpts)
const JobsShow = universal(import('./containers/Jobs/Show'), uniOpts)
const JobsNew = universal(import('./containers/Jobs/New'), uniOpts)
const BridgesIndex = universal(import('./containers/Bridges/Index'), uniOpts)
const BridgesNew = universal(import('./containers/Bridges/New'), uniOpts)
const BridgesShow = universal(import('./containers/Bridges/Show'), uniOpts)
const BridgesEdit = universal(import('./containers/Bridges/Edit'), uniOpts)
const JobRunsIndex = universal(import('./containers/JobRuns/Index'), uniOpts)
const JobRunsShow = universal(import('./containers/JobRuns/Show'), uniOpts)
const Configuration = universal(import('./containers/Configuration'), uniOpts)
const About = universal(import('./containers/About'), uniOpts)
const SignIn = universal(import('./containers/SignIn'), uniOpts)
const SignOut = universal(import('./containers/SignOut'), uniOpts)

const styles = theme => {
  return {
    content: {
      margin: theme.spacing.unit * 5,
      marginTop: 0
    }
  }
}

class Layout extends Component {
  state = {headerHeight: 0}

  onHeaderResize = (_width, height) => {
    this.setState({headerHeight: height})
  }

  render () {
    const {classes, redirectTo} = this.props

    return (
      <Router>
        <Grid container>
          <CssBaseline />
          <Grid item xs={12}>
            <Header
              onResize={this.onHeaderResize}
              drawerContainer={this.drawerContainer}
            />

            <main
              ref={ref => { this.drawerContainer = ref }}
              style={{paddingTop: this.state.headerHeight}}
            >
              <Notifications />

              <div className={classes.content}>
                <Switch>
                  <Route exact path='/signin' component={SignIn} />
                  <PrivateRoute exact path='/signout' component={SignOut} />
                  {redirectTo && <Redirect to={redirectTo} />}
                  <PrivateRoute
                    exact
                    path='/'
                    render={props => <DashboardsIndex {...props} recentlyCreatedPageSize={4} />}
                  />
                  <PrivateRoute exact path='/jobs' component={JobsIndex} />
                  <PrivateRoute exact path='/jobs/page/:jobPage' component={JobsIndex} />
                  <PrivateRoute exact path='/jobs/new' component={JobsNew} />
                  <PrivateRoute
                    exact
                    path='/jobs/:jobSpecId'
                    render={props => <JobsShow {...props} showJobRunsCount={5} />}
                  />
                  <PrivateRoute exact path='/jobs/:jobSpecId/runs' component={JobRunsIndex} />
                  <PrivateRoute exact path='/jobs/:jobSpecId/runs/page/:jobRunsPage' component={JobRunsIndex} />
                  <PrivateRoute exact path='/jobs/:jobSpecId/runs/id/:jobRunId' component={JobRunsShow} />
                  <PrivateRoute exact path='/bridges' component={BridgesIndex} />
                  <PrivateRoute exact path='/bridges/page/:bridgePage' component={BridgesIndex} />
                  <PrivateRoute exact path='/bridges/new' component={BridgesNew} />
                  <PrivateRoute exact path='/bridges/:bridgeId' component={BridgesShow} />
                  <PrivateRoute exact path='/bridges/:bridgeId/edit' component={BridgesEdit} />
                  <PrivateRoute exact path='/about' component={About} />
                  <PrivateRoute exact path='/config' component={Configuration} />
                  <Routes />
                </Switch>
              </div>
            </main>
          </Grid>
        </Grid>
      </Router>
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

export const ConnectedLayout = connect(mapStateToProps, mapDispatchToProps)(Layout)

export default hot(module)(withStyles(styles)(ConnectedLayout))
