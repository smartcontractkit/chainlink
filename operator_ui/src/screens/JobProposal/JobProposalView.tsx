import React from 'react'

import { gql } from '@apollo/client'

import Grid from '@material-ui/core/Grid'

import Content from 'src/components/Content'
import { Props as EditJobSpecDialogProps } from './EditJobSpecDialog'
import { Heading1 } from 'src/components/Heading/Heading1'
import { Heading2 } from 'src/components/Heading/Heading2'
import { JobProposalCard } from './JobProposalCard'
import { SpecsView, JOB_PROPOSAL__SPECS_FIELDS } from './SpecsView'

export const JOB_PROPOSAL_PAYLOAD_FIELDS = gql`
  ${JOB_PROPOSAL__SPECS_FIELDS}
  fragment JobProposalPayloadFields on JobProposal {
    id
    externalJobID
    remoteUUID
    jobID
    specs {
      ...JobProposal_SpecsFields
    }
    status
  }
`

interface Props {
  proposal: JobProposalPayloadFields
  onUpdateSpec: EditJobSpecDialogProps['onSubmit']
  onApprove: (specID: string) => void
  onCancel: (specID: string) => void
  onReject: (specID: string) => void
}

export const JobProposalView: React.FC<Props> = ({
  onApprove,
  onCancel,
  onReject,
  onUpdateSpec,
  proposal,
}) => {
  return (
    <Content>
      <Grid container spacing={32}>
        <Grid item xs={9}>
          <Heading1>Job Proposal #{proposal.id}</Heading1>
        </Grid>
      </Grid>

      <JobProposalCard proposal={proposal} />

      <Grid container spacing={32}>
        <Grid item xs={9}>
          <Heading2>Specs</Heading2>
        </Grid>
      </Grid>

      <Grid container spacing={32}>
        <Grid item xs={12}>
          <SpecsView
            specs={proposal.specs}
            onReject={onReject}
            onApprove={onApprove}
            onCancel={onCancel}
            onUpdateSpec={onUpdateSpec}
          />
        </Grid>
      </Grid>
    </Content>
  )
}
