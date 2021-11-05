import React from 'react'

import Grid from '@material-ui/core/Grid'

import { FeedsManagerCard } from './FeedsManagerCard'
import { FeedsManager } from 'types/generated/graphql'
import { JobProposalsCard } from './JobProposalsCard'

interface Props {
  manager: FeedsManager
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
