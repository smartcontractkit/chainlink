import React from 'react'
import Copy from 'components/Copy'

const CopyJobSpec = ({ JobSpec }) => (
  <Copy buttonText="Copy JobSpec" data={JSON.stringify(JobSpec, null, '\t')} />
)

export default CopyJobSpec
