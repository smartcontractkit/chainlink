import React from 'react'
import { PaddedCard } from '@chainlink/styleguide'
import PrettyJson from 'components/PrettyJson'
import { DirectRequestJobRun, PipelineJobRun } from '../sharedTypes'

export const Json = ({
  jobRun,
}: {
  jobRun: DirectRequestJobRun | PipelineJobRun
}) => {
  return (
    <PaddedCard>
      <PrettyJson object={jobRun} />
    </PaddedCard>
  )
}
