import React from 'react'
import { PaddedCard } from '@chainlink/styleguide'
import PrettyJson from 'components/PrettyJson'
import { DirectRequestJobRun, OffChainReportingJobRun } from '../sharedTypes'

export const Json = ({
  jobRun,
}: {
  jobRun: DirectRequestJobRun | OffChainReportingJobRun
}) => {
  return (
    <PaddedCard>
      <PrettyJson object={jobRun} />
    </PaddedCard>
  )
}
