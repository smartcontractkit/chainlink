import React from 'react'

import { gql } from '@apollo/client'

import Grid from '@material-ui/core/Grid'

import { FeedsManagerCard, FEEDS_MANAGER_FIELDS } from './FeedsManagerCard'
import {
  JobProposalsCard,
  FEEDS_MANAGER__JOB_PROPOSAL_FIELDS,
} from './JobProposalsCard'
import { Heading1 } from 'src/components/Heading/Heading1'
import { Heading2 } from 'src/components/Heading/Heading2'

export const FEEDS_MANAGERS_PAYLOAD__RESULTS_FIELDS = gql`
  ${FEEDS_MANAGER_FIELDS}
  ${FEEDS_MANAGER__JOB_PROPOSAL_FIELDS}
  fragment FeedsManagerPayload_ResultsFields on FeedsManager {
    ...FeedsManagerFields
    jobProposals {
      ...FeedsManager_JobProposalsFields
    }
  }
`

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
        <Heading2>Job Proposals</Heading2>
      </Grid>
      <Grid item xs={12}>
        <JobProposalsCard proposals={manager.jobProposals} />
      </Grid>
    </Grid>
  )
}
