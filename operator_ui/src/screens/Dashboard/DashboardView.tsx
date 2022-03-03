import React from 'react'
import Grid from '@material-ui/core/Grid'

import Content from 'components/Content'

import { AccountBalance } from './AccountBalance'
import { Activity } from './Activity'
import { BuildInfoFooter } from './BuildInfoFooter'
import { RecentJobs } from './RecentJobs'

export const DashboardView = () => {
  return (
    <Content>
      <Grid container>
        <Grid item xs={8}>
          <Activity />
        </Grid>
        <Grid item xs={4}>
          <Grid container>
            <Grid item xs={12}>
              <AccountBalance />
            </Grid>
            <Grid item xs={12}>
              <RecentJobs />
            </Grid>
          </Grid>
        </Grid>
      </Grid>

      <BuildInfoFooter />
    </Content>
  )
}
