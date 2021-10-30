import React from 'react'

import Grid from '@material-ui/core/Grid'

import { FeedsManagerCard } from './FeedsManagerCard'
import { JobProposalsCard } from './JobProposalsCard'

import { FeedsManager } from 'types/feeds_manager'

interface Props {
  data: FeedsManager
}

export const FeedsManagerView: React.FC<Props> = ({ data }) => {
  return (
    <Grid container>
      <Grid item xs={12} lg={8}>
        <JobProposalsCard />
      </Grid>
      <Grid item xs={12} lg={4}>
        <FeedsManagerCard manager={data} />
      </Grid>
    </Grid>
  )
}
