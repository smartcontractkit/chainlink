import React from 'react'
import { Theme } from '@material-ui/core/styles'
import Grid from '@material-ui/core/Grid'

import { Switch } from 'react-router-dom'
import { hot } from 'react-hot-loader'
import { withStyles } from '@material-ui/core/styles'
import Header from 'pages/Header'
import Notifications from 'pages/Notifications'
import PrivateRoute from './PrivateRoute'

import DashboardIndex from 'pages/Dashboards/Index'
import BridgesIndex from 'pages/Bridges/Index'
import BridgesShow from 'pages/Bridges/Show'
import BridgesNew from 'pages/Bridges/New'
import BridgesEdit from 'pages/Bridges/Edit'
import Configuration from 'pages/Configuration/Index'
import JobsIndex from 'pages/JobsIndex/JobsIndex'
import JobsShow from 'pages/Jobs/Show'
import JobsNew from 'pages/Jobs/New'
import JobRunsIndex from 'pages/JobRuns/Index'
import JobRunsShowOverview from 'pages/Jobs/Runs/Show'
import ChainsIndex from 'pages/ChainsIndex/ChainsIndex'
import KeysIndex from 'pages/Keys/Index'
import NotFound from 'pages/NotFound'
import TransactionsIndex from 'pages/Transactions/Index'
import TransactionsShow from 'pages/Transactions/Show'
import { FeedsManagerScreen } from 'pages/FeedsManager/FeedsManagerScreen'
import { JobProposalScreen } from 'pages/JobProposal/JobProposalScreen'

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
              <PrivateRoute exact path="/">
                <DashboardIndex
                  recentJobRunsCount={5}
                  recentlyCreatedPageSize={4}
                />
              </PrivateRoute>

              <PrivateRoute exact path="/jobs">
                <JobsIndex />
              </PrivateRoute>
              <PrivateRoute exact path="/jobs/new">
                <JobsNew />
              </PrivateRoute>

              <PrivateRoute
                path="/jobs/:jobId/runs/:jobRunId"
                component={JobRunsShowOverview}
              />

              <PrivateRoute path="/jobs/:jobId">
                <JobsShow />
              </PrivateRoute>

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

              <PrivateRoute exact path="/chains">
                <ChainsIndex />
              </PrivateRoute>

              <PrivateRoute exact path="/bridges" component={BridgesIndex} />

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
              <PrivateRoute exact path="/keys" component={KeysIndex} />
              <PrivateRoute exact path="/config" component={Configuration} />
              <PrivateRoute path="/feeds_manager">
                <FeedsManagerScreen />
              </PrivateRoute>

              <PrivateRoute path="/job_proposals/:id">
                <JobProposalScreen />
              </PrivateRoute>

              <PrivateRoute component={NotFound} />
            </Switch>
          </div>
        </main>
      </Grid>
    </Grid>
  )
}

export default hot(module)(withStyles(styles)(Private))
