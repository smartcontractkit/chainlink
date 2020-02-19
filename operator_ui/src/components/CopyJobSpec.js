import React from 'react'
import Copy from 'components/Copy'

const CopyJobSpec = ({ JobSpec, ...props }) => (
  <Copy
    buttonText="Copy JobSpec"
    data={JSON.stringify(JobSpec, null, '\t')}
    {...props}
  />
)

export default CopyJobSpec
