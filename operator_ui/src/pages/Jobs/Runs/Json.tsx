import React from 'react'
import { PaddedCard } from 'components/PaddedCard'
import PrettyJson from 'components/PrettyJson'

export const Json = ({ jobRun }: { jobRun: object }) => {
  return (
    <PaddedCard>
      <PrettyJson object={jobRun} />
    </PaddedCard>
  )
}
