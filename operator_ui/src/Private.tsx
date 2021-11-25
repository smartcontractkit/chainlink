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
import Configuration from 'pages/Configuration/Index'
import JobsShow from 'pages/Jobs/Show'
import JobsNew from 'pages/Jobs/New'
import JobRunsIndex from 'pages/JobRuns/Index'
import JobRunsShowOverview from 'pages/Jobs/Runs/Show'
import ChainsIndex from 'pages/ChainsIndex/ChainsIndex'
import ChainsNew from 'pages/Chains/New'
import ChainShow from 'pages/Chains/Show'
import { NodeScreen } from 'screens/Node/NodeScreen'
import KeysIndex from 'pages/Keys/Index'
import NotFound from 'pages/NotFound'
import TransactionsIndex from 'pages/Transactions/Index'
import TransactionsShow from 'pages/Transactions/Show'
import NodesIndex from './pages/NodesIndex/NodesIndex'

import { BridgesPage } from 'pages/bridges'
import { JobsPage } from 'pages/JobsIndex'
import { FeedsManagerPage } from 'pages/feeds_manager'
import { JobProposalsPage } from 'pages/job_proposals'

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
              <PrivateRoute exact path="/chains/new">
                <ChainsNew />
              </PrivateRoute>

              <PrivateRoute path="/chains/:chainId">
                <ChainShow />
              </PrivateRoute>

              <PrivateRoute exact path="/nodes">
                <NodesIndex />
              </PrivateRoute>

              <PrivateRoute path="/nodes/:id">
                <NodeScreen />
              </PrivateRoute>

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

              <PrivateRoute path="/bridges">
                <BridgesPage />
              </PrivateRoute>

              <PrivateRoute path="/feeds_manager">
                <FeedsManagerPage />
              </PrivateRoute>

              <PrivateRoute path="/job_proposals">
                <JobProposalsPage />
              </PrivateRoute>

              <PrivateRoute exact path="/jobs">
                <JobsPage />
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
