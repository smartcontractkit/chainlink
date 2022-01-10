import React from 'react'
import { Theme } from '@material-ui/core/styles'
import Grid from '@material-ui/core/Grid'

import { Switch } from 'react-router-dom'
import { hot } from 'react-hot-loader'
import { withStyles } from '@material-ui/core/styles'
import Header from 'pages/Header'
import Notifications from 'pages/Notifications'
import PrivateRoute from './PrivateRoute'

import { ChainsScreen } from 'screens/Chains/ChainsScreen'
import ChainsNew from 'pages/Chains/New'
import ChainShow from 'pages/Chains/Show'
import NotFound from 'pages/NotFound'

import { BridgesPage } from 'pages/bridges'
import { ConfigPage } from 'pages/config'
import { DashboardPage } from 'pages/dashboard'
import { JobsPage } from 'pages/JobsIndex'
import { KeysPage } from 'pages/keys'
import { JobRunsPage } from 'pages/job_runs'
import { FeedsManagerPage } from 'pages/feeds_manager'
import { JobProposalsPage } from 'pages/job_proposals'
import { NodesPage } from 'pages/nodes'
import { TransactionsPage } from 'pages/Transactions'

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
                <DashboardPage />
              </PrivateRoute>

              <PrivateRoute exact path="/chains">
                <ChainsScreen />
              </PrivateRoute>
              <PrivateRoute exact path="/chains/new">
                <ChainsNew />
              </PrivateRoute>

              <PrivateRoute path="/chains/:chainId">
                <ChainShow />
              </PrivateRoute>

              <PrivateRoute path="/bridges">
                <BridgesPage />
              </PrivateRoute>

              <PrivateRoute path="/config">
                <ConfigPage />
              </PrivateRoute>

              <PrivateRoute path="/feeds_manager">
                <FeedsManagerPage />
              </PrivateRoute>

              <PrivateRoute path="/job_proposals">
                <JobProposalsPage />
              </PrivateRoute>

              <PrivateRoute path="/jobs">
                <JobsPage />
              </PrivateRoute>

              <PrivateRoute path="/runs">
                <JobRunsPage />
              </PrivateRoute>

              <PrivateRoute path="/keys">
                <KeysPage />
              </PrivateRoute>

              <PrivateRoute path="/nodes">
                <NodesPage />
              </PrivateRoute>

              <PrivateRoute path="/transactions">
                <TransactionsPage />
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
