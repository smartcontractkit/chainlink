import React from 'react'

import { FeedsManagerCard } from './FeedsManagerCard'
import { JobProposalsCard } from './JobProposalsCard'

import Grid from '@material-ui/core/Grid'
import { FeedsManagerGQL } from './FeedsManagerScreen'

interface Props {
  manager: FeedsManagerGQL
}

export const FeedsManagerView: React.FC<Props> = ({ manager }) => {
  return (
    <Grid container>
      <Grid item xs={12} lg={8}>
        <JobProposalsCard />
      </Grid>
      <Grid item xs={12} lg={4}>
        <FeedsManagerCard manager={manager} />
      </Grid>
    </Grid>
  )
}
