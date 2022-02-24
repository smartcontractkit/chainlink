import React from 'react'

import Grid from '@material-ui/core/Grid'

import {
  DetailsCard,
  DetailsCardItemTitle,
  DetailsCardItemValue,
} from 'src/components/Cards/DetailsCard'
import Link from 'src/components/Link'
import titleize from 'src/utils/titleize'

interface Props {
  proposal: JobProposalPayloadFields
}

export const JobProposalCard = ({ proposal }: Props) => {
  const approvedSpec = proposal.specs.find((spec) => spec.status === 'APPROVED')

  return (
    <DetailsCard>
      <Grid container>
        <Grid item xs={12} sm={4} md={2}>
          <DetailsCardItemTitle title="Status" />
          <DetailsCardItemValue value={titleize(proposal.status)} />
        </Grid>
        <Grid item xs={12} sm={4} md={4}>
          <DetailsCardItemTitle title="FMS UUID" />
          <DetailsCardItemValue value={proposal.remoteUUID} />
        </Grid>
        <Grid item xs={12} sm={4} md={4}>
          <DetailsCardItemTitle title="External Job ID" />
          {proposal.jobID && proposal.externalJobID ? (
            <DetailsCardItemValue>
              <Link color="primary" href={`/jobs/${proposal.jobID}`}>
                {proposal.externalJobID}
              </Link>
            </DetailsCardItemValue>
          ) : (
            <DetailsCardItemValue value={proposal.externalJobID || '--'} />
          )}
        </Grid>
        <Grid item xs={12} sm={4} md={2}>
          <DetailsCardItemTitle title="Approved Version" />
          <DetailsCardItemValue
            value={approvedSpec ? approvedSpec.version : '--'}
          />
        </Grid>
      </Grid>
    </DetailsCard>
  )
}
