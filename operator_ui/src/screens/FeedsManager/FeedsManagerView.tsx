import React from 'react'

import Grid from '@material-ui/core/Grid'

import { SupportedChainsCard } from './SupportedChainsCard'
import { FeedsManagerCard } from './FeedsManagerCard'
import { JobProposalsCard } from './JobProposalsCard'
import { Heading1 } from 'src/components/Heading/Heading1'
import { Heading2 } from 'src/components/Heading/Heading2'

interface Props {
  manager: FeedsManagerPayload_ResultsFields
}

export const FeedsManagerView: React.FC<Props> = ({ manager }) => {
  return (
    <Grid container spacing={16}>
      <Grid item xs={12}>
        <Heading1>Feeds Manager</Heading1>
      </Grid>
      <Grid item xs={12}>
        <FeedsManagerCard manager={manager} />
      </Grid>
      <Grid item xs={12}>
        <SupportedChainsCard cfgs={manager.chainConfigs} mgrID={manager.id} />
      </Grid>
      <Grid item xs={12}>
        <Heading2>Job Proposals</Heading2>
      </Grid>
      <Grid item xs={12}>
        <JobProposalsCard proposals={manager.jobProposals} />
      </Grid>
    </Grid>
  )
}
