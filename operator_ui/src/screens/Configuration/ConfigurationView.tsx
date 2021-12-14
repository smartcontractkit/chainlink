import React from 'react'
import Grid from '@material-ui/core/Grid'
import Content from 'components/Content'

import { LoggingCard } from 'src/pages/Configuration/LoggingCard'
import { JobRuns } from 'src/pages/Configuration/JobRuns'

import { ConfigurationCard } from './ConfigurationCard/ConfigurationCard'
import { NodeInfoCard } from './NodeInfoCard/NodeInfoCard'

export const ConfigurationView = () => {
  return (
    <Content>
      <Grid container>
        <Grid item sm={12} md={8}>
          <ConfigurationCard />
        </Grid>
        <Grid item sm={12} md={4}>
          <Grid container>
            <Grid item xs={12}>
              <NodeInfoCard />
            </Grid>
            <Grid item xs={12}>
              <JobRuns />
            </Grid>
            <Grid item xs={12}>
              <LoggingCard />
            </Grid>
          </Grid>
        </Grid>
      </Grid>
    </Content>
  )
}
