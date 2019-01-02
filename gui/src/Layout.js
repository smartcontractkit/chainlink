import React from 'react'
import { Root } from 'react-static'
import CssBaseline from '@material-ui/core/CssBaseline'
import Grid from '@material-ui/core/Grid'
import PrivateRoute from './PrivateRoute'
import Header from 'containers/Header'
import Loading from 'components/Loading'
import Notifications from 'containers/Notifications'
import universal from 'react-universal-component'
import { Route, BrowserRouter, Switch, Redirect } from 'react-router-dom'
import { withStyles } from '@material-ui/core/styles'
import { connect } from 'react-redux'
import { bindActionCreators } from 'redux'
import { useHooks, useState } from 'use-react-hooks'

// Asynchronously load routes that are chunked via code-splitting
// 'import' as a function must take a string. It can't take a variable.
const uniOpts = { loading: Loading }
const DashboardsIndex = universal(import('./containers/Dashboards/Index'), uniOpts)
const JobsIndex = universal(import('./containers/Jobs/Index'), uniOpts)
const JobsShow = universal(import('./containers/Jobs/Show'), uniOpts)
const JobsDefinition = universal(import('./containers/Jobs/Definition'), uniOpts)
const JobsNew = universal(import('./containers/Jobs/New'), uniOpts)
const BridgesIndex = universal(import('./containers/Bridges/Index'), uniOpts)
const BridgesNew = universal(import('./containers/Bridges/New'), uniOpts)
const BridgesShow = universal(import('./containers/Bridges/Show'), uniOpts)
const BridgesEdit = universal(import('./containers/Bridges/Edit'), uniOpts)
const JobRunsIndex = universal(import('./containers/JobRuns/Index'), uniOpts)
const JobRunsShow = universal(import('./containers/JobRuns/Show'), uniOpts)
const JobRunsShowJson = universal(import('./containers/JobRuns/ShowJson'), uniOpts)
const Configuration = universal(import('./containers/Configuration'), uniOpts)
const SignIn = universal(import('./containers/SignIn'), uniOpts)
const SignOut = universal(import('./containers/SignOut'), uniOpts)

const styles = theme => {
  return {
    content: {
      marginTop: 0,
      marginBottom: theme.spacing.unit * 5
    }
  }
}

const Layout = useHooks(props => {
  const [headerHeight, resizeHeaderHeight] = useState(0)

  const onHeaderResize = (_width, height) => {
    resizeHeaderHeight(height)
  }

  const { classes, redirectTo } = props

  return (<Root>
    <BrowserRouter>
    <Grid container>
      <CssBaseline />
      <Grid item xs={12}>
        <Header
          onResize={onHeaderResize}
          drawerContainer={props.drawerContainer}
          />
        <main
          ref={ref => { props.drawerContainer = ref }}
          style={{ paddingTop: headerHeight }}
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
                render={props => (
                  <DashboardsIndex
                  {...props}
                  recentJobRunsCount={5}
                  recentlyCreatedPageSize={4}
                  />
                  )}
                  />
              <PrivateRoute exact path='/jobs' component={JobsIndex} />
              <PrivateRoute exact path='/jobs/page/:jobPage' component={JobsIndex} />
              <PrivateRoute exact path='/jobs/new' component={JobsNew} />
              <PrivateRoute
                exact
                path='/jobs/:jobSpecId'
                render={props => <JobsShow {...props} showJobRunsCount={5} />}
                />
              <PrivateRoute exact path='/jobs/:jobSpecId/definition' component={JobsDefinition} />
              <PrivateRoute exact path='/jobs/:jobSpecId/runs' component={JobRunsIndex} />
              <PrivateRoute exact path='/jobs/:jobSpecId/runs/page/:jobRunsPage' component={JobRunsIndex} />
              <PrivateRoute exact path='/jobs/:jobSpecId/runs/id/:jobRunId' component={JobRunsShow} />
              <PrivateRoute exact path='/jobs/:jobSpecId/runs/id/:jobRunId/json' component={JobRunsShowJson} />
              <PrivateRoute exact path='/bridges' component={BridgesIndex} />
              <PrivateRoute exact path='/bridges/page/:bridgePage' component={BridgesIndex} />
              <PrivateRoute exact path='/bridges/new' component={BridgesNew} />
              <PrivateRoute exact path='/bridges/:bridgeId' component={BridgesShow} />
              <PrivateRoute exact path='/bridges/:bridgeId/edit' component={BridgesEdit} />
              <PrivateRoute exact path='/config' component={Configuration} />
            </Switch>
          </div>
        </main>
      </Grid>
    </Grid>
    </BrowserRouter>
    </Root>
)
})

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

export default withStyles(styles)(ConnectedLayout)
