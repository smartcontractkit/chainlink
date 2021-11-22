import React from 'react'

import { gql } from '@apollo/client'

import Grid from '@material-ui/core/Grid'

import { FeedsManagerCard, FEEDS_MANAGER_FIELDS } from './FeedsManagerCard'
import {
  JobProposalsCard,
  FEEDS_MANAGER__JOB_PROPOSAL_FIELDS,
} from './JobProposalsCard'

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
    <Grid container>
      <Grid item xs={12} lg={8}>
        <JobProposalsCard proposals={manager.jobProposals} />
      </Grid>
      <Grid item xs={12} lg={4}>
        <FeedsManagerCard manager={manager} />
      </Grid>
    </Grid>
  )
}
