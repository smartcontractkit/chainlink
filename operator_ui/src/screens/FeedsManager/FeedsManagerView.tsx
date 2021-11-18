import React from 'react'

import Grid from '@material-ui/core/Grid'

import { FeedsManagerCard } from './FeedsManagerCard'
import { JobProposalsCard } from './JobProposalsCard'

interface Props {
  manager: any // This is a hack to fix the typecheck until we can merged the fix
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
