import React from 'react'
import { RouteComponentProps } from 'react-router-dom'
import { Theme } from '@material-ui/core/styles'
import Grid from '@material-ui/core/Grid'
import universal, { UniversalComponent } from 'react-universal-component'
import { StyledComponentProps } from '@material-ui/core/styles'

import { Switch } from 'react-router-dom'
import { hot } from 'react-hot-loader'
import { withStyles } from '@material-ui/core/styles'
import Header from 'containers/Header'
import Loading from 'components/Loading'
import Notifications from 'containers/Notifications'
import { ShowComponentProps } from 'containers/Jobs/Show'
import PrivateRoute from './PrivateRoute'

// Asynchronously load routes that are chunked via code-splitting
// 'import' as a function must take a string. It can't take a variable.

const uniOpts = { loading: Loading }
const DashboardsIndex: UniversalComponent<
  Pick<
    {
      classes: any
    },
    never
  > &
    StyledComponentProps<'wrapper' | 'text'> & {
      recentJobRunsCount: number
      recentlyCreatedPageSize: number
    }
> = universal(import('./containers/Dashboards/Index'), uniOpts)
const JobsIndex = universal(import('./containers/Jobs/Index'), uniOpts)
const JobsShow: UniversalComponent<ShowComponentProps> &
  StyledComponentProps<'wrapper' | 'text'> = universal(
  import('./containers/Jobs/Show'),
  uniOpts,
)
const JobsDefinition = universal(
  import('./containers/Jobs/Definition'),
  uniOpts,
)
const JobsErrors = universal(import('./containers/Jobs/Errors'), uniOpts)
const JobsNew = universal(import('./containers/Jobs/New'), uniOpts)
const BridgesIndex = universal(import('./containers/Bridges/Index'), uniOpts)
const BridgesNew = universal(import('./containers/Bridges/New'), uniOpts)
const BridgesShow = universal(import('./containers/Bridges/Show'), uniOpts)
const BridgesEdit = universal(import('./containers/Bridges/Edit'), uniOpts)
const JobRunsIndex: UniversalComponent<
  Pick<
    {
      classes: any
    },
    never
  > &
    StyledComponentProps<'wrapper' | 'text'> & {
      pagePath: string
    }
> = universal(import('./containers/JobRuns/Index'), uniOpts)
const JobRunsShowOverview = universal(
  import('./containers/JobRuns/Show/Overview'),
  uniOpts,
)
const JobRunsShowJson = universal(
  import('./containers/JobRuns/Show/Json'),
  uniOpts,
)
const JobRunsShowErrorLog = universal(
  import('./containers/JobRuns/Show/ErrorLog'),
  uniOpts,
)
const TransactionsIndex = universal(
  import('./containers/Transactions/Index'),
  uniOpts,
)
const TransactionsShow = universal(
  import('./containers/Transactions/Show'),
  uniOpts,
)
const Configuration = universal(
  import('./containers/Configuration/Index'),
  uniOpts,
)
const NotFound = universal(import('./containers/NotFound'), uniOpts)

const styles = (theme: Theme) => {
  return {
    content: {
      marginTop: 0,
      marginBottom: theme.spacing.unit * 5,
    },
  }
}

const Private = ({ classes }: { classes: { content: string } }) => {
  const [headerHeight, setHeaderHeight] = React.useState(0)
  let drawerContainerRef: HTMLElement | null = null

  return (
    <Grid container spacing={0}>
      <Grid item xs={12}>
        <Header
          onResize={(_width, height) => setHeaderHeight(height)}
          drawerContainer={drawerContainerRef}
        />
        <main
          ref={(ref) => {
            drawerContainerRef = ref
          }}
          style={{ paddingTop: headerHeight }}
        >
          <Notifications />

          <div className={classes.content}>
            <Switch>
              <PrivateRoute
                exact
                path="/"
                render={(props) => (
                  <DashboardsIndex
                    {...props}
                    recentJobRunsCount={5}
                    recentlyCreatedPageSize={4}
                  />
                )}
              />
              <PrivateRoute exact path="/jobs" component={JobsIndex} />
              <PrivateRoute
                exact
                path="/jobs/page/:jobPage"
                component={JobsIndex}
              />
              <PrivateRoute exact path="/jobs/new" component={JobsNew} />
              <PrivateRoute
                exact
                path="/jobs/:jobSpecId"
                render={(
                  props: RouteComponentProps<{
                    jobSpecId: string
                  }>,
                ) => <JobsShow {...props} showJobRunsCount={5} />}
              />
              <PrivateRoute
                exact
                path="/jobs/:jobSpecId/json"
                component={JobsDefinition}
              />
              <PrivateRoute
                exact
                path="/jobs/:jobSpecId/errors"
                component={JobsErrors}
              />
              <PrivateRoute
                exact
                path="/jobs/:jobSpecId/runs"
                render={(props) => (
                  <JobRunsIndex
                    {...props}
                    pagePath="/jobs/:jobSpecId/runs/page"
                  />
                )}
              />
              <PrivateRoute
                exact
                path="/jobs/:jobSpecId/runs/page/:jobRunsPage"
                render={(props) => (
                  <JobRunsIndex
                    {...props}
                    pagePath="/jobs/:jobSpecId/runs/page"
                  />
                )}
              />
              <PrivateRoute
                exact
                path="/jobs/:jobSpecId/runs/id/:jobRunId"
                component={JobRunsShowOverview}
              />
              <PrivateRoute
                exact
                path="/jobs/:jobSpecId/runs/id/:jobRunId/json"
                component={JobRunsShowJson}
              />
              <PrivateRoute
                exact
                path="/jobs/:jobSpecId/runs/id/:jobRunId/error_log"
                component={JobRunsShowErrorLog}
              />
              <PrivateRoute
                exact
                path="/runs"
                render={(props) => (
                  <JobRunsIndex {...props} pagePath="/runs/page" />
                )}
              />
              <PrivateRoute
                exact
                path="/runs/page/:jobRunsPage"
                render={(props) => (
                  <JobRunsIndex {...props} pagePath="/runs/page" />
                )}
              />
              ;<PrivateRoute exact path="/bridges" component={BridgesIndex} />
              <PrivateRoute
                exact
                path="/bridges/page/:bridgePage"
                component={BridgesIndex}
              />
              <PrivateRoute exact path="/bridges/new" component={BridgesNew} />
              <PrivateRoute
                exact
                path="/bridges/:bridgeId"
                component={BridgesShow}
              />
              <PrivateRoute
                exact
                path="/bridges/:bridgeId/edit"
                component={BridgesEdit}
              />
              <PrivateRoute
                exact
                path="/transactions"
                component={TransactionsIndex}
              />
              <PrivateRoute
                exact
                path="/transactions/page/:transactionsPage"
                component={TransactionsIndex}
              />
              <PrivateRoute
                exact
                path="/transactions/:transactionId"
                component={TransactionsShow}
              />
              <PrivateRoute exact path="/config" component={Configuration} />
              <PrivateRoute component={NotFound} />
            </Switch>
          </div>
        </main>
      </Grid>
    </Grid>
  )
}

export default hot(module)(withStyles(styles)(Private))
