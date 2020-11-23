import React from 'react'
import { PaddedCard } from '@chainlink/styleguide'
import PrettyJson from 'components/PrettyJson'
import { DirectRequestJobRun } from '../sharedTypes'

export const Json = ({ jobRun }: { jobRun: DirectRequestJobRun }) => {
  return (
    <PaddedCard>
      <PrettyJson object={jobRun} />
    </PaddedCard>
  )
}
